package meshio

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func subscriptionJSON() map[string]any {
	return map[string]any{
		"status": "active", "subscription_id": "sub_1", "seats": 2, "scheduled_seats": nil, "connected_numbers": 1,
		"number_limit": map[string]any{"limit": 2, "reason": "seats"},
		"tier":         map[string]any{"min_numbers": 1, "max_numbers": 9, "unit_price_micros": "49900000"},
		"current_cycle": map[string]any{
			"id": "cyc_1", "kind": "paid", "status": "open", "period_start": "2026-10-21T04:29:07Z", "period_end": "2026-11-21T04:29:07Z",
			"seats": 2, "max_numbers": 2, "included_messages": 20000, "used_messages": 150, "held_messages": 2, "remaining_messages": 19848,
		},
		"overage":       map[string]any{"price_micros": nil, "available": false, "unavailable_reason": "price_undefined", "enabled": false, "limit_micros": nil, "used_micros": "0"},
		"sending":       map[string]any{"allowed": true, "reason": nil},
		"open_invoices": []any{},
	}
}

func invoiceJSON(id string) map[string]any {
	return map[string]any{
		"id": id, "subscription_id": "sub_1", "kind": "subscription", "status": "open", "total_micros": "99800000",
		"period_start": "2026-10-21T04:29:07Z", "period_end": "2026-11-21T04:29:07Z", "issued_at": "2026-10-21T04:29:07Z",
		"due_at": "2026-10-28T04:29:07Z", "checkout_url": "https://pay.example/inv", "lines": []any{},
	}
}

