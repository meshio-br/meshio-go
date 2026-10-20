package meshio

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type recorded struct {
	Method string
	Path   string
	Query  string
	Header http.Header
	Body   string
}

// newTestClient cria um cliente apontando para um servidor simulado, sem esperas reais.
func newTestClient(t *testing.T, h http.HandlerFunc, opts ...Option) (*Client, *[]recorded, *[]time.Duration) {
	t.Helper()
	var mu sync.Mutex
	reqs := []recorded{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		reqs = append(reqs, recorded{r.Method, r.URL.EscapedPath(), r.URL.RawQuery, r.Header.Clone(), string(b)})
		mu.Unlock()
		h(w, r)
	}))
	t.Cleanup(srv.Close)
	sleeps := []time.Duration{}
	c := NewClient(append([]Option{WithAPIKey("msk_test_abc"), WithBaseURL(srv.URL + "/")}, opts...)...)
	c.sleep = func(_ context.Context, d time.Duration) error {
		sleeps = append(sleeps, d)
		return nil
	}
	return c, &reqs, &sleeps
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func errBody(code, msg string, details map[string]any) map[string]any {
	if details == nil {
		details = map[string]any{}
	}
	return map[string]any{"error": map[string]any{"code": code, "message": msg, "details": details}, "request_id": "req_body"}
}

func TestDefaultHeaders(t *testing.T) {
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"version": "1.0.0", "ref": nil, "commit": nil, "released_at": nil})
	})
	v, err := c.Version.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if v.Version == nil || *v.Version != "1.0.0" {
		t.Fatalf("versão inesperada: %+v", v)
	}
	h := (*reqs)[0].Header
	if got := h.Get("Authorization"); got != "Bearer msk_test_abc" {
		t.Errorf("Authorization = %q", got)
	}
	if got := h.Get("User-Agent"); got != "meshio-go/"+SDKVersion {
		t.Errorf("User-Agent = %q", got)
	}
	if got := h.Get("Accept"); got != "application/json" {
		t.Errorf("Accept = %q", got)
	}
	if (*reqs)[0].Path != "/v1/version" {
		t.Errorf("path = %q", (*reqs)[0].Path)
	}
}

func TestEnvironmentConfig(t *testing.T) {
	t.Setenv("MESHIO_API_KEY", "msk_test_env")
	t.Setenv("MESHIO_BASE_URL", "http://example.invalid/")
	c := NewClient()
	if c.apiKey != "msk_test_env" || c.baseURL != "http://example.invalid" {
		t.Fatalf("config do ambiente não aplicada: %q %q", c.apiKey, c.baseURL)
	}
	if c.timeout != 30*time.Second || c.maxRetries != 2 {
		t.Fatalf("padrões inesperados: %v %d", c.timeout, c.maxRetries)
	}
	t.Setenv("MESHIO_API_KEY", "")
	t.Setenv("MESHIO_BASE_URL", "")
	c = NewClient()
	if c.baseURL != DefaultBaseURL {
		t.Fatalf("baseURL padrão = %q", c.baseURL)
	}
	if _, err := c.Version.Get(context.Background()); !errors.Is(err, ErrMissingAPIKey) {
		t.Fatalf("esperava ErrMissingAPIKey, veio %v", err)
	}
}

