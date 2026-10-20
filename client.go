package meshio

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// SDKVersion é a versão desta biblioteca.
const SDKVersion = "0.1.0"

// DefaultBaseURL é o endereço da API de produção.
const DefaultBaseURL = "https://api.meshio.com.br"

const (
	defaultTimeout    = 30 * time.Second
	defaultMaxRetries = 2
	maxBackoff        = 8 * time.Second
	initialBackoff    = 500 * time.Millisecond
)

// Client é o cliente da API Meshio. Crie com NewClient; é seguro para uso concorrente.
type Client struct {
	apiKey     string
	baseURL    string
	timeout    time.Duration
	maxRetries int
	httpClient *http.Client
	userAgent  string
	sleep      func(ctx context.Context, d time.Duration) error

	Version           *VersionService
	Businesses        *BusinessesService
	Channels          *ChannelsService
	ConnectLinks      *ConnectLinksService
	Messages          *MessagesService
	Conversations     *ConversationsService
	Templates         *TemplatesService
	WebhookEndpoints  *WebhookEndpointsService
	WebhookDeliveries *WebhookDeliveriesService
	Events            *EventsService
	Usage             *UsageService
	Compliance        *ComplianceService
	Consents          *ConsentsService
	OptOuts           *OptOutsService
	Billing           *BillingService
}

// Option configura o Client.
type Option func(*Client)

// WithAPIKey define a chave de API (msk_live_… ou msk_test_…). Padrão: variável MESHIO_API_KEY.
func WithAPIKey(key string) Option { return func(c *Client) { c.apiKey = key } }

// WithBaseURL define o endereço da API. Padrão: MESHIO_BASE_URL ou https://api.meshio.com.br.
func WithBaseURL(u string) Option { return func(c *Client) { c.baseURL = u } }

// WithTimeout define o tempo máximo de cada tentativa (padrão 30 s).
func WithTimeout(d time.Duration) Option { return func(c *Client) { c.timeout = d } }

// WithMaxRetries define quantas novas tentativas são feitas em falhas transitórias (padrão 2; 0 desativa).
func WithMaxRetries(n int) Option {
	return func(c *Client) {
		if n < 0 {
			n = 0
		}
		c.maxRetries = n
	}
}

// WithHTTPClient usa um *http.Client próprio (proxy, transporte, testes). O timeout por tentativa continua valendo.
func WithHTTPClient(hc *http.Client) Option { return func(c *Client) { c.httpClient = hc } }

// NewClient cria um cliente. Sem WithAPIKey, usa MESHIO_API_KEY; sem chave, as chamadas falham com ErrMissingAPIKey.
func NewClient(opts ...Option) *Client {
	c := &Client{
		apiKey:     os.Getenv("MESHIO_API_KEY"),
		baseURL:    os.Getenv("MESHIO_BASE_URL"),
		timeout:    defaultTimeout,
		maxRetries: defaultMaxRetries,
		userAgent:  "meshio-go/" + SDKVersion,
		sleep:      sleepCtx,
	}
	for _, o := range opts {
		o(c)
	}
	if c.baseURL == "" {
		c.baseURL = DefaultBaseURL
	}
	c.baseURL = strings.TrimRight(c.baseURL, "/")
	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}
	c.Version = &VersionService{c}
	c.Businesses = &BusinessesService{c}
	c.Channels = &ChannelsService{client: c, SignupSessions: &SignupSessionsService{c}, MetaLogs: &MetaLogsService{c}}
	c.ConnectLinks = &ConnectLinksService{c}
	c.Messages = &MessagesService{c}
	c.Conversations = &ConversationsService{c}
	c.Templates = &TemplatesService{c}
	c.WebhookEndpoints = &WebhookEndpointsService{c}
	c.WebhookDeliveries = &WebhookDeliveriesService{c}
	c.Events = &EventsService{c}
	c.Usage = &UsageService{c}
	c.Compliance = &ComplianceService{c}
	c.Consents = &ConsentsService{c}
	c.OptOuts = &OptOutsService{c}
	c.Billing = &BillingService{
		client:     c,
		Entries:    &BillingEntriesService{c},
		NumberFees: &BillingNumberFeesService{c},
		Settings:   &BillingSettingsService{c},
		Topups:     &TopupsService{c},
	}
	return c
}

