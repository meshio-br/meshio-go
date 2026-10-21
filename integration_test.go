package meshio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

// Testes de integração contra uma pilha local. Pulados sem MESHIO_TEST_API_URL.
// A pilha é compartilhada: os testes criam e apagam o que usam e não dependem de contagens globais.

func integrationClient(t *testing.T) *Client {
	t.Helper()
	base := os.Getenv("MESHIO_TEST_API_URL")
	if base == "" {
		t.Skip("MESHIO_TEST_API_URL não definido")
	}
	key := os.Getenv("MESHIO_TEST_API_KEY")
	if key == "" {
		t.Skip("MESHIO_TEST_API_KEY não definido")
	}
	return NewClient(WithAPIKey(key), WithBaseURL(base), WithTimeout(20*time.Second))
}

func ctxT(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func uniq() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }

func TestIntegrationReadOnly(t *testing.T) {
	c := integrationClient(t)
	ctx := ctxT(t)

	v, err := c.Version.Get(ctx)
	if err != nil || v.Version == nil {
		t.Fatalf("version: %v", err)
	}

	pg, err := c.Businesses.List(ctx, &BusinessesListParams{Limit: 2})
	if err != nil {
		t.Fatalf("businesses.list: %v", err)
	}
	if len(pg.Data) > 0 {
		b, err := c.Businesses.Get(ctx, pg.Data[0].ID)
		if err != nil || b.ID != pg.Data[0].ID {
			t.Fatalf("businesses.get: %v", err)
		}
	}
	// iteração: lê no máximo 5 itens em páginas de 2
	it := c.Businesses.ListAll(ctx, &BusinessesListParams{Limit: 2})
	seen := map[string]bool{}
	for i := 0; i < 5 && it.Next(); i++ {
		id := it.Current().ID
		if seen[id] {
			t.Fatalf("item repetido na paginação: %s", id)
		}
		seen[id] = true
	}
	if err := it.Err(); err != nil {
		t.Fatalf("iter: %v", err)
	}

	chs, err := c.Channels.List(ctx, &ChannelsListParams{Limit: 5})
	if err != nil {
		t.Fatalf("channels.list: %v", err)
	}
	if len(chs.Data) > 0 {
		ch := chs.Data[0]
		if _, err := c.Channels.Get(ctx, ch.ID); err != nil {
			t.Fatalf("channels.get: %v", err)
		}
		if _, err := c.Channels.History(ctx, ch.ID, &ChannelsHistoryParams{Limit: 5}); err != nil {
			t.Fatalf("channels.history: %v", err)
		}
		if _, err := c.Channels.MetaLogs.List(ctx, ch.ID, &MetaLogsListParams{Limit: 5}); err != nil {
			t.Fatalf("meta_logs.list: %v", err)
		}
	}

	checks := []struct {
		name string
		fn   func() error
	}{
		{"messages.list", func() error { _, err := c.Messages.List(ctx, &MessagesListParams{Limit: 5}); return err }},
		{"conversations.list", func() error { _, err := c.Conversations.List(ctx, &ConversationsListParams{Limit: 5}); return err }},
		{"templates.list", func() error { _, err := c.Templates.List(ctx, &TemplatesListParams{Limit: 5}); return err }},
		{"connect_links.list", func() error { _, err := c.ConnectLinks.List(ctx, &ConnectLinksListParams{Limit: 5}); return err }},
		{"webhook_endpoints.list", func() error {
			_, err := c.WebhookEndpoints.List(ctx, &WebhookEndpointsListParams{Limit: 5})
			return err
		}},
		{"webhook_deliveries.list", func() error {
			_, err := c.WebhookDeliveries.List(ctx, &WebhookDeliveriesListParams{Limit: 5})
			return err
		}},
		{"events.list", func() error { _, err := c.Events.List(ctx, &EventsListParams{Limit: 5}); return err }},
		{"usage.get", func() error { _, err := c.Usage.Get(ctx, nil); return err }},
		{"usage.export_csv", func() error { _, err := c.Usage.ExportCSV(ctx, nil); return err }},
		{"overview.get", func() error { _, err := c.Usage.Overview(ctx, nil); return err }},
		{"status.integration", func() error { _, err := c.Usage.IntegrationStatus(ctx); return err }},
		{"compliance.status", func() error { _, err := c.Compliance.Status(ctx); return err }},
		{"consents.list", func() error { _, err := c.Consents.List(ctx, &ConsentsListParams{Limit: 5}); return err }},
		{"opt_outs.list", func() error { _, err := c.OptOuts.List(ctx, &OptOutsListParams{Limit: 5}); return err }},
		{"billing.wallet", func() error {
			w, err := c.Billing.Wallet(ctx)
			if err == nil {
				if _, err := MicrosToReais(w.PaidBalanceMicros); err != nil {
					return err
				}
			}
			return err
		}},
		{"billing.pricing", func() error { _, err := c.Billing.Pricing(ctx); return err }},
		{"billing.subscription.get", func() error { _, err := c.Billing.Subscription.Get(ctx); return err }},
		{"billing.allowance.get", func() error { _, err := c.Billing.Allowance.Get(ctx); return err }},
		{"billing.allowance.consumptions", func() error {
			_, err := c.Billing.Allowance.Consumptions(ctx, &BillingAllowanceConsumptionsListParams{Limit: 5})
			return err
		}},
		{"billing.cycles.list", func() error { _, err := c.Billing.Cycles.List(ctx, &BillingCyclesListParams{Limit: 5}); return err }},
		{"billing.invoices.list", func() error {
			_, err := c.Billing.Invoices.List(ctx, &BillingInvoicesListParams{Limit: 5})
			return err
		}},
		{"billing.entries.list", func() error { _, err := c.Billing.Entries.List(ctx, &BillingEntriesListParams{Limit: 5}); return err }},
		{"billing.entries.export_csv", func() error { _, err := c.Billing.Entries.ExportCSV(ctx, nil); return err }},
		{"billing.topups.list", func() error { _, err := c.Billing.Topups.List(ctx, &TopupsListParams{Limit: 5}); return err }},
	}
	for _, ck := range checks {
		if err := ck.fn(); err != nil {
			t.Errorf("%s: %v", ck.name, err)
		}
	}
}

