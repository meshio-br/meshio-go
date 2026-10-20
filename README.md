# meshio-go

SDK oficial da [API Meshio](https://app.meshio.com.br/docs) para Go: WhatsApp para integradores — envio e recebimento de mensagens, conexão de números, templates, webhooks, uso e cobrança.

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

**Idempotência:** `Messages.Send` sempre envia `Idempotency-Key`. Sem `WithIdempotencyKey`, o SDK gera um UUID v4 e reutiliza a mesma chave nas novas tentativas. Repetir uma chamada com a mesma chave devolve a mesma mensagem com `IdempotentReplay == true`. `Billing.Topups.Create` segue a mesma regra. Use `meshio.WithResponseMeta(&meta)` para ler status, cabeçalhos, `RequestID` e a chave usada.

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
	fmt.Println(apiErr.Fields) // campo → mensagem
case errors.As(err, &apiErr):
	fmt.Println(apiErr.StatusCode, apiErr.Code, apiErr.Message, apiErr.RequestID)
case errors.Is(err, meshio.ErrTimeout), errors.Is(err, meshio.ErrConnection):
	// rede
}
```

| Status | Sentinela |
| --- | --- |
| 400 | `ErrBadRequest` |
| 401 | `ErrAuthentication` |
| 403 | `ErrPermissionDenied` |
| 404 | `ErrNotFound` |
| 409 | `ErrConflict` (+ `ErrIdempotencyConflict` para `IDEMPOTENCY_CONFLICT`) |
| 422 | `ErrValidation` |
| 429 | `ErrRateLimit` (`apiErr.RetryAfter`) |
| 5xx | `ErrServer` |
| `code == POLICY_BLOCKED` | `ErrPolicyBlocked` (mantém a categoria do status) |

Falhas de rede, timeout, 429, 500, 502, 503 e 504 são repetidas até `WithMaxRetries` vezes (padrão 2), esperando `Retry-After` quando presente ou 0,5 s, 1 s, 2 s… (máx. 8 s) com jitter. Outros 4xx nunca são repetidos. POST, PATCH e DELETE só são repetidos em falha de rede ou 5xx quando levam `Idempotency-Key`.

## Dinheiro em micros

Valores `*_micros` são strings inteiras em milionésimos de real (R$ 1 = `"1000000"`). Não converta para `float64`:

```go
w, _ := client.Billing.Wallet(ctx)
reais, _ := meshio.MicrosToReais(w.AvailableMicros) // "150.00"
r, _ := meshio.MicrosToRat(w.AvailableMicros)        // *big.Rat para cálculos exatos
```

## Versões suportadas

Go 1.21, 1.22 e 1.23 (testadas no CI). Versão do SDK: `meshio.SDKVersion`. Segue SemVer; até a 1.0, versões menores podem trazer mudanças incompatíveis, sempre descritas no [CHANGELOG](CHANGELOG.md).

## Desenvolvimento

```sh
go test ./...                              # unitários
MESHIO_TEST_API_URL=http://localhost:49700 \
MESHIO_TEST_API_KEY=msk_test_… go test ./... # inclui integração com uma pilha local
go generate ./...                          # regenera os tipos a partir de api/ (requer python3)
```

Documentação completa da API: https://app.meshio.com.br/docs

## Licença

[MIT](LICENSE) © Meshio