// ResponseMeta traz metadados da resposta HTTP (ver WithResponseMeta).
type ResponseMeta struct {
	StatusCode int
	Header     http.Header
	RequestID  string
	// IdempotencyKey é a chave enviada (gerada pelo SDK quando não informada).
	IdempotencyKey string
	// IdempotentReplay indica que o servidor devolveu o resultado de uma chamada anterior com a mesma chave.
	IdempotentReplay bool
}

// RequestOption ajusta uma chamada específica.
type RequestOption func(*requestOptions)

type requestOptions struct {
	idempotencyKey string
	metas          []*ResponseMeta
}

// WithIdempotencyKey define a chave de idempotência (messages.send e topups.create).
// Sem ela, o SDK gera um UUID v4 e reutiliza a mesma chave nas novas tentativas.
func WithIdempotencyKey(key string) RequestOption {
	return func(o *requestOptions) { o.idempotencyKey = key }
}

// WithResponseMeta preenche meta com status, cabeçalhos, request id e dados de idempotência.
func WithResponseMeta(meta *ResponseMeta) RequestOption {
	return func(o *requestOptions) {
		if meta != nil {
			o.metas = append(o.metas, meta)
		}
	}
}

type call struct {
	method     string
	path       string
	query      url.Values
	body       any
	idempotent bool // envia Idempotency-Key e permite novas tentativas em POST
	raw        bool // resposta não-JSON (CSV)
}

// ErrMissingAPIKey indica que nenhuma chave de API foi configurada.
var ErrMissingAPIKey = errors.New("meshio: chave de API ausente (use WithAPIKey ou MESHIO_API_KEY)")

func (c *Client) do(ctx context.Context, cl call, out any, opts []RequestOption) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if c.apiKey == "" {
		return ErrMissingAPIKey
	}
	ro := requestOptions{}
	for _, o := range opts {
		o(&ro)
	}
	var payload []byte
	if cl.body != nil {
		b, err := json.Marshal(cl.body)
		if err != nil {
			return fmt.Errorf("meshio: serializar corpo: %w", err)
		}
		payload = b
	}
	u := c.baseURL + cl.path
	if len(cl.query) > 0 {
		u += "?" + cl.query.Encode()
	}
	if cl.idempotent && ro.idempotencyKey == "" {
		ro.idempotencyKey = newUUID()
	}
	safe := cl.method == http.MethodGet || cl.method == http.MethodHead
	canRetryAfterSend := safe || (cl.idempotent && ro.idempotencyKey != "")

	for attempt := 0; ; attempt++ {
		resp, err := c.attempt(ctx, cl, u, payload, ro.idempotencyKey)
		if err != nil {
			// Falha de conexão antes do envio: sempre repetível. Timeout/erro após o envio:
			// só em GET ou com Idempotency-Key.
			if ctx.Err() == nil && attempt < c.maxRetries && (notSent(err) || canRetryAfterSend) {
				if serr := c.sleep(ctx, backoff(attempt, 0)); serr != nil {
					return serr
				}
				continue
			}
			return err
		}
		body, rerr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if rerr != nil {
			return &ConnectionError{Err: rerr}
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			for _, m := range ro.metas {
				fillMeta(m, resp, ro.idempotencyKey)
			}
			return decodeOK(body, out, cl.raw)
		}
		apiErr := newAPIError(resp, body)
		retryable := resp.StatusCode == 429 || (canRetryAfterSend && isRetryableStatus(resp.StatusCode))
		if retryable && attempt < c.maxRetries {
			if serr := c.sleep(ctx, backoff(attempt, apiErr.RetryAfter)); serr != nil {
				return serr
			}
			continue
		}
		for _, m := range ro.metas {
			fillMeta(m, resp, ro.idempotencyKey)
		}
		return apiErr
	}
}