// TestBillingRoutes confere método, caminho e corpo de cada operação da cobrança.
func TestBillingRoutes(t *testing.T) {
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/export.csv"):
			w.Header().Set("Content-Type", "text/csv")
			_, _ = w.Write([]byte("id\n"))
		case strings.HasPrefix(r.URL.Path, "/v1/billing/invoices/"):
			writeJSON(w, 200, invoiceJSON("inv_1"))
		case r.URL.Path == "/v1/billing/subscription" && r.Method == http.MethodPost:
			writeJSON(w, 201, map[string]any{"subscription": subscriptionJSON(), "invoice": invoiceJSON("inv_1")})
		case strings.HasPrefix(r.URL.Path, "/v1/billing/subscription/seats/preview"):
			writeJSON(w, 200, map[string]any{"addition": map[string]any{"add": 1}, "reduction": nil})
		case r.URL.Path == "/v1/billing/subscription/seats":
			writeJSON(w, 200, map[string]any{"subscription": subscriptionJSON(), "invoice": nil})
		case strings.HasPrefix(r.URL.Path, "/v1/billing/subscription"):
			writeJSON(w, 200, subscriptionJSON())
		default:
			writeJSON(w, 200, map[string]any{"data": []any{}, "next_cursor": nil})
		}
	})
	ctx := context.Background()
	one, three := int64(1), int64(3)
	yes := true
	steps := []struct {
		method, path, body string
		fn                 func() error
	}{
		{"GET", "/v1/billing/subscription", "", func() error {
			s, err := c.Billing.Subscription.Get(ctx)
			if err == nil && (s.Status != "active" || *s.Seats != 2 || s.CurrentCycle.RemainingMessages != 19848 || s.Overage.Available) {
				t.Errorf("assinatura mal decodificada: %+v", s)
			}
			return err
		}},
		{"POST", "/v1/billing/subscription/preview", `{"seats":3}`, func() error {
			_, err := c.Billing.Subscription.Preview(ctx, &SubscriptionPreviewParams{Seats: 3})
			return err
		}},
		{"POST", "/v1/billing/subscription", `{"seats":2,"payment_method":"invoice_link"}`, func() error {
			r, err := c.Billing.Subscription.Create(ctx, &CreateSubscriptionParams{Seats: 2, PaymentMethod: "invoice_link"})
			if err == nil && (r.Invoice == nil || r.Invoice.CheckoutURL == nil) {
				t.Errorf("fatura ausente: %+v", r)
			}
			return err
		}},
		{"POST", "/v1/billing/subscription/seats/preview", `{"add":1}`, func() error {
			p, err := c.Billing.Subscription.PreviewSeats(ctx, &SeatsPreviewParams{Add: &one})
			if err == nil && (p.Addition == nil || p.Reduction != nil) {
				t.Errorf("prévia mal decodificada: %+v", p)
			}
			return err
		}},
		{"POST", "/v1/billing/subscription/seats", `{"target":3}`, func() error {
			_, err := c.Billing.Subscription.ChangeSeats(ctx, &SeatsParams{Target: &three})
			return err
		}},
		{"POST", "/v1/billing/subscription/cancel", "", func() error { _, err := c.Billing.Subscription.Cancel(ctx); return err }},
		{"POST", "/v1/billing/subscription/resume", "", func() error { _, err := c.Billing.Subscription.Resume(ctx); return err }},
		{"PATCH", "/v1/billing/subscription/overage", `{"enabled":true,"limit_micros":"50000000"}`, func() error {
			_, err := c.Billing.Subscription.UpdateOverage(ctx, &OverageParams{Enabled: true, LimitMicros: "50000000"})
			return err
		}},
		{"PATCH", "/v1/billing/subscription/payment-method", `{"payment_method":"wallet","accept_wallet_terms":true}`, func() error {
			_, err := c.Billing.Subscription.UpdatePaymentMethod(ctx, &PaymentMethodParams{PaymentMethod: "wallet", AcceptWalletTerms: &yes})
			return err
		}},
		{"GET", "/v1/billing/allowance", "", func() error { _, err := c.Billing.Allowance.Get(ctx); return err }},
		{"GET", "/v1/billing/allowance/consumptions", "", func() error {
			_, err := c.Billing.Allowance.Consumptions(ctx, &BillingAllowanceConsumptionsListParams{Direction: "inbound"})
			return err
		}},
		{"GET", "/v1/billing/cycles", "", func() error { _, err := c.Billing.Cycles.List(ctx, nil); return err }},
		{"GET", "/v1/billing/invoices", "", func() error {
			_, err := c.Billing.Invoices.List(ctx, &BillingInvoicesListParams{Status: "open"})
			return err
		}},
		{"GET", "/v1/billing/invoices/inv_1", "", func() error { _, err := c.Billing.Invoices.Get(ctx, "inv_1"); return err }},
		{"POST", "/v1/billing/invoices/inv_1/pay", `{"method":"link"}`, func() error {
			_, err := c.Billing.Invoices.Pay(ctx, "inv_1", &PayInvoiceParams{Method: "link"})
			return err
		}},
		{"GET", "/v1/billing/pricing", "", func() error { _, err := c.Billing.Pricing(ctx); return err }},
		{"GET", "/v1/billing/wallet", "", func() error { _, err := c.Billing.Wallet(ctx); return err }},
		{"GET", "/v1/billing/entries", "", func() error { _, err := c.Billing.Entries.List(ctx, nil); return err }},
		{"GET", "/v1/billing/entries/export.csv", "", func() error { _, err := c.Billing.Entries.ExportCSV(ctx, nil); return err }},
		{"GET", "/v1/billing/topups", "", func() error { _, err := c.Billing.Topups.List(ctx, nil); return err }},
		{"GET", "/v1/billing/topups/top_1", "", func() error {
			// a resposta simulada não é uma recarga; só o caminho importa aqui
			_, err := c.Billing.Topups.Get(ctx, "top_1")
			return err
		}},
	}
	for i, st := range steps {
		if err := st.fn(); err != nil {
			t.Fatalf("%s %s: %v", st.method, st.path, err)
		}
		r := (*reqs)[i]
		if r.Method != st.method || r.Path != st.path {
			t.Errorf("passo %d: %s %s, esperava %s %s", i, r.Method, r.Path, st.method, st.path)
		}
		if st.body != "" && r.Body != st.body {
			t.Errorf("%s: corpo %s, esperava %s", st.path, r.Body, st.body)
		}
	}
	if q := (*reqs)[10].Query; q != "direction=inbound" {
		t.Errorf("filtro de consumo = %q", q)
	}
}