func TestSendGeneratesIdempotencyKeyAndReusesOnRetry(t *testing.T) {
	n := 0
	c, reqs, sleeps := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n++
		if n == 1 {
			writeJSON(w, 503, errBody("UNAVAILABLE", "indisponível", nil))
			return
		}
		w.Header().Set("Idempotent-Replayed", "true")
		writeJSON(w, 202, map[string]any{"id": "msg_1", "status": "queued", "channel_id": "chn_1"})
	})
	var meta ResponseMeta
	msg, err := c.Messages.Send(context.Background(), &SendMessageParams{
		ChannelID: "chn_1", To: "+5511988887777", Type: "text", Text: &TextParams{Body: "Olá"},
	}, WithResponseMeta(&meta))
	if err != nil {
		t.Fatal(err)
	}
	if len(*reqs) != 2 || len(*sleeps) != 1 {
		t.Fatalf("esperava 2 tentativas e 1 espera, veio %d e %d", len(*reqs), len(*sleeps))
	}
	k1, k2 := (*reqs)[0].Header.Get("Idempotency-Key"), (*reqs)[1].Header.Get("Idempotency-Key")
	if k1 == "" || k1 != k2 {
		t.Fatalf("chave deve ser gerada e reutilizada: %q %q", k1, k2)
	}
	if len(k1) != 36 || k1[14] != '4' {
		t.Fatalf("chave não parece UUID v4: %q", k1)
	}
	if !msg.IdempotentReplay || !meta.IdempotentReplay || meta.IdempotencyKey != k1 {
		t.Fatalf("replay não exposto: %+v %+v", msg.IdempotentReplay, meta)
	}
	var body map[string]any
	_ = json.Unmarshal([]byte((*reqs)[0].Body), &body)
	if body["type"] != "text" || body["to"] != "+5511988887777" {
		t.Fatalf("corpo inesperado: %s", (*reqs)[0].Body)
	}
	if _, ok := body["template"]; ok {
		t.Fatalf("campos opcionais vazios não devem ir no corpo: %s", (*reqs)[0].Body)
	}
}

func TestSendUsesProvidedIdempotencyKey(t *testing.T) {
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 202, map[string]any{"id": "msg_1"})
	})
	_, err := c.Messages.Send(context.Background(), &SendMessageParams{ChannelID: "c", To: "+1", Type: "text"}, WithIdempotencyKey("order-1234"))
	if err != nil {
		t.Fatal(err)
	}
	if got := (*reqs)[0].Header.Get("Idempotency-Key"); got != "order-1234" {
		t.Fatalf("Idempotency-Key = %q", got)
	}
}

func TestRetryAfterHeaderIsHonored(t *testing.T) {
	n := 0
	c, reqs, sleeps := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n++
		if n < 3 {
			w.Header().Set("Retry-After", "3")
			writeJSON(w, 429, errBody("RATE_LIMITED", "muitas requisições", nil))
			return
		}
		writeJSON(w, 200, map[string]any{"data": []any{}, "next_cursor": nil})
	})
	if _, err := c.Businesses.List(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if len(*reqs) != 3 {
		t.Fatalf("tentativas = %d", len(*reqs))
	}
	for _, d := range *sleeps {
		if d != 3*time.Second {
			t.Fatalf("espera = %v, esperava 3s", d)
		}
	}
}

func TestRetriesExhaustedReturnsRateLimitError(t *testing.T) {
	c, reqs, sleeps := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req_hdr")
		writeJSON(w, 429, errBody("RATE_LIMITED", "muitas requisições", nil))
	})
	_, err := c.Businesses.Get(context.Background(), "biz_1")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || !errors.Is(err, ErrRateLimit) {
		t.Fatalf("esperava RateLimit, veio %v", err)
	}
	if apiErr.RequestID != "req_hdr" {
		t.Fatalf("request id do cabeçalho deve prevalecer: %q", apiErr.RequestID)
	}
	if len(*reqs) != 3 || len(*sleeps) != 2 {
		t.Fatalf("tentativas=%d esperas=%d", len(*reqs), len(*sleeps))
	}
	// backoff exponencial com jitter: ~0,5 s e ~1 s
	if (*sleeps)[0] < 375*time.Millisecond || (*sleeps)[0] > 625*time.Millisecond || (*sleeps)[1] < 750*time.Millisecond || (*sleeps)[1] > 1250*time.Millisecond {
		t.Fatalf("backoff fora do esperado: %v", *sleeps)
	}
}

func TestBackoffCap(t *testing.T) {
	for i := 0; i < 10; i++ {
		if d := backoff(i, 0); d > maxBackoff || d <= 0 {
			t.Fatalf("backoff(%d) = %v", i, d)
		}
	}
}

