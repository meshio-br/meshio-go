package meshio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Sentinelas por categoria. Use errors.Is(err, meshio.ErrNotFound) etc.
// Para os detalhes (code, request id, campos), use errors.As com *APIError.
var (
	ErrBadRequest          = errors.New("meshio: requisição inválida")      // 400
	ErrAuthentication      = errors.New("meshio: falha de autenticação")    // 401
	ErrPermissionDenied    = errors.New("meshio: permissão negada")         // 403
	ErrNotFound            = errors.New("meshio: recurso não encontrado")   // 404
	ErrConflict            = errors.New("meshio: conflito")                 // 409
	ErrIdempotencyConflict = errors.New("meshio: conflito de idempotência") // 409 IDEMPOTENCY_CONFLICT
	ErrValidation          = errors.New("meshio: dados inválidos")          // 422
	ErrRateLimit           = errors.New("meshio: limite de requisições")    // 429
	ErrServer              = errors.New("meshio: erro no servidor")         // 5xx
	ErrPolicyBlocked       = errors.New("meshio: bloqueado por política")   // code POLICY_BLOCKED
	ErrConnection          = errors.New("meshio: falha de conexão")         // rede
	ErrTimeout             = errors.New("meshio: tempo esgotado")           // timeout
	ErrWebhookVerification = errors.New("meshio: assinatura de webhook inválida")
)

// APIError é o erro devolvido pela API (corpo {"error":{"code","message","details"},"request_id"}).
type APIError struct {
	// StatusCode é o status HTTP.
	StatusCode int
	// Code é o código técnico (ex.: POLICY_BLOCKED, IDEMPOTENCY_CONFLICT).
	Code string
	// Message é a mensagem técnica do servidor; é informativa e pode vir em inglês.
	// Para mostrar ao usuário, prefira MessagePt e Fields.
	Message string
	// MessagePt é a explicação em pt-BR (details.message_pt), quando houver.
	MessagePt string
	// Details são os detalhes adicionais do erro, quando houver.
	Details map[string]any
	// Fields mapeia campo → mensagem nos erros de validação (details.fields).
	Fields map[string]string
	// RequestID vem do cabeçalho X-Request-Id ou do corpo.
	RequestID string
	// RetryAfter é o tempo indicado pelo servidor (429/503), quando houver.
	RetryAfter time.Duration
	// Body é o corpo bruto da resposta.
	Body []byte
}

func (e *APIError) Error() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "meshio: %d", e.StatusCode)
	if e.Code != "" {
		sb.WriteString(" " + e.Code)
	}
	if e.Message != "" {
		sb.WriteString(": " + e.Message)
	}
	if len(e.Fields) > 0 {
		keys := make([]string, 0, len(e.Fields))
		for k := range e.Fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, k+": "+e.Fields[k])
		}
		sb.WriteString(" (" + strings.Join(parts, "; ") + ")")
	}
	if e.RequestID != "" {
		sb.WriteString(" [request_id=" + e.RequestID + "]")
	}
	return sb.String()
}

// Unwrap expõe as sentinelas da categoria do erro para errors.Is.
func (e *APIError) Unwrap() []error {
	var errs []error
	switch {
	case e.StatusCode == 400:
		errs = append(errs, ErrBadRequest)
	case e.StatusCode == 401:
		errs = append(errs, ErrAuthentication)
	case e.StatusCode == 403:
		errs = append(errs, ErrPermissionDenied)
	case e.StatusCode == 404:
		errs = append(errs, ErrNotFound)
	case e.StatusCode == 409:
		errs = append(errs, ErrConflict)
		if e.Code == "IDEMPOTENCY_CONFLICT" {
			errs = append(errs, ErrIdempotencyConflict)
		}
	case e.StatusCode == 422:
		errs = append(errs, ErrValidation)
	case e.StatusCode == 429:
		errs = append(errs, ErrRateLimit)
	case e.StatusCode >= 500:
		errs = append(errs, ErrServer)
	}
	if e.Code == "POLICY_BLOCKED" {
		errs = append(errs, ErrPolicyBlocked)
	}
	return errs
}