// TestBillingIdempotentOperations: Subscription.Create, ChangeSeats e Invoices.Pay geram a
// chave e a reutilizam nas novas tentativas; as prévias não são repetidas em 5xx.
func TestBillingIdempotentOperations(t *testing.T) {
	ctx := context.Background()
	t.Run("subscription.create", func(t *testing.T) {
		n := 0
		c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			n++
			if n == 1 {
				writeJSON(w, 503, errBody("UNAVAILABLE", "unavailable", nil))
				return
			}
			w.Header().Set("Idempotent-Replayed", "true")
			writeJSON(w, 201, map[string]any{"subscription": subscriptionJSON(), "invoice": invoiceJSON("inv_1")})
		})
		var meta ResponseMeta
		if _, err := c.Billing.Subscription.Create(ctx, &CreateSubscriptionParams{Seats: 2}, WithResponseMeta(&meta)); err != nil {
			t.Fatal(err)
		}
		k1, k2 := (*reqs)[0].Header.Get("Idempotency-Key"), (*reqs)[1].Header.Get("Idempotency-Key")
		if len(*reqs) != 2 || k1 == "" || k1 != k2 || len(k1) != 36 {
			t.Fatalf("chave deve ser gerada e reutilizada: %d %q %q", len(*reqs), k1, k2)
		}
		if !meta.IdempotentReplay || meta.IdempotencyKey != k1 {
			t.Fatalf("meta: %+v", meta)
		}
	})
	t.Run("invoices.pay", func(t *testing.T) {
		n := 0
		c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			n++
			if n == 1 {
				writeJSON(w, 502, errBody("PAYMENT_PROVIDER_UNAVAILABLE", "provider unavailable", nil))
				return
			}
			writeJSON(w, 200, invoiceJSON("inv_9"))
		})
		inv, err := c.Billing.Invoices.Pay(ctx, "inv_9", nil, WithIdempotencyKey("pay-inv_9"))
		if err != nil || inv.ID != "inv_9" {
			t.Fatalf("%v %+v", err, inv)
		}
		for _, r := range *reqs {
			if r.Header.Get("Idempotency-Key") != "pay-inv_9" || r.Body != "{}" {
				t.Fatalf("requisição inesperada: %+v", r)
			}
		}
		if len(*reqs) != 2 {
			t.Fatalf("tentativas = %d", len(*reqs))
		}
	})
	t.Run("seats.change", func(t *testing.T) {
		c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 503, errBody("UNAVAILABLE", "unavailable", nil))
		})
		add := int64(1)
		_, err := c.Billing.Subscription.ChangeSeats(ctx, &SeatsParams{Add: &add})
		k := (*reqs)[0].Header.Get("Idempotency-Key")
		if !errors.Is(err, ErrServer) || len(*reqs) != 3 || k == "" || (*reqs)[2].Header.Get("Idempotency-Key") != k {
			t.Fatalf("seats deve gerar a chave e repetir com ela: %v %d %q", err, len(*reqs), k)
		}
	})
	t.Run("subscription.preview não é repetido", func(t *testing.T) {
		c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, 500, errBody("INTERNAL", "internal", nil))
		})
		_, err := c.Billing.Subscription.Preview(ctx, &SubscriptionPreviewParams{Seats: 1})
		if !errors.Is(err, ErrServer) || len(*reqs) != 1 || (*reqs)[0].Header.Get("Idempotency-Key") != "" {
			t.Fatalf("%v %d", err, len(*reqs))
		}
	})
}

func TestBillingErrors(t *testing.T) {
	cases := []struct {
		status    int
		code      string
		reason    string
		sentinels []error
		not       []error
	}{
		{402, "SUBSCRIPTION_REQUIRED", "trial_ended", []error{ErrPaymentRequired, ErrSubscriptionRequired}, []error{ErrAllowanceExhausted}},
		{402, "SUBSCRIPTION_PAST_DUE", "subscription_past_due", []error{ErrPaymentRequired, ErrSubscriptionPastDue}, nil},
		{402, "ALLOWANCE_EXHAUSTED", "overage_unavailable", []error{ErrPaymentRequired, ErrAllowanceExhausted}, []error{ErrSubscriptionRequired}},
		{409, "NUMBER_LIMIT_REACHED", "seats", []error{ErrConflict, ErrNumberLimitReached}, []error{ErrIdempotencyConflict}},
		{409, "OVERAGE_UNAVAILABLE", "price_undefined", []error{ErrConflict, ErrOverageUnavailable}, nil},
		{422, "PAYMENT_METHOD_UNAVAILABLE", "wallet_terms_not_accepted", []error{ErrValidation, ErrPaymentMethodUnavailable}, nil},
		{503, "BILLING_NOT_CONFIGURED", "", []error{ErrServer, ErrBillingNotConfigured}, nil},
		{410, "DISCONTINUED", "discontinued", []error{ErrGone}, []error{ErrNotFound}},
	}
	for _, tc := range cases {
		t.Run(tc.code, func(t *testing.T) {
			details := map[string]any{"message_pt": "Explicação em português.", "path": "/billing"}
			if tc.reason != "" {
				details["reason"] = tc.reason
			}
			if tc.code == "DISCONTINUED" {
				details["replacement"] = "/v1/billing/subscription"
			}
			c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				writeJSON(w, tc.status, errBody(tc.code, "technical message", details))
			}, WithMaxRetries(0))
			_, err := c.Messages.Send(context.Background(), &SendMessageParams{ChannelID: "chn_1", To: "+5511988887777", Type: "text"})
			for _, s := range tc.sentinels {
				if !errors.Is(err, s) {
					t.Errorf("errors.Is(%v, %v) = false", err, s)
				}
			}
			for _, s := range tc.not {
				if errors.Is(err, s) {
					t.Errorf("errors.Is(%v, %v) não deveria ser verdadeiro", err, s)
				}
			}
			var apiErr *APIError
			if !errors.As(err, &apiErr) || apiErr.Code != tc.code || apiErr.Reason != tc.reason || apiErr.MessagePt != "Explicação em português." {
				t.Fatalf("campos inesperados: %+v", apiErr)
			}
			if tc.code == "DISCONTINUED" && apiErr.Replacement != "/v1/billing/subscription" {
				t.Errorf("replacement = %q", apiErr.Replacement)
			}
			if len(*reqs) != 1 {
				t.Errorf("tentativas = %d", len(*reqs))
			}
		})
	}
}