func TestIntegrationErrors(t *testing.T) {
	c := integrationClient(t)
	ctx := ctxT(t)

	_, err := c.Businesses.Get(ctx, "biz_doesnotexist0000000000")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || !errors.Is(err, ErrNotFound) || apiErr.RequestID == "" {
		t.Fatalf("esperava 404 tipado com request id, veio %v", err)
	}

	_, err = c.Businesses.List(ctx, &BusinessesListParams{Limit: 999})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("esperava 422, veio %v", err)
	}

	bad := NewClient(WithAPIKey("msk_test_invalid"), WithBaseURL(os.Getenv("MESHIO_TEST_API_URL")))
	if _, err := bad.Version.Get(ctx); err != nil && !errors.Is(err, ErrAuthentication) {
		t.Fatalf("version com chave inválida: %v", err)
	}
	if _, err := bad.Businesses.List(ctx, nil); !errors.Is(err, ErrAuthentication) {
		t.Fatalf("esperava 401, veio %v", err)
	}
}

func TestIntegrationWebhookEndpointLifecycle(t *testing.T) {
	c := integrationClient(t)
	ctx := ctxT(t)

	created, err := c.WebhookEndpoints.Create(ctx, &CreateWebhookEndpointParams{
		URL:         "https://example.com/meshio-go-sdk/" + uniq(),
		EventTypes:  []string{"webhook.test", "message.received"},
		Description: "meshio-go integration test",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	t.Cleanup(func() { _ = c.WebhookEndpoints.Delete(context.Background(), created.ID) })
	if !strings.HasPrefix(created.Secret, "whsec_") {
		t.Fatalf("segredo inesperado: %q", created.Secret)
	}
	if _, err := c.WebhookEndpoints.Get(ctx, created.ID); err != nil {
		t.Fatalf("get: %v", err)
	}
	upd, err := c.WebhookEndpoints.Update(ctx, created.ID, &UpdateWebhookEndpointParams{Description: "atualizado"})
	if err != nil || upd.Description == nil || *upd.Description != "atualizado" {
		t.Fatalf("update: %v %+v", err, upd)
	}
	if _, err := c.WebhookEndpoints.Test(ctx, created.ID); err != nil {
		t.Errorf("test: %v", err)
	}
	rot, err := c.WebhookEndpoints.RotateSecret(ctx, created.ID)
	if err != nil || rot.Secret == created.Secret {
		t.Errorf("rotate: %v", err)
	}
	dl, err := c.WebhookDeliveries.List(ctx, &WebhookDeliveriesListParams{EndpointID: created.ID, Limit: 5})
	if err != nil {
		t.Errorf("deliveries.list: %v", err)
	} else if len(dl.Data) > 0 {
		if _, err := c.WebhookDeliveries.Get(ctx, dl.Data[0].ID); err != nil {
			t.Errorf("deliveries.get: %v", err)
		}
		if _, err := c.WebhookDeliveries.Replay(ctx, dl.Data[0].ID); err != nil {
			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Errorf("deliveries.replay: %v", err)
			}
		}
	}
	if err := c.WebhookEndpoints.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := c.WebhookEndpoints.Get(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("endpoint apagado deveria dar 404: %v", err)
	}
}

func firstBusiness(t *testing.T, c *Client, ctx context.Context) string {
	pg, err := c.Businesses.List(ctx, &BusinessesListParams{Limit: 1, Status: "active"})
	if err != nil {
		t.Fatalf("businesses.list: %v", err)
	}
	if len(pg.Data) == 0 {
		t.Skip("nenhuma empresa ativa na pilha")
	}
	return pg.Data[0].ID
}

func TestIntegrationConnectLinkLifecycle(t *testing.T) {
	c := integrationClient(t)
	ctx := ctxT(t)
	biz := firstBusiness(t, c, ctx)
	link, err := c.ConnectLinks.Create(ctx, &CreateConnectLinkParams{BusinessID: biz})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	t.Cleanup(func() { _, _ = c.ConnectLinks.Revoke(context.Background(), link.ID) })
	if link.URL == nil || *link.URL == "" {
		t.Errorf("link sem URL: %+v", link)
	}
	if _, err := c.ConnectLinks.Get(ctx, link.ID); err != nil {
		t.Fatalf("get: %v", err)
	}
	rv, err := c.ConnectLinks.Revoke(ctx, link.ID)
	if err != nil || rv.Status != "revoked" {
		t.Fatalf("revoke: %v %+v", err, rv)
	}
}

func TestIntegrationOptOutAndConsent(t *testing.T) {
	c := integrationClient(t)
	ctx := ctxT(t)
	biz := firstBusiness(t, c, ctx)
	phone := fmt.Sprintf("+55119%08d", time.Now().UnixNano()%100000000)

	oo, err := c.OptOuts.Create(ctx, &OptOutParams{BusinessID: biz, Phone: phone})
	if err != nil {
		t.Fatalf("opt_outs.create: %v", err)
	}
	_ = oo
	if err := c.OptOuts.Delete(ctx, &OptOutsDeleteParams{BusinessID: biz, Phone: phone}); err != nil {
		t.Fatalf("opt_outs.delete: %v", err)
	}

	cb, err := c.Consents.Create(ctx, &ConsentBulkParams{Consents: []ConsentParams{{
		BusinessID: biz, Phone: phone, Source: "other", CapturedAt: time.Now().UTC().Truncate(time.Second),
	}}})
	if err != nil {
		t.Fatalf("consents.create: %v", err)
	}
	list, err := c.Consents.List(ctx, &ConsentsListParams{BusinessID: biz, Phone: phone})
	if err != nil {
		t.Fatalf("consents.list: %v", err)
	}
	for _, cs := range list.Data {
		if _, err := c.Consents.Delete(ctx, cs.ID); err != nil {
			t.Errorf("consents.delete: %v", err)
		}
	}
	_ = cb
}

func TestIntegrationSendMessage(t *testing.T) {
	c := integrationClient(t)
	ctx := ctxT(t)
	chs, err := c.Channels.List(ctx, &ChannelsListParams{Limit: 20, Status: "active"})
	if err != nil {
		t.Fatalf("channels.list: %v", err)
	}
	if len(chs.Data) == 0 {
		// sem canal ativo, tenta qualquer canal: a recusa precisa vir como erro tipado
		chs, err = c.Channels.List(ctx, &ChannelsListParams{Limit: 1})
		if err != nil {
			t.Fatalf("channels.list: %v", err)
		}
		if len(chs.Data) == 0 {
			t.Skip("nenhum canal na pilha")
		}
	}
	key := "meshio-go-" + uniq()
	params := &SendMessageParams{ChannelID: chs.Data[0].ID, To: "+5511988887777", Type: "text", Text: &TextParams{Body: "Teste do SDK Go"}}
	msg, err := c.Messages.Send(ctx, params, WithIdempotencyKey(key))
	if err != nil {
		// envio pode ser recusado por política/canal; o erro precisa ser tipado
		var apiErr *APIError
		if !errors.As(err, &apiErr) || apiErr.Code == "" {
			t.Fatalf("erro não tipado: %v", err)
		}
		t.Logf("envio recusado (aceito): %v", err)
		return
	}
	again, err := c.Messages.Send(ctx, params, WithIdempotencyKey(key))
	if err != nil {
		t.Fatalf("repetição: %v", err)
	}
	if again.ID != msg.ID || !again.IdempotentReplay {
		t.Errorf("repetição deveria devolver a mesma mensagem com replay: %s %s %v", msg.ID, again.ID, again.IdempotentReplay)
	}
	if _, err := c.Messages.Get(ctx, msg.ID); err != nil {
		t.Errorf("messages.get: %v", err)
	}
}

// TestIntegrationContractDrift decodifica respostas reais de forma estrita e registra campos fora do contrato.
func TestIntegrationContractDrift(t *testing.T) {
	c := integrationClient(t)
	ctx := ctxT(t)
	targets := map[string]any{
		"/v1/version":                                &Version{},
		"/v1/businesses?limit=5":                     &BusinessPage{},
		"/v1/channels?limit=5":                       &ChannelPage{},
		"/v1/messages?limit=5":                       &MessagePage{},
		"/v1/conversations?limit=5":                  &ConversationPage{},
		"/v1/templates?limit=5":                      &TemplatePage{},
		"/v1/connect-links?limit=5":                  &ConnectLinkPage{},
		"/v1/webhook-endpoints?limit=5":              &WebhookEndpointPage{},
		"/v1/webhook-deliveries?limit=5":             &WebhookDeliveryPage{},
		"/v1/events?limit=5":                         &EventPage{},
		"/v1/usage":                                  &Usage{},
		"/v1/overview":                               &Overview{},
		"/v1/status/integration":                     &IntegrationStatus{},
		"/v1/compliance/status":                      &ComplianceStatus{},
		"/v1/consents?limit=5":                       &ConsentPage{},
		"/v1/opt-outs?limit=5":                       &OptOutPage{},
		"/v1/billing/wallet":                         &Wallet{},
		"/v1/billing/pricing":                        &Pricing{},
		"/v1/billing/subscription":                   &Subscription{},
		"/v1/billing/allowance":                      &Allowance{},
		"/v1/billing/allowance/consumptions?limit=5": &ConsumptionPage{},
		"/v1/billing/cycles?limit=5":                 &CyclePage{},
		"/v1/billing/invoices?limit=5":               &InvoicePage{},
		"/v1/billing/entries?limit=5":                &EntryPage{},
		"/v1/billing/topups?limit=5":                 &TopupPage{},
	}
	for path, target := range targets {
		req, _ := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Logf("DRIFT %s: status %d", path, resp.StatusCode)
			continue
		}
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()
		if err := dec.Decode(target); err != nil {
			t.Logf("DRIFT %s: %v", path, err)
		}
	}
}