// IsPolicyBlocked indica code == POLICY_BLOCKED.
func (e *APIError) IsPolicyBlocked() bool { return e.Code == "POLICY_BLOCKED" }

// ConnectionError indica falha de rede antes de obter resposta.
type ConnectionError struct {
	Err        error
	beforeSend bool
}

// BeforeSend indica que a falha ocorreu antes de a requisição ser enviada.
func (e *ConnectionError) BeforeSend() bool { return e.beforeSend }

func (e *ConnectionError) Error() string   { return "meshio: falha de conexão: " + e.Err.Error() }
func (e *ConnectionError) Unwrap() []error { return []error{ErrConnection, e.Err} }

// TimeoutError indica que a tentativa excedeu o tempo configurado.
type TimeoutError struct {
	Err        error
	beforeSend bool
}

func (e *TimeoutError) Error() string   { return "meshio: tempo esgotado: " + e.Err.Error() }
func (e *TimeoutError) Unwrap() []error { return []error{ErrTimeout, e.Err} }

// WebhookVerificationError indica assinatura, cabeçalho ou horário inválidos.
type WebhookVerificationError struct{ Reason string }

func (e *WebhookVerificationError) Error() string {
	return "meshio: verificação do webhook falhou: " + e.Reason
}
func (e *WebhookVerificationError) Unwrap() error { return ErrWebhookVerification }

func newAPIError(resp *http.Response, body []byte) *APIError {
	e := &APIError{
		StatusCode: resp.StatusCode,
		RequestID:  resp.Header.Get("X-Request-Id"),
		RetryAfter: parseRetryAfter(resp.Header.Get("Retry-After")),
		Body:       body,
	}
	var env struct {
		Error *struct {
			Code    string         `json:"code"`
			Message string         `json:"message"`
			Details map[string]any `json:"details"`
		} `json:"error"`
		RequestID string `json:"request_id"`
		Detail    any    `json:"detail"`
	}
	if json.Unmarshal(body, &env) == nil {
		if env.Error != nil {
			e.Code = env.Error.Code
			e.Message = env.Error.Message
			e.Details = env.Error.Details
		}
		if e.RequestID == "" {
			e.RequestID = env.RequestID
		}
		if e.Details != nil {
			if m, ok := e.Details["message_pt"].(string); ok {
				e.MessagePt = m
			}
			e.Fields = parseFields(e.Details["fields"])
			if f, ok := e.Details["field"].(string); ok && e.Fields == nil && f != "" {
				// Variante com um único campo: {"details":{"field":"limit"}}
				e.Fields = map[string]string{f: e.Message}
			}
		}
		if env.Error == nil && env.Detail != nil {
			// Formato de validação padrão do framework: {"detail":[{"loc","msg"}]}
			e.Fields = parseFields(env.Detail)
			if s, ok := env.Detail.(string); ok {
				e.Message = s
			}
		}
	}
	if e.Message == "" {
		e.Message = http.StatusText(resp.StatusCode)
	}
	return e
}

// parseFields aceita {"campo":"mensagem"}, {"campo":["m1"]} ou [{"field"|"loc", "message"|"msg"}].
func parseFields(v any) map[string]string {
	out := map[string]string{}
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			out[k] = stringify(val)
		}
	case []any:
		for _, it := range t {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			name := ""
			if f, ok := m["field"].(string); ok {
				name = f
			} else if loc, ok := m["loc"].([]any); ok {
				parts := []string{}
				for i, l := range loc {
					if i == 0 && (l == "body" || l == "query" || l == "path") && len(loc) > 1 {
						continue
					}
					parts = append(parts, fmt.Sprint(l))
				}
				name = strings.Join(parts, ".")
			}
			msg := ""
			if s, ok := m["message"].(string); ok {
				msg = s
			} else if s, ok := m["msg"].(string); ok {
				msg = s
			}
			if name != "" {
				out[name] = msg
			}
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func stringify(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []any:
		parts := make([]string, 0, len(t))
		for _, x := range t {
			parts = append(parts, stringify(x))
		}
		return strings.Join(parts, "; ")
	case map[string]any:
		if s, ok := t["message"].(string); ok {
			return s
		}
	}
	b, _ := json.Marshal(v)
	return string(b)
}
