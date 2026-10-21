# Changelog

Todas as mudanças relevantes deste SDK. Formato baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/); versões seguem [SemVer](https://semver.org/lang/pt-BR/).

## [0.1.0] - 2026-10-21

### Adicionado

- Instalação direto do GitHub (`go get github.com/meshio-br/meshio-go@v0.1.0`).
- Cliente `meshio.NewClient` com opções funcionais (`WithAPIKey`, `WithBaseURL`, `WithTimeout`, `WithMaxRetries`, `WithHTTPClient`) e leitura de `MESHIO_API_KEY`/`MESHIO_BASE_URL`.
- Cobertura das 73 operações do contrato público v1 (empresas, canais, sessões de cadastro, registros da Meta, links de conexão, mensagens, conversas, templates, endpoints e entregas de webhook, eventos, uso, conformidade, consentimentos, opt-outs e cobrança).
- Cobrança por assinatura de números: `Billing.Subscription` (`Get`, `Preview`, `Create`, `PreviewSeats`, `ChangeSeats`, `Cancel`, `Resume`, `UpdateOverage`, `UpdatePaymentMethod`), `Billing.Allowance` (`Get`, `Consumptions`, `ConsumptionsAll`), `Billing.Cycles`, `Billing.Invoices` (`List`, `ListAll`, `Get`, `Pay`), `Billing.Pricing`, `Billing.Wallet`, `Billing.Entries` (com `ExportCSV`) e `Billing.Topups` (somente histórico).
- Tipos gerados a partir do contrato (`models_gen.go`, `params_gen.go`) e dos esquemas de eventos (`events_gen.go`), incluindo os eventos `billing.allowance_warning`, `billing.allowance_exhausted`, `billing.overage_limit_reached`, `billing.invoice_created`, `billing.invoice_paid`, `billing.invoice_past_due`, `billing.subscription_updated`, `billing.trial_ending` e `billing.trial_ended`.
- Paginação por página (`List`) e iteração automática (`ListAll` → `*Iter`).
- Idempotência automática em `Messages.Send`, `Billing.Subscription.Create`, `Billing.Subscription.ChangeSeats` e `Billing.Invoices.Pay`, com `IdempotentReplay` e `WithResponseMeta`.
- Regra única de novas tentativas: 429 e falha de conexão antes do envio sempre; 500/502/503/504 e timeout só em GET ou com `Idempotency-Key`; 503 `BILLING_NOT_CONFIGURED` (não transitório) e outros 4xx nunca. Espera por `Retry-After` ou backoff exponencial com jitter.
- Erros tipados (`*APIError`, sentinelas por categoria, `*ConnectionError` com `BeforeSend()`, `*TimeoutError`). `APIError.MessagePt` traz a explicação em pt-BR (`details.message_pt`), `APIError.Reason` o motivo (`details.reason`) e `APIError.Replacement` a rota nova de um recurso descontinuado; `Message` é informativa (pode vir em inglês).
- Sentinelas da cobrança: `ErrPaymentRequired` (402), `ErrSubscriptionRequired`, `ErrSubscriptionPastDue`, `ErrAllowanceExhausted`, `ErrNumberLimitReached`, `ErrOverageUnavailable`, `ErrPaymentMethodUnavailable`, `ErrBillingNotConfigured` e `ErrGone` (410 `DISCONTINUED`).
- Verificação de webhooks (`VerifyWebhook`, `ConstructEvent`, `VerifyWebhookSignature`) com rotação de segredo.
- Utilitários de dinheiro `MicrosToReais` e `MicrosToRat`.
- Testes de integração opcionais contra uma pilha local, incluindo assinatura, franquia, faturas, prévias de vagas, rota descontinuada (410) e envio de ponta a ponta com número conectado pelo Simulador.