// TestIntegrationBilling lê assinatura, franquia, ciclos e faturas e calcula prévias,
// sem contratar nem pagar (o integrador da pilha já tem assinatura ativa).
func TestIntegrationBilling(t *testing.T) {
	c := integrationClient(t)
	ctx := ctxT(t)

	sub, err := c.Billing.Subscription.Get(ctx)
	if err != nil {
		t.Fatalf("subscription.get: %v", err)
	}
	if sub.Status != "active" || sub.Seats == nil || *sub.Seats < 1 || sub.CurrentCycle == nil {
		t.Fatalf("esperava assinatura ativa com vagas e ciclo: %+v", sub)
	}
	if sub.NumberLimit.Limit < *sub.Seats {
		t.Errorf("limite de números inconsistente: %+v", sub.NumberLimit)
	}
	if sub.MonthlyTotalMicros == nil {
		t.Fatalf("assinatura sem total mensal")
	}
	if _, err := MicrosToReais(*sub.MonthlyTotalMicros); err != nil {
		t.Errorf("total mensal: %v", err)
	}
	t.Logf("assinatura %s: %d vagas, %d números, ciclo %s até %s",
		sub.Status, *sub.Seats, sub.ConnectedNumbers, sub.CurrentCycle.ID, sub.CurrentCycle.PeriodEnd.Format(time.RFC3339))

	al, err := c.Billing.Allowance.Get(ctx)
	if err != nil {
		t.Fatalf("allowance.get: %v", err)
	}
	if al.CycleID == nil || *al.CycleID != sub.CurrentCycle.ID {
		t.Errorf("franquia em outro ciclo: %v x %s", al.CycleID, sub.CurrentCycle.ID)
	}
	want := al.IncludedMessages - al.UsedMessages - al.HeldMessages
	if want < 0 {
		want = 0
	}
	if al.RemainingMessages != want {
		t.Errorf("restantes = %d, esperava %d", al.RemainingMessages, want)
	}
	t.Logf("franquia: %d incluídas, %d usadas, %d restantes (%d%%)", al.IncludedMessages, al.UsedMessages, al.RemainingMessages, al.Percent)

	it := c.Billing.Allowance.ConsumptionsAll(ctx, &BillingAllowanceConsumptionsListParams{CycleID: *al.CycleID, Limit: 2})
	for i := 0; i < 3 && it.Next(); i++ {
		if cs := it.Current(); cs.CycleID == nil || *cs.CycleID != *al.CycleID {
			t.Errorf("evidência fora do ciclo filtrado: %+v", cs)
		}
	}
	if err := it.Err(); err != nil {
		t.Errorf("allowance.consumptions: %v", err)
	}

	cycles, err := c.Billing.Cycles.List(ctx, &BillingCyclesListParams{Limit: 5})
	if err != nil || len(cycles.Data) == 0 {
		t.Fatalf("cycles.list: %v (%d)", err, len(cycles.Data))
	}

	invs, err := c.Billing.Invoices.List(ctx, &BillingInvoicesListParams{Limit: 5})
	if err != nil {
		t.Fatalf("invoices.list: %v", err)
	}
	if len(invs.Data) > 0 {
		inv, err := c.Billing.Invoices.Get(ctx, invs.Data[0].ID)
		if err != nil || inv.ID != invs.Data[0].ID {
			t.Fatalf("invoices.get: %v", err)
		}
	}
	if _, err := c.Billing.Invoices.Get(ctx, "inv_doesnotexist0000000000"); !errors.Is(err, ErrNotFound) {
		t.Errorf("fatura inexistente deveria dar 404: %v", err)
	}

	pv, err := c.Billing.Subscription.Preview(ctx, &SubscriptionPreviewParams{Seats: 10})
	if err != nil {
		t.Fatalf("subscription.preview: %v", err)
	}
	if pv.Seats != 10 || pv.IncludedMessages != 10*pv.MessagesPerNumber {
		t.Errorf("prévia inconsistente: %+v", pv)
	}

	add := int64(1)
	sp, err := c.Billing.Subscription.PreviewSeats(ctx, &SeatsPreviewParams{Add: &add})
	if err != nil {
		t.Fatalf("subscription.seats_preview: %v", err)
	}
	if sp.Addition == nil {
		t.Errorf("prévia de acréscimo vazia: %+v", sp)
	}

	// a prévia não altera a assinatura
	after, err := c.Billing.Subscription.Get(ctx)
	if err != nil || *after.Seats != *sub.Seats {
		t.Errorf("prévia não deveria alterar vagas: %v", err)
	}

	if _, err := c.Billing.Subscription.Preview(ctx, &SubscriptionPreviewParams{Seats: 0}); !errors.Is(err, ErrValidation) {
		t.Errorf("0 vagas deveria dar 422: %v", err)
	}
}