func (c *Client) attempt(ctx context.Context, cl call, u string, payload []byte, idemKey string) (*http.Response, error) {
	actx := ctx
	var cancel context.CancelFunc
	if c.timeout > 0 {
		actx, cancel = context.WithTimeout(ctx, c.timeout)
	}
	var rd io.Reader
	if payload != nil {
		rd = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(actx, cl.method, u, rd)
	if err != nil {
		if cancel != nil {
			cancel()
		}
		return nil, fmt.Errorf("meshio: montar requisição: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)
	if cl.raw {
		req.Header.Set("Accept", "text/csv, application/json")
	} else {
		req.Header.Set("Accept", "application/json")
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if idemKey != "" {
		req.Header.Set("Idempotency-Key", idemKey)
	}
	var wrote atomic.Bool
	trace := &httptrace.ClientTrace{WroteHeaders: func() { wrote.Store(true) }}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if cancel != nil {
			cancel()
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		var ne net.Error
		if errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &ne) && ne.Timeout()) {
			return nil, &TimeoutError{Err: err, beforeSend: !wrote.Load()}
		}
		return nil, &ConnectionError{Err: err, beforeSend: !wrote.Load()}
	}
	if cancel != nil {
		resp.Body = &cancelBody{ReadCloser: resp.Body, cancel: cancel}
	}
	return resp, nil
}

// notSent indica falha de conexão ocorrida antes de a requisição ser enviada.
func notSent(err error) bool {
	var ce *ConnectionError
	if errors.As(err, &ce) {
		return ce.beforeSend
	}
	var te *TimeoutError
	return errors.As(err, &te) && te.beforeSend
}

type cancelBody struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (b *cancelBody) Close() error {
	err := b.ReadCloser.Close()
	b.cancel()
	return err
}

func fillMeta(m *ResponseMeta, resp *http.Response, key string) {
	m.StatusCode = resp.StatusCode
	m.Header = resp.Header
	m.RequestID = resp.Header.Get("X-Request-Id")
	m.IdempotencyKey = key
	m.IdempotentReplay = strings.EqualFold(resp.Header.Get("Idempotent-Replayed"), "true")
}

func decodeOK(body []byte, out any, raw bool) error {
	if out == nil {
		return nil
	}
	if raw {
		if p, ok := out.(*[]byte); ok {
			*p = body
			return nil
		}
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("meshio: decodificar resposta: %w", err)
	}
	return nil
}

func isRetryableStatus(s int) bool {
	switch s {
	case 429, 500, 502, 503, 504:
		return true
	}
	return false
}

// backoff: Retry-After quando houver; senão 0,5 s, 1 s, 2 s… (máx. 8 s) com jitter.
func backoff(attempt int, retryAfter time.Duration) time.Duration {
	if retryAfter > 0 {
		return retryAfter
	}
	d := initialBackoff << uint(attempt)
	if d > maxBackoff || d <= 0 {
		d = maxBackoff
	}
	// jitter de ±25 %
	j := time.Duration(mrand.Int63n(int64(d)/2+1)) - d/4
	d += j
	if d > maxBackoff {
		d = maxBackoff
	}
	return d
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func parseRetryAfter(v string) time.Duration {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	if s, err := strconv.ParseFloat(v, 64); err == nil && s >= 0 {
		return time.Duration(s * float64(time.Second))
	}
	if t, err := http.ParseTime(v); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

// newUUID gera um UUID v4.
func newUUID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func ep(segments ...string) string {
	var sb strings.Builder
	for i, s := range segments {
		if i%2 == 0 {
			sb.WriteString(s)
		} else {
			sb.WriteString(url.PathEscape(s))
		}
	}
	return sb.String()
}
