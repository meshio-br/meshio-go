package meshio

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// Cabeçalhos enviados pela Meshio em cada entrega de webhook.
const (
	HeaderWebhookTimestamp = "Meshio-Timestamp"
	HeaderWebhookSignature = "Meshio-Signature"
	HeaderWebhookEventID   = "Meshio-Event-Id"
)

// DefaultWebhookTolerance é a diferença máxima aceita entre o relógio local e Meshio-Timestamp.
const DefaultWebhookTolerance = 300 * time.Second

// WebhookEvent é o envelope de um evento entregue por webhook.
type WebhookEvent struct {
	EventID    string          `json:"event_id"`
	EventType  string          `json:"event_type"`
	OccurredAt time.Time       `json:"occurred_at"`
	BusinessID *string         `json:"business_id"`
	ChannelID  *string         `json:"channel_id"`
	Data       json.RawMessage `json:"data"`
}

// ParsedData decodifica Data no tipo do evento (ex.: *MessageReceivedData para "message.received").
// Para tipos desconhecidos devolve map[string]any, sem erro.
func (e *WebhookEvent) ParsedData() (any, error) {
	v, ok, err := decodeEventData(e.EventType, e.Data)
	if ok || err != nil {
		return v, err
	}
	var m map[string]any
	if len(e.Data) == 0 {
		return map[string]any{}, nil
	}
	err = json.Unmarshal(e.Data, &m)
	return m, err
}

// DecodeData decodifica Data em v (ponteiro para o tipo desejado).
func (e *WebhookEvent) DecodeData(v any) error { return json.Unmarshal(e.Data, v) }

// WebhookOption ajusta a verificação de webhooks.
type WebhookOption func(*webhookConfig)

type webhookConfig struct {
	tolerance time.Duration
	now       func() time.Time
}

// WithTolerance altera a tolerância de horário (padrão 300 s).
func WithTolerance(d time.Duration) WebhookOption {
	return func(c *webhookConfig) { c.tolerance = d }
}

// WithClock injeta o relógio usado na verificação (útil em testes).
func WithClock(now func() time.Time) WebhookOption {
	return func(c *webhookConfig) { c.now = now }
}

// VerifyWebhookSignature confere a assinatura HMAC-SHA256 do corpo bruto sem decodificar o evento.
// headers aceita http.Header ou qualquer map[string][]string; nomes sem diferenciar maiúsculas.
func VerifyWebhookSignature(payload []byte, headers map[string][]string, secret string, opts ...WebhookOption) error {
	cfg := webhookConfig{tolerance: DefaultWebhookTolerance, now: time.Now}
	for _, o := range opts {
		o(&cfg)
	}
	if secret == "" {
		return &WebhookVerificationError{Reason: "segredo vazio"}
	}
	tsRaw := headerValue(headers, HeaderWebhookTimestamp)
	sigRaw := headerValue(headers, HeaderWebhookSignature)
	if tsRaw == "" {
		return &WebhookVerificationError{Reason: "cabeçalho Meshio-Timestamp ausente"}
	}
	if sigRaw == "" {
		return &WebhookVerificationError{Reason: "cabeçalho Meshio-Signature ausente"}
	}
	ts, err := strconv.ParseInt(strings.TrimSpace(tsRaw), 10, 64)
	if err != nil {
		return &WebhookVerificationError{Reason: "Meshio-Timestamp inválido"}
	}
	diff := cfg.now().Unix() - ts
	if diff < 0 {
		diff = -diff
	}
	if cfg.tolerance > 0 && time.Duration(diff)*time.Second > cfg.tolerance {
		return &WebhookVerificationError{Reason: "horário fora da tolerância"}
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(ts, 10)))
	mac.Write([]byte("."))
	mac.Write(payload)
	expected := mac.Sum(nil)
	for _, part := range strings.Split(sigRaw, ",") {
		part = strings.TrimSpace(part)
		if !strings.HasPrefix(part, "v1=") {
			continue
		}
		got, err := hex.DecodeString(strings.TrimPrefix(part, "v1="))
		if err != nil {
			continue
		}
		if hmac.Equal(got, expected) {
			return nil
		}
	}
	return &WebhookVerificationError{Reason: "nenhuma assinatura v1 confere"}
}

// VerifyWebhook confere a assinatura e devolve o evento decodificado.
// payload deve ser o corpo bruto, exatamente como recebido.
func VerifyWebhook(payload []byte, headers map[string][]string, secret string, opts ...WebhookOption) (*WebhookEvent, error) {
	if err := VerifyWebhookSignature(payload, headers, secret, opts...); err != nil {
		return nil, err
	}
	return ParseWebhookEvent(payload)
}

// ConstructEvent é sinônimo de VerifyWebhook.
func ConstructEvent(payload []byte, headers map[string][]string, secret string, opts ...WebhookOption) (*WebhookEvent, error) {
	return VerifyWebhook(payload, headers, secret, opts...)
}

// ParseWebhookEvent decodifica o envelope SEM verificar a assinatura.
func ParseWebhookEvent(payload []byte) (*WebhookEvent, error) {
	var ev WebhookEvent
	if err := json.Unmarshal(payload, &ev); err != nil {
		return nil, &WebhookVerificationError{Reason: "corpo não é um evento JSON válido"}
	}
	return &ev, nil
}

func headerValue(h map[string][]string, name string) string {
	for k, v := range h {
		if strings.EqualFold(k, name) && len(v) > 0 {
			return strings.Join(v, ",")
		}
	}
	return ""
}