// TestIntegrationDiscontinuedRoute confere o 410 DISCONTINUED de uma rota do modelo anterior.
func TestIntegrationDiscontinuedRoute(t *testing.T) {
	c := integrationClient(t)
	ctx := ctxT(t)

	// HTTP cru: o SDK não expõe mais as rotas descontinuadas.
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/billing/trial", nil)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusGone {
		t.Fatalf("esperava 410, veio %d: %s", resp.StatusCode, body)
	}

	// o mesmo corpo, interpretado pelo SDK
	apiErr := newAPIError(resp, body)
	if !errors.Is(apiErr, ErrGone) || apiErr.Code != "DISCONTINUED" || apiErr.Reason != "discontinued" ||
		apiErr.MessagePt == "" || apiErr.Replacement == "" || apiErr.RequestID == "" {
		t.Fatalf("410 mal interpretado: %+v", apiErr)
	}

	// um POST descontinuado, pelo caminho normal do cliente: não é repetido
	err = c.do(ctx, call{method: http.MethodPost, path: "/v1/billing/topups", body: map[string]any{}}, nil, nil)
	if !errors.Is(err, ErrGone) {
		t.Fatalf("POST /v1/billing/topups: esperava 410, veio %v", err)
	}
	t.Logf("410 DISCONTINUED → %s", apiErr.Replacement)
}