func TestNoRetryForNonIdempotentPOSTOn5xx(t *testing.T) {
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 500, errBody("INTERNAL", "erro interno", nil))
	})
	_, err := c.ConnectLinks.Create(context.Background(), &CreateConnectLinkParams{BusinessID: "biz_1"})
	if !errors.Is(err, ErrServer) {
		t.Fatalf("esperava ErrServer, veio %v", err)
	}
	if len(*reqs) != 1 {
		t.Fatalf("POST sem idempotência não deve ser repetido: %d", len(*reqs))
	}
}

func TestNoRetryOn4xx(t *testing.T) {
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 404, errBody("NOT_FOUND", "não encontrado", nil))
	})
	_, err := c.Channels.Get(context.Background(), "chn_x")
	if !errors.Is(err, ErrNotFound) || len(*reqs) != 1 {
		t.Fatalf("err=%v tentativas=%d", err, len(*reqs))
	}
}

func TestConnectionErrorRetriedForGET(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := srv.URL
	srv.Close()
	c := NewClient(WithAPIKey("k"), WithBaseURL(url), WithMaxRetries(1))
	waits := 0
	c.sleep = func(context.Context, time.Duration) error { waits++; return nil }
	_, err := c.Version.Get(context.Background())
	var ce *ConnectionError
	if !errors.As(err, &ce) || !errors.Is(err, ErrConnection) {
		t.Fatalf("esperava ConnectionError, veio %T %v", err, err)
	}
	if waits != 1 {
		t.Fatalf("esperava 1 nova tentativa, veio %d", waits)
	}
}

func TestTimeoutError(t *testing.T) {
	c, _, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(2 * time.Second):
		case <-r.Context().Done():
		}
	}, WithTimeout(50*time.Millisecond), WithMaxRetries(0))
	_, err := c.Version.Get(context.Background())
	var te *TimeoutError
	if !errors.As(err, &te) || !errors.Is(err, ErrTimeout) {
		t.Fatalf("esperava TimeoutError, veio %T %v", err, err)
	}
}

func TestErrorMapping(t *testing.T) {
	cases := []struct {
		status   int
		code     string
		sentinel error
	}{
		{400, "BAD_REQUEST", ErrBadRequest},
		{401, "UNAUTHENTICATED", ErrAuthentication},
		{403, "FORBIDDEN", ErrPermissionDenied},
		{404, "NOT_FOUND", ErrNotFound},
		{409, "CONFLICT", ErrConflict},
		{409, "IDEMPOTENCY_CONFLICT", ErrIdempotencyConflict},
		{422, "VALIDATION_FAILED", ErrValidation},
		{429, "RATE_LIMITED", ErrRateLimit},
		{502, "BAD_GATEWAY", ErrServer},
		{422, "POLICY_BLOCKED", ErrPolicyBlocked},
	}
	for _, tc := range cases {
		t.Run(tc.code, func(t *testing.T) {
			c, _, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				writeJSON(w, tc.status, errBody(tc.code, "mensagem do servidor", map[string]any{"x": 1}))
			}, WithMaxRetries(0))
			_, err := c.Templates.Get(context.Background(), "tpl_1")
			if !errors.Is(err, tc.sentinel) {
				t.Fatalf("errors.Is(%v, %v) = false", err, tc.sentinel)
			}
			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("esperava *APIError")
			}
			if apiErr.StatusCode != tc.status || apiErr.Code != tc.code || apiErr.Message != "mensagem do servidor" || apiErr.RequestID != "req_body" || apiErr.Details["x"] != float64(1) {
				t.Fatalf("campos inesperados: %+v", apiErr)
			}
			if tc.code == "POLICY_BLOCKED" && (!apiErr.IsPolicyBlocked() || !errors.Is(err, ErrValidation)) {
				t.Fatalf("POLICY_BLOCKED deve preservar a categoria do status")
			}
		})
	}
}

