# Changelog

Todas as mudanças relevantes deste SDK. Formato baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/); versões seguem [SemVer](https://semver.org/lang/pt-BR/).

## [0.1.0] - 2026-10-20

### Adicionado
- Cliente `meshio.NewClient` com opções funcionais (`WithAPIKey`, `WithBaseURL`, `WithTimeout`, `WithMaxRetries`, `WithHTTPClient`) e leitura de `MESHIO_API_KEY`/`MESHIO_BASE_URL`.
- Cobertura das 67 operações do contrato público v1 (empresas, canais, sessões de cadastro, registros da Meta, links de conexão, mensagens, conversas, templates, endpoints e entregas de webhook, eventos, uso, conformidade, consentimentos, opt-outs e cobrança).
- Tipos gerados a partir do contrato (`models_gen.go`, `params_gen.go`) e dos esquemas de eventos (`events_gen.go`).
- Paginação por página (`List`) e iteração automática (`ListAll` → `*Iter`).
- Idempotência automática em `Messages.Send` e `Billing.Topups.Create`, com `IdempotentReplay` e `WithResponseMeta`.
- Novas tentativas com `Retry-After` ou backoff exponencial com jitter.
- Erros tipados (`*APIError`, sentinelas por categoria, `*ConnectionError`, `*TimeoutError`).
- Verificação de webhooks (`VerifyWebhook`, `ConstructEvent`, `VerifyWebhookSignature`) com rotação de segredo.
- Utilitários de dinheiro `MicrosToReais` e `MicrosToRat`.