// simPost envia JSON ao Simulador da Meta e decodifica a resposta.
func simPost(t *testing.T, ctx context.Context, url string, body any, out any) {
	t.Helper()
	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("simulador: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("simulador %s: %v", url, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		t.Fatalf("simulador %s: %d %s", url, resp.StatusCode, raw)
	}
	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			t.Fatalf("simulador %s: %v", url, err)
		}
	}
}

// skipIfNumberLimit encerra o teste quando a conta já ocupou todas as vagas: a pilha
// compartilhada pode ter mais números conectados que vagas contratadas.
func skipIfNumberLimit(t *testing.T, err error) {
	t.Helper()
	var apiErr *APIError
	if errors.Is(err, ErrNumberLimitReached) && errors.As(err, &apiErr) {
		if apiErr.StatusCode != http.StatusConflict || apiErr.MessagePt == "" {
			t.Fatalf("NUMBER_LIMIT_REACHED mal formado: %+v", apiErr)
		}
		t.Skipf("409 NUMBER_LIMIT_REACHED (%s): %s", apiErr.Reason, apiErr.MessagePt)
	}
}

// TestIntegrationConnectNumberAndSend conecta um número pelo Simulador, simula uma
// mensagem recebida (abre a janela de 24 h) e envia texto livre com idempotência.
// verifiedBusiness: envio exige cliente verificado (CNPJ e atividade); usa um já verificado.
func verifiedBusiness(t *testing.T, c *Client, ctx context.Context) string {
	t.Helper()
	status, err := c.Compliance.Status(ctx)
	if err != nil {
		t.Fatalf("compliance.status: %v", err)
	}
	for _, b := range status.Businesses {
		if b.VerificationStatus == "verified" && !b.Suspended {
			return b.BusinessID
		}
	}
	t.Fatalf("nenhum cliente verificado na pilha")
	return ""
}

