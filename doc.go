// Package meshio é o SDK oficial da API Meshio para Go.
//
// A API Meshio permite enviar e receber mensagens de WhatsApp em nome das
// empresas atendidas por um integrador, conectar números, gerenciar templates,
// receber eventos por webhook e acompanhar uso e cobrança.
//
// # Uso básico
//
//	client := meshio.NewClient(meshio.WithAPIKey(os.Getenv("MESHIO_API_KEY")))
//	msg, err := client.Messages.Send(ctx, &meshio.SendMessageParams{
//		ChannelID: "chn_…",
//		To:        "+5511988887777",
//		Type:      "text",
//		Text:      &meshio.TextParams{Body: "Olá!"},
//	})
//
// Todas as chamadas recebem um context.Context. O cliente é configurado com
// opções funcionais (WithAPIKey, WithBaseURL, WithTimeout, WithMaxRetries,
// WithHTTPClient) e lê MESHIO_API_KEY e MESHIO_BASE_URL do ambiente.
//
// # Paginação
//
// Métodos List devolvem uma página ({Data, NextCursor}). Métodos ListAll
// devolvem um *Iter que busca as páginas sob demanda:
//
//	it := client.Businesses.ListAll(ctx, &meshio.BusinessesListParams{Limit: 100})
//	for it.Next() {
//		fmt.Println(it.Current().Name)
//	}
//	if err := it.Err(); err != nil { … }
//
// # Erros
//
// Respostas de erro viram *APIError (status, código, mensagem, detalhes,
// request id). Use errors.Is com as sentinelas por categoria (ErrNotFound,
// ErrValidation, ErrRateLimit, ErrPolicyBlocked, ErrIdempotencyConflict…) e
// errors.As para ler os campos. Falhas de rede são *ConnectionError ou
// *TimeoutError.
//
// # Novas tentativas e idempotência
//
// Até WithMaxRetries novas tentativas (padrão 2): 429 e falha de conexão antes
// do envio são sempre repetidos; 500, 502, 503, 504 e timeout só em GET ou
// quando a requisição leva Idempotency-Key (Messages.Send e
// Billing.Topups.Create, que o SDK gera quando não informada); outros 4xx
// nunca. A espera respeita Retry-After ou usa backoff exponencial com jitter
// (0,5 s, 1 s, 2 s…, máx. 8 s).
//
// # Webhooks
//
// VerifyWebhook confere a assinatura Meshio-Signature (HMAC-SHA256) sobre o
// corpo bruto e devolve o *WebhookEvent; ParsedData decodifica o campo data no
// tipo do evento.
//
// # Dinheiro
//
// Campos *_micros são strings inteiras (R$ 1 = "1000000"). Use MicrosToReais
// ou MicrosToRat; nunca converta para float.
//
// Os tipos dos arquivos *_gen.go são gerados a partir do contrato público
// (api/openapi.public.json e api/events/).
package meshio

//go:generate sh -c "python3 internal/gen/genmodels.py && gofmt -w models_gen.go params_gen.go events_gen.go"