func TestAllowanceConsumptionsIterator(t *testing.T) {
	pages := map[string]any{
		"":   map[string]any{"data": []any{map[string]any{"id": "cns_1", "direction": "outbound", "classification": "included"}}, "next_cursor": "c2"},
		"c2": map[string]any{"data": []any{map[string]any{"id": "cns_2", "direction": "inbound", "classification": "excess_unbilled"}}, "next_cursor": nil},
	}
	c, reqs, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, pages[r.URL.Query().Get("cursor")])
	})
	items, err := c.Billing.Allowance.ConsumptionsAll(context.Background(), &BillingAllowanceConsumptionsListParams{CycleID: "cyc_1", Limit: 1}).Collect()
	if err != nil || len(items) != 2 || items[1].Classification != "excess_unbilled" {
		t.Fatalf("%v %+v", err, items)
	}
	for _, r := range *reqs {
		if !strings.Contains(r.Query, "cycle_id=cyc_1") || !strings.Contains(r.Query, "limit=1") {
			t.Fatalf("filtros devem ser mantidos: %q", r.Query)
		}
	}
	invs, err := c.Billing.Invoices.ListAll(context.Background(), nil).Collect()
	if err != nil {
		t.Fatal(err)
	}
	_ = invs
}

func TestBillingEventsDecode(t *testing.T) {
	raw := []byte(`{"event_id":"evt_1","event_type":"billing.allowance_exhausted","occurred_at":"2026-10-20T12:00:00Z","business_id":null,"channel_id":null,
		"data":{"cycle_id":"cyc_1","cycle_kind":"paid","included_messages":20000,"used_messages":20000,"percent":100,"cycle_ends_at":"2026-11-21T04:29:07Z"}}`)
	ev, err := ParseWebhookEvent(raw)
	if err != nil {
		t.Fatal(err)
	}
	d, err := ev.ParsedData()
	ae, ok := d.(*BillingAllowanceExhaustedData)
	if err != nil || !ok || ae.CycleID != "cyc_1" || ae.Percent != 100 {
		t.Fatalf("billing.allowance_exhausted mal decodificado: %T %+v %v", d, d, err)
	}
	for _, et := range []string{
		EventTypeBillingAllowanceWarning, EventTypeBillingAllowanceExhausted, EventTypeBillingOverageLimitReached,
		EventTypeBillingInvoiceCreated, EventTypeBillingInvoicePaid, EventTypeBillingInvoicePastDue,
		EventTypeBillingSubscriptionUpdated, EventTypeBillingTrialEnding, EventTypeBillingTrialEnded,
	} {
		if _, ok, _ := decodeEventData(et, json.RawMessage(`{}`)); !ok {
			t.Errorf("tipo sem estrutura: %s", et)
		}
	}
	for _, old := range []string{"billing.balance_low", "billing.balance_added", "billing.charges_pending", "channel.billing_suspended", "business.allowance_reached"} {
		if _, ok, _ := decodeEventData(old, json.RawMessage(`{}`)); ok {
			t.Errorf("evento removido ainda tipado: %s", old)
		}
	}
}