func TestIntegrationConnectNumberAndSend(t *testing.T) {
	c := integrationClient(t)
	sim := strings.TrimRight(os.Getenv("MESHIO_TEST_SIM_URL"), "/")
	if sim == "" {
		t.Skip("MESHIO_TEST_SIM_URL não definido")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	biz := verifiedBusiness(t, c, ctx)
	session, err := c.Channels.SignupSessions.Create(ctx, &CreateSignupSessionParams{BusinessID: biz, Mode: "cloud_api"})
	skipIfNumberLimit(t, err)
	if err != nil {
		t.Fatalf("signup_sessions.create: %v", err)
	}
	if got, err := c.Channels.SignupSessions.Get(ctx, session.SignupSessionID); err != nil || got.ID != session.SignupSessionID {
		t.Fatalf("signup_sessions.get: %v", err)
	}
	if session.Launch.URL == nil {
		t.Fatalf("sessão sem launch.url")
	}
	launch, err := url.Parse(*session.Launch.URL)
	if err != nil {
		t.Fatalf("launch.url: %v", err)
	}
	q := launch.Query()

	// O Simulador faz o papel da janela de cadastro da Meta.
	phone := fmt.Sprintf("+55 11 9%04d-%04d", 1000+mrand.Intn(9000), 1000+mrand.Intn(9000))
	var submitted struct {
		Assets struct {
			Code          string `json:"code"`
			WABAID        string `json:"waba_id"`
			PhoneNumberID string `json:"phone_number_id"`
		} `json:"assets"`
	}
	simPost(t, ctx, sim+"/es/dialog/submit", map[string]any{
		"app_id": q.Get("app_id"), "config_id": q.Get("config_id"), "state": q.Get("state"),
		"feature_type": nil, "origin": q.Get("origin"), "scenario": "success",
		"business_name": "SDK go", "display_phone_number": phone, "verified_name": "SDK go",
	}, &submitted)

	ch, err := c.Channels.SignupSessions.Complete(ctx, session.SignupSessionID, &CompleteSignupParams{
		State: session.State, Code: submitted.Assets.Code, WABAID: submitted.Assets.WABAID,
		PhoneNumberID: submitted.Assets.PhoneNumberID, Event: "FINISH",
	})
	skipIfNumberLimit(t, err)
	if err != nil {
		t.Fatalf("signup_sessions.complete: %v", err)
	}
	defer func() {
		confirm := true
		if _, err := c.Channels.Disconnect(context.Background(), ch.ID, &DisconnectParams{Confirm: &confirm}); err != nil {
			t.Errorf("channels.disconnect: %v", err)
		}
	}()
	if ch.BusinessID != biz {
		t.Fatalf("canal em outra empresa: %s", ch.BusinessID)
	}
	if _, err := c.Channels.HealthCheck(ctx, ch.ID); err != nil {
		t.Fatalf("channels.health_check: %v", err)
	}
	// A ativação é assíncrona (validação na Meta); aguarda até 20 s.
	for i := 0; i < 40 && ch.Status != "active"; i++ {
		time.Sleep(500 * time.Millisecond)
		if ch, err = c.Channels.Get(ctx, ch.ID); err != nil {
			t.Fatalf("channels.get: %v", err)
		}
	}
	if ch.Status != "active" {
		t.Fatalf("canal não ativou: %s", ch.Status)
	}

	// Mensagem do cliente abre a janela de 24 h para texto livre.
	contact := fmt.Sprintf("+55119%08d", 10000000+mrand.Intn(89999999))
	simPost(t, ctx, sim+"/_sim/inbound-message", map[string]any{
		"phone_number_id": submitted.Assets.PhoneNumberID, "from": strings.TrimPrefix(contact, "+"), "text": "Oi",
	}, nil)
	for i := 0; i < 40; i++ {
		tl, err := c.Conversations.Timeline(ctx, ch.ID, contact, &ConversationsTimelineParams{Limit: 1})
		if err == nil && len(tl.Data) > 0 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	params := &SendMessageParams{ChannelID: ch.ID, To: contact, Type: "text", Text: &TextParams{Body: "Teste do SDK Go"}, ClientReference: "sdk-go"}
	key := "sdk-go-" + uniq()
	first, err := c.Messages.Send(ctx, params, WithIdempotencyKey(key))
	if err != nil {
		t.Fatalf("messages.send: %v", err)
	}
	if first.Status != "queued" {
		t.Fatalf("status inesperado: %s", first.Status)
	}
	var meta ResponseMeta
	second, err := c.Messages.Send(ctx, params, WithIdempotencyKey(key), WithResponseMeta(&meta))
	if err != nil {
		t.Fatalf("repetição: %v", err)
	}
	if second.ID != first.ID || !meta.IdempotentReplay {
		t.Fatalf("repetição deveria devolver a mesma mensagem com replay: %s %s %v", first.ID, second.ID, meta.IdempotentReplay)
	}
	other := *params
	other.ClientReference = "other"
	_, err = c.Messages.Send(ctx, &other, WithIdempotencyKey(key))
	if !errors.Is(err, ErrIdempotencyConflict) {
		t.Fatalf("esperava IDEMPOTENCY_CONFLICT, veio %v", err)
	}
	if got, err := c.Messages.Get(ctx, first.ID); err != nil || got.ID != first.ID {
		t.Fatalf("messages.get: %v", err)
	}
	t.Logf("mensagem %s aceita (%s); repetição idempotente e conflito confirmados", first.ID, first.Status)
}
