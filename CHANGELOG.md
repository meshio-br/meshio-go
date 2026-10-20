# Changelog

Todas as mudanças relevantes deste SDK. Formato baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/); versões seguem [SemVer](https://semver.org/lang/pt-BR/).

## [Não lançado]

### Alterado
- Regra única de novas tentativas: 429 e falha de conexão antes do envio sempre; 500/502/503/504 e timeout só em GET ou com `Idempotency-Key`; outros 4xx nunca.

### Adicionado
- `APIError.MessagePt` com a explicação em pt-BR (`details.message_pt`); `Message` passa a ser documentada como informativa (pode vir em inglês).
- `ConnectionError.BeforeSend()`.
- Teste de integração de ponta a ponta: conecta um número pelo Simulador, simula mensagem recebida e envia, conferindo repetição idempotente e `IDEMPOTENCY_CONFLICT`.

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
