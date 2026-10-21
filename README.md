# meshio-go

SDK oficial da [API Meshio](https://app.meshio.com.br/docs) para Go: WhatsApp para integradores — envio e recebimento de mensagens, conexão de números, templates, webhooks, uso e assinatura por número.

Somente biblioteca padrão. Go 1.21 ou superior.

## Instalação

```sh
go get github.com/meshio-br/meshio-go@v0.1.0
```

## Configuração

```go
import "github.com/meshio-br/meshio-go"

client := meshio.NewClient(
	meshio.WithAPIKey(os.Getenv("MESHIO_API_KEY")), // padrão: variável MESHIO_API_KEY
	// meshio.WithBaseURL("https://api.meshio.com.br"), // padrão; ou MESHIO_BASE_URL
	// meshio.WithTimeout(30*time.Second),              // por tentativa
	// meshio.WithMaxRetries(2),
)
```

- **Teste vs. produção:** o ambiente é definido pela chave. `msk_test_…` usa o ambiente de teste (números simulados, sem cobrança real); `msk_live_…` usa produção. O endereço é o mesmo.
- `WithBaseURL` serve para apontar para outra instância (por exemplo, uma pilha local).
- `WithHTTPClient` aceita um `*http.Client` próprio (proxy, transporte).

## Enviar mensagens

Texto (dentro da janela de 24 h de atendimento):

```go
msg, err := client.Messages.Send(ctx, &meshio.SendMessageParams{
	ChannelID: "chn_…",
	To:        "+5511988887777",
	Type:      "text",
	Text:      &meshio.TextParams{Body: "Seu pedido saiu para entrega."},
}, meshio.WithIdempotencyKey("pedido-1234-entrega"))
if err != nil {
	return err
}
fmt.Println(msg.ID, msg.Status, msg.IdempotentReplay)
```

Template (iniciado pela empresa; exige consentimento registrado ou `OptIn`):

```go
msg, err := client.Messages.Send(ctx, &meshio.SendMessageParams{
	ChannelID: "chn_…",
	To:        "+5511988887777",
	Type:      "template",
	Template: &meshio.TemplateParams{
		Name:      "order_update",
		Language:  "pt_BR",
		Variables: &meshio.TemplateVariablesParams{Body: []string{"Ana", "1234"}},
	},
	OptIn: &meshio.OptInParams{Source: "checkout", CapturedAt: time.Now()},
})
```

**Idempotência:** `Messages.Send` sempre envia `Idempotency-Key`. Sem `WithIdempotencyKey`, o SDK gera um UUID v4 e reutiliza a mesma chave nas novas tentativas. Repetir uma chamada com a mesma chave devolve a mesma mensagem com `IdempotentReplay == true`. `Billing.Subscription.Create`, `Billing.Subscription.ChangeSeats` e `Billing.Invoices.Pay` seguem a mesma regra (mesma chave, sem cobrança em dobro). Use `meshio.WithResponseMeta(&meta)` para ler status, cabeçalhos, `RequestID` e a chave usada.

## Listar com paginação

```go
// uma página
page, err := client.Messages.List(ctx, &meshio.MessagesListParams{ChannelID: "chn_…", Limit: 50})
// page.Data, page.NextCursor

// todos os itens, buscando as páginas sob demanda
it := client.Businesses.ListAll(ctx, &meshio.BusinessesListParams{Status: "active", Limit: 200})
for it.Next() {
	b := it.Current()
	fmt.Println(b.ID, b.Name)
}
if err := it.Err(); err != nil {
	return err
}
```

`limit` aceita de 1 a 200. Todo `List` paginado tem um `ListAll` (ou `HistoryAll`/`TimelineAll`).

## Conectar um número (link de conexão)

O dono da empresa abre uma página hospedada pela Meshio e conecta o próprio número:

```go
link, err := client.ConnectLinks.Create(ctx, &meshio.CreateConnectLinkParams{BusinessID: "biz_…"})
if err != nil {
	return err
}
fmt.Println("Envie este link ao cliente:", *link.URL)
// Quando o número for conectado, chega o evento channel.connected.
```

## Recursos

| Serviço | Operações |
| --- | --- |
| `Version` | `Get` |
| `Businesses` | `List`, `ListAll`, `Get` |
| `Channels` | `List`, `ListAll`, `Get`, `History`, `HistoryAll`, `HealthCheck`, `Disconnect`; `SignupSessions` (`Create`, `Get`, `Complete`); `MetaLogs` |
| `ConnectLinks` | criar, listar, consultar e revogar links de conexão |
| `Messages`, `Conversations`, `Templates` | envio, consulta, conversas e linha do tempo, templates |
| `WebhookEndpoints`, `WebhookDeliveries`, `Events` | endpoints, entregas (com reenvio) e eventos |
| `Usage`, `Compliance`, `Consents`, `OptOuts` | uso e exportação CSV, visão geral, saúde da integração, conformidade, consentimentos e opt-outs |
| `Billing` | `Pricing`, `Wallet` |
| `Billing.Subscription` | `Get`, `Preview`, `Create`, `PreviewSeats`, `ChangeSeats`, `Cancel`, `Resume`, `UpdateOverage`, `UpdatePaymentMethod` |
| `Billing.Allowance` | `Get`, `Consumptions`, `ConsumptionsAll` |
| `Billing.Cycles` | `List`, `ListAll` |
| `Billing.Invoices` | `List`, `ListAll`, `Get`, `Pay` |
| `Billing.Entries` | `List`, `ListAll`, `ExportCSV` (extrato do saldo anterior) |
| `Billing.Topups` | `List`, `ListAll`, `Get` (somente histórico; não há novas recargas) |

## Assinatura e franquia

A Meshio cobra uma assinatura mensal por número contratado (vaga). Cada vaga inclui mensagens no ciclo, e a soma forma uma franquia compartilhada entre todos os números do integrador:

```go
sub, err := client.Billing.Subscription.Get(ctx)
if err != nil {
	return err
}
fmt.Println(sub.Status, *sub.Seats, "vagas,", sub.ConnectedNumbers, "números conectados")
if sub.CurrentCycle != nil {
	fmt.Println("ciclo até", sub.CurrentCycle.PeriodEnd.Format("02/01/2006"))
}

al, err := client.Billing.Allowance.Get(ctx)
if err != nil {
	return err
}
fmt.Printf("franquia: %d de %d mensagens usadas (%d%%), %d restantes\n",
	al.UsedMessages, al.IncludedMessages, al.Percent, al.RemainingMessages)

// prévia antes de acrescentar uma vaga (não altera nada)
add := int64(1)
pv, err := client.Billing.Subscription.PreviewSeats(ctx, &meshio.SeatsPreviewParams{Add: &add})
```

Com a franquia esgotada (e sem excedente autorizado), o envio é recusado com `402 ALLOWANCE_EXHAUSTED`. Nada é criado: repita depois com a mesma `Idempotency-Key`.

```go
_, err := client.Messages.Send(ctx, params, meshio.WithIdempotencyKey("pedido-1234"))
var apiErr *meshio.APIError
switch {
case errors.Is(err, meshio.ErrAllowanceExhausted) && errors.As(err, &apiErr):
	// apiErr.Reason: allowance_exhausted, trial_allowance_exhausted, overage_limit_reached…
	fmt.Println(apiErr.MessagePt) // "A franquia de mensagens deste ciclo acabou."
case errors.Is(err, meshio.ErrSubscriptionRequired), errors.Is(err, meshio.ErrSubscriptionPastDue):
	// sem assinatura ativa ou fatura em atraso: veja client.Billing.Subscription.Get
}
```

Contratar (`Billing.Subscription.Create`) e pagar faturas (`Billing.Invoices.Pay`) devolvem a fatura com `CheckoutURL` (Pix ou cartão). Mensagens recebidas nunca são bloqueadas.

## Webhooks

Crie o endpoint (o segredo `whsec_…` só aparece na criação e na rotação):

```go
ep, err := client.WebhookEndpoints.Create(ctx, &meshio.CreateWebhookEndpointParams{
	URL:        "https://seu-sistema.com.br/webhooks/meshio",
	EventTypes: []string{"message.received", "message.status_changed", "channel.connected"},
})
// guarde ep.Secret com segurança
```

Servidor mínimo que verifica a assinatura:

```go
package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/meshio-br/meshio-go"
)

func main() {
	secret := os.Getenv("MESHIO_WEBHOOK_SECRET")
	http.HandleFunc("/webhooks/meshio", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body) // corpo bruto, sem alterar
		if err != nil {
			http.Error(w, "erro de leitura", http.StatusBadRequest)
			return
		}
		event, err := meshio.VerifyWebhook(body, r.Header, secret)
		if errors.Is(err, meshio.ErrWebhookVerification) {
			http.Error(w, "assinatura inválida", http.StatusBadRequest)
			return
		}
		// entrega pelo menos uma vez: deduplique por event.EventID
		data, _ := event.ParsedData()
		switch d := data.(type) {
		case *meshio.MessageReceivedData:
			if d.Direction == "inbound" && d.Origin == "customer" {
				log.Printf("mensagem de %s", d.From)
			}
		case *meshio.MessageStatusChangedData:
			log.Printf("mensagem %s: %s", d.MessageID, d.Status)
		default:
			log.Printf("evento %s", event.EventType) // tipos novos não quebram
		}
		w.WriteHeader(http.StatusNoContent)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

A verificação usa HMAC-SHA256 sobre `"{Meshio-Timestamp}.{corpo}"`, aceita várias assinaturas `v1=` (rotação de segredo), compara em tempo constante e recusa horários com diferença maior que 5 minutos (`meshio.WithTolerance` altera).

## Erros e novas tentativas

```go
_, err := client.Messages.Send(ctx, params)
var apiErr *meshio.APIError
switch {
case errors.Is(err, meshio.ErrPolicyBlocked):
	// envio recusado por política (opt-out, consentimento, limites…)
case errors.Is(err, meshio.ErrValidation) && errors.As(err, &apiErr):
	fmt.Println(apiErr.Fields) // campo → mensagem em pt-BR
case errors.As(err, &apiErr):
	fmt.Println(apiErr.StatusCode, apiErr.Code, apiErr.MessagePt, apiErr.RequestID)
case errors.Is(err, meshio.ErrTimeout), errors.Is(err, meshio.ErrConnection):
	// rede
}
```

| Status | Sentinela |
| --- | --- |
| 400 | `ErrBadRequest` |
| 401 | `ErrAuthentication` |
| 402 | `ErrPaymentRequired` (+ `ErrSubscriptionRequired`, `ErrSubscriptionPastDue`, `ErrAllowanceExhausted`) |
| 403 | `ErrPermissionDenied` |
| 404 | `ErrNotFound` |
| 409 | `ErrConflict` (+ `ErrIdempotencyConflict`, `ErrNumberLimitReached`, `ErrOverageUnavailable`) |
| 410 | `ErrGone` (`DISCONTINUED`; `apiErr.Replacement` indica a rota nova) |
| 422 | `ErrValidation` (+ `ErrPaymentMethodUnavailable`) |
| 429 | `ErrRateLimit` (`apiErr.RetryAfter`) |
| 5xx | `ErrServer` (+ `ErrBillingNotConfigured` para `503 BILLING_NOT_CONFIGURED`) |
| `code == POLICY_BLOCKED` | `ErrPolicyBlocked` (mantém a categoria do status) |

Regra de novas tentativas (até `WithMaxRetries`, padrão 2):

- **429** é sempre repetido, esperando `Retry-After` quando presente.
- **Falha de conexão antes do envio** (DNS, conexão recusada) é sempre repetida, em qualquer método.
- **500, 502, 503, 504 e timeout** só são repetidos em GET ou quando a requisição leva `Idempotency-Key` (`Messages.Send`, `Billing.Subscription.Create`, `Billing.Subscription.ChangeSeats`, `Billing.Invoices.Pay`). Os demais POST, PATCH e DELETE não são repetidos nesses casos.
- **503 `BILLING_NOT_CONFIGURED`** nunca é repetido, porque não é transitório (a cobrança precisa ser configurada).
- **Outros 4xx** nunca são repetidos.

Sem `Retry-After`, a espera é 0,5 s, 1 s, 2 s… (máx. 8 s) com jitter.

`apiErr.Message` é a mensagem técnica do servidor, informativa e às vezes em inglês. Para o usuário final, use `apiErr.MessagePt` (de `details.message_pt`) e, em erros de validação, `apiErr.Fields` (de `details.fields`, sempre presente nesses erros). `apiErr.Reason` traz o motivo detalhado (`details.reason`).

## Dinheiro em micros

Valores `*_micros` são strings inteiras em milionésimos de real (R$ 1 = `"1000000"`). Não converta para `float64`:

```go
sub, _ := client.Billing.Subscription.Get(ctx)
reais, _ := meshio.MicrosToReais(*sub.MonthlyTotalMicros) // "99.80"
r, _ := meshio.MicrosToRat(*sub.MonthlyTotalMicros)        // *big.Rat para cálculos exatos
```

## Versões suportadas

Go 1.21, 1.22 e 1.23 (testadas no CI). Versão do SDK: `meshio.SDKVersion`. Segue SemVer; até a 1.0, versões menores podem trazer mudanças incompatíveis, sempre descritas no [CHANGELOG](CHANGELOG.md).

## Desenvolvimento

```sh
go test ./...                              # unitários
MESHIO_TEST_API_URL=http://localhost:49800 \
MESHIO_TEST_SIM_URL=http://localhost:49890 \
MESHIO_TEST_API_KEY=msk_test_… go test ./... # inclui integração com uma pilha local
go generate ./...                          # regenera os tipos a partir de api/ (requer python3)
```

Documentação completa da API: https://app.meshio.com.br/docs

## Licença

[MIT](LICENSE) © Meshio