func TestValidationFields(t *testing.T) {
	bodies := []map[string]any{
		errBody("VALIDATION_FAILED", "dados inválidos", map[string]any{"fields": map[string]any{"to": "telefone inválido"}}),
		errBody("VALIDATION_FAILED", "dados inválidos", map[string]any{"fields": []any{map[string]any{"field": "to", "message": "telefone inválido"}}}),
		{"detail": []any{map[string]any{"loc": []any{"body", "to"}, "msg": "telefone inválido", "type": "x"}}},
		errBody("VALIDATION_FAILED", "telefone inválido", map[string]any{"field": "to"}),
	}
	for i, b := range bodies {
		c, _, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { writeJSON(w, 422, b) })
		_, err := c.Messages.Send(context.Background(), &SendMessageParams{})
		var apiErr *APIError
		if !errors.As(err, &apiErr) || apiErr.Fields["to"] != "telefone inválido" {
			t.Fatalf("caso %d: fields = %+v (%v)", i, apiErr, err)
		}
		if !strings.Contains(err.Error(), "to: telefone inválido") {
			t.Fatalf("caso %d: mensagem sem campos: %s", i, err)
		}
	}
}

func TestNonJSONErrorBody(t *testing.T) {
	c, _, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(502)
		_, _ = w.Write([]byte("<html>bad gateway</html>"))
	}, WithMaxRetries(0))
	_, err := c.Version.Get(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 502 || apiErr.Message == "" {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func TestQueryAndPathEncoding(t *testing.T) {
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"data": []any{}, "next_cursor": nil})
	})
	from := time.Date(2026, 10, 1, 12, 0, 0, 0, time.FixedZone("BRT", -3*3600))
	_, err := c.Messages.List(context.Background(), &MessagesListParams{ChannelID: "chn_1", Direction: "inbound", From: from, Limit: 50})
	if err != nil {
		t.Fatal(err)
	}
	q := (*reqs)[0].Query
	for _, want := range []string{"channel_id=chn_1", "direction=inbound", "limit=50", "from=2026-10-01T15%3A00%3A00Z"} {
		if !strings.Contains(q, want) {
			t.Errorf("query %q sem %q", q, want)
		}
	}
	if strings.Contains(q, "cursor") || strings.Contains(q, "status") {
		t.Errorf("filtros vazios não devem ser enviados: %q", q)
	}
	_, _ = c.Conversations.Timeline(context.Background(), "chn_1", "+5511988887777", nil)
	if p := (*reqs)[1].Path; p != "/v1/conversations/chn_1/+5511988887777" {
		t.Errorf("path = %q", p)
	}
	_, _ = c.Businesses.Get(context.Background(), "a/b")
	if p := (*reqs)[2].Path; p != "/v1/businesses/a%2Fb" {
		t.Errorf("id deve ser escapado: %q", p)
	}
}

func TestPaginationIterator(t *testing.T) {
	pages := map[string]any{
		"":   map[string]any{"data": []any{map[string]any{"id": "biz_1"}, map[string]any{"id": "biz_2"}}, "next_cursor": "c2"},
		"c2": map[string]any{"data": []any{}, "next_cursor": "c3"},
		"c3": map[string]any{"data": []any{map[string]any{"id": "biz_3"}}, "next_cursor": nil},
	}
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, pages[r.URL.Query().Get("cursor")])
	})
	it := c.Businesses.ListAll(context.Background(), &BusinessesListParams{Limit: 2, Status: "active"})
	var ids []string
	for it.Next() {
		ids = append(ids, it.Current().ID)
	}
	if err := it.Err(); err != nil {
		t.Fatal(err)
	}
	if strings.Join(ids, ",") != "biz_1,biz_2,biz_3" {
		t.Fatalf("ids = %v", ids)
	}
	if len(*reqs) != 3 {
		t.Fatalf("requisições = %d", len(*reqs))
	}
	for _, r := range *reqs {
		if !strings.Contains(r.Query, "limit=2") || !strings.Contains(r.Query, "status=active") {
			t.Fatalf("filtros devem ser mantidos entre páginas: %q", r.Query)
		}
	}
}

func TestPaginationIteratorError(t *testing.T) {
	n := 0
	c, _, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n++
		if n == 1 {
			writeJSON(w, 200, map[string]any{"data": []any{map[string]any{"id": "tpl_1"}}, "next_cursor": "x"})
			return
		}
		writeJSON(w, 403, errBody("FORBIDDEN", "sem permissão", nil))
	})
	items, err := c.Templates.ListAll(context.Background(), nil).Collect()
	if len(items) != 1 || !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("items=%d err=%v", len(items), err)
	}
}

func TestCSVExport(t *testing.T) {
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte("date,count\n2026-10-01,3\n"))
	})
	b, err := c.Usage.ExportCSV(context.Background(), &UsageExportCSVParams{From: "2026-10-01"})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "date,count\n2026-10-01,3\n" || (*reqs)[0].Path != "/v1/usage/export.csv" {
		t.Fatalf("csv inesperado: %q", b)
	}
}

func TestDeleteNoContentAndMethods(t *testing.T) {
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			w.WriteHeader(204)
			return
		}
		writeJSON(w, 200, map[string]any{"id": "whe_1", "status": "disabled"})
	})
	if err := c.WebhookEndpoints.Delete(context.Background(), "whe_1"); err != nil {
		t.Fatal(err)
	}
	ep, err := c.WebhookEndpoints.Update(context.Background(), "whe_1", &UpdateWebhookEndpointParams{Status: "disabled"})
	if err != nil || ep.Status != "disabled" {
		t.Fatalf("update: %v %+v", err, ep)
	}
	if err := c.OptOuts.Delete(context.Background(), &OptOutsDeleteParams{BusinessID: "biz_1", Phone: "+5511999999999"}); err != nil {
		t.Fatal(err)
	}
	r := *reqs
	if r[0].Method != "DELETE" || r[1].Method != "PATCH" || r[1].Body != `{"status":"disabled"}` || r[2].Query != "business_id=biz_1&phone=%2B5511999999999" {
		t.Fatalf("requisições inesperadas: %+v", r)
	}
}

func TestTopupIdempotency(t *testing.T) {
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 201, map[string]any{"id": "top_1", "status": "pending", "amount_cents": 10000})
	})
	tp, err := c.Billing.Topups.Create(context.Background(), &TopupParams{AmountCents: 10000})
	if err != nil || tp.ID != "top_1" {
		t.Fatalf("%v %+v", err, tp)
	}
	if (*reqs)[0].Header.Get("Idempotency-Key") == "" {
		t.Fatal("topups.create deve enviar Idempotency-Key")
	}
}

func TestMicrosToReais(t *testing.T) {
	cases := map[string]string{
		"0": "0.00", "1000000": "1.00", "12500000": "12.50", "1": "0.000001",
		"-2500000": "-2.50", "123456789": "123.456789", "99999999999999999999": "99999999999999.999999",
	}
	for in, want := range cases {
		got, err := MicrosToReais(in)
		if err != nil || got != want {
			t.Errorf("MicrosToReais(%q) = %q, %v; esperava %q", in, got, err, want)
		}
	}
	if _, err := MicrosToReais("1.5"); err == nil {
		t.Error("esperava erro para valor não inteiro")
	}
	r, err := MicrosToRat("1500000")
	if err != nil || r.FloatString(2) != "1.50" {
		t.Errorf("MicrosToRat = %v %v", r, err)
	}
}

func TestOperationCoverage(t *testing.T) {
	// Cada operationId do contrato deve ter um método documentado no código.
	raw, err := readFile("api/openapi.public.json")
	if err != nil {
		t.Fatal(err)
	}
	var spec struct {
		Paths map[string]map[string]struct {
			OperationID string `json:"operationId"`
		} `json:"paths"`
	}
	if err := json.Unmarshal(raw, &spec); err != nil {
		t.Fatal(err)
	}
	src, err := readFile("services.go")
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, ops := range spec.Paths {
		for _, op := range ops {
			if op.OperationID == "" {
				continue
			}
			count++
			if !strings.Contains(string(src), "("+op.OperationID+")") {
				t.Errorf("operação sem método: %s", op.OperationID)
			}
		}
	}
	if count != 67 {
		t.Errorf("contrato com %d operações, esperava 67", count)
	}
}

func TestConnectionErrorBeforeSendRetriedForNonIdempotentPOST(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := srv.URL
	srv.Close()
	c := NewClient(WithAPIKey("k"), WithBaseURL(url))
	waits := 0
	c.sleep = func(context.Context, time.Duration) error { waits++; return nil }
	_, err := c.ConnectLinks.Create(context.Background(), &CreateConnectLinkParams{BusinessID: "biz_1"})
	var ce *ConnectionError
	if !errors.As(err, &ce) || !ce.BeforeSend() {
		t.Fatalf("esperava ConnectionError antes do envio, veio %T %v", err, err)
	}
	if waits != 2 {
		t.Fatalf("esperava 2 novas tentativas, veio %d", waits)
	}
}

func TestTimeoutAfterSendNotRetriedForNonIdempotentPOST(t *testing.T) {
	var calls atomic.Int32
	c, _, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		select {
		case <-time.After(2 * time.Second):
		case <-r.Context().Done():
		}
	}, WithTimeout(50*time.Millisecond))
	_, err := c.ConnectLinks.Create(context.Background(), &CreateConnectLinkParams{BusinessID: "biz_1"})
	if !errors.Is(err, ErrTimeout) || calls.Load() != 1 {
		t.Fatalf("err=%v tentativas=%d", err, calls.Load())
	}
}

func TestTimeoutRetriedForGETAndIdempotentPOST(t *testing.T) {
	var mu sync.Mutex
	n := 0
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		n++
		first := n%2 == 1
		mu.Unlock()
		if first {
			select {
			case <-time.After(2 * time.Second):
			case <-r.Context().Done():
			}
			return
		}
		if r.Method == http.MethodGet {
			writeJSON(w, 200, map[string]any{"version": "1"})
			return
		}
		writeJSON(w, 202, map[string]any{"id": "msg_1", "status": "queued"})
	}, WithTimeout(100*time.Millisecond))
	if _, err := c.Version.Get(context.Background()); err != nil {
		t.Fatalf("GET: %v", err)
	}
	if _, err := c.Messages.Send(context.Background(), &SendMessageParams{ChannelID: "chn_1", To: "+5511", Type: "text"}); err != nil {
		t.Fatalf("send: %v", err)
	}
	if len(*reqs) != 4 {
		t.Fatalf("esperava 4 requisições, veio %d", len(*reqs))
	}
	if (*reqs)[2].Header.Get("Idempotency-Key") != (*reqs)[3].Header.Get("Idempotency-Key") {
		t.Fatalf("a chave deve ser reutilizada")
	}
}

func TestRateLimitRetriedForNonIdempotentPOST(t *testing.T) {
	var mu sync.Mutex
	n := 0
	c, reqs, sleeps := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		n++
		cur := n
		mu.Unlock()
		if cur == 1 {
			w.Header().Set("Retry-After", "3")
			writeJSON(w, 429, errBody("RATE_LIMITED", "too many requests", nil))
			return
		}
		writeJSON(w, 201, map[string]any{"id": "cl_1"})
	})
	if _, err := c.ConnectLinks.Create(context.Background(), &CreateConnectLinkParams{BusinessID: "biz_1"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(*reqs) != 2 || len(*sleeps) != 1 || (*sleeps)[0] != 3*time.Second {
		t.Fatalf("tentativas=%d esperas=%v", len(*reqs), *sleeps)
	}
}

func TestServerErrorRetriedForIdempotentSend(t *testing.T) {
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 503, errBody("UNAVAILABLE", "unavailable", nil))
	})
	_, err := c.Messages.Send(context.Background(), &SendMessageParams{})
	if !errors.Is(err, ErrServer) || len(*reqs) != 3 {
		t.Fatalf("err=%v tentativas=%d", err, len(*reqs))
	}
}

func TestMessagePtAndFieldsWithField(t *testing.T) {
	c, _, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 422, errBody("VALIDATION_FAILED", "Invalid phone", map[string]any{
			"field": "to", "message_pt": "Telefone inválido.", "fields": map[string]any{"to": "Use o formato +5511999999999."},
		}))
	})
	_, err := c.Messages.Send(context.Background(), &SendMessageParams{})
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("esperava *APIError: %v", err)
	}
	if apiErr.Message != "Invalid phone" || apiErr.MessagePt != "Telefone inválido." || apiErr.Fields["to"] != "Use o formato +5511999999999." {
		t.Fatalf("campos inesperados: %+v", apiErr)
	}
}
