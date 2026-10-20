// Código gerado por internal/gen/genmodels.py a partir do contrato. NÃO EDITE.

package meshio

import (
	"net/url"
	"strconv"
	"time"
)

// BusinessesListParams são os filtros de consulta de GET /v1/businesses (businesses_list).
type BusinessesListParams struct {
	Status string
	Q      string
	Limit  int
	Cursor string
}

func (p *BusinessesListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Status != "" {
		q.Set("status", p.Status)
	}
	if p.Q != "" {
		q.Set("q", p.Q)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *BusinessesListParams) withCursor(c string) *BusinessesListParams {
	cp := BusinessesListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *BusinessesListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// ChannelsListParams são os filtros de consulta de GET /v1/channels (channels_list).
type ChannelsListParams struct {
	BusinessID string
	Status     string
	// Valores: "cloud_api", "coexistence".
	Mode   string
	Limit  int
	Cursor string
}

func (p *ChannelsListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.Status != "" {
		q.Set("status", p.Status)
	}
	if p.Mode != "" {
		q.Set("mode", p.Mode)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ChannelsListParams) withCursor(c string) *ChannelsListParams {
	cp := ChannelsListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *ChannelsListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// ChannelsHistoryParams são os filtros de consulta de GET /v1/channels/{channel_id}/history (channels_history).
type ChannelsHistoryParams struct {
	Limit  int
	Cursor string
}

func (p *ChannelsHistoryParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ChannelsHistoryParams) withCursor(c string) *ChannelsHistoryParams {
	cp := ChannelsHistoryParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *ChannelsHistoryParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// MetaLogsListParams são os filtros de consulta de GET /v1/channels/{channel_id}/meta-logs (meta_logs_list).
type MetaLogsListParams struct {
	// all | problems | notices
	// Valores: "all", "problems", "notices".
	Filter string
	// Default: 7 days ago; at most 90 days
	Since  time.Time
	Limit  int
	Cursor string
}

func (p *MetaLogsListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Filter != "" {
		q.Set("filter", p.Filter)
	}
	if !p.Since.IsZero() {
		q.Set("since", p.Since.UTC().Format(time.RFC3339))
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *MetaLogsListParams) withCursor(c string) *MetaLogsListParams {
	cp := MetaLogsListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *MetaLogsListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// ConnectLinksListParams são os filtros de consulta de GET /v1/connect-links (connect_links_list).
type ConnectLinksListParams struct {
	BusinessID string
	// Valores: "active", "used", "revoked", "expired".
	Status string
	Limit  int
	Cursor string
}

func (p *ConnectLinksListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.Status != "" {
		q.Set("status", p.Status)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ConnectLinksListParams) withCursor(c string) *ConnectLinksListParams {
	cp := ConnectLinksListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *ConnectLinksListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// MessagesListParams são os filtros de consulta de GET /v1/messages (messages_list).
type MessagesListParams struct {
	BusinessID string
	ChannelID  string
	// Valores: "outbound", "inbound".
	Direction string
	// Valores: "queued", "sending", "accepted_by_meta", "delivered", "read", "failed", "canceled", "unknown", "received".
	Status string
	// Valores: "api", "panel", "business_app", "customer", "unknown".
	Origin string
	// E.164 phone of the contact
	Counterparty string
	From         time.Time
	To           time.Time
	Limit        int
	Cursor       string
}

func (p *MessagesListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.ChannelID != "" {
		q.Set("channel_id", p.ChannelID)
	}
	if p.Direction != "" {
		q.Set("direction", p.Direction)
	}
	if p.Status != "" {
		q.Set("status", p.Status)
	}
	if p.Origin != "" {
		q.Set("origin", p.Origin)
	}
	if p.Counterparty != "" {
		q.Set("counterparty", p.Counterparty)
	}
	if !p.From.IsZero() {
		q.Set("from", p.From.UTC().Format(time.RFC3339))
	}
	if !p.To.IsZero() {
		q.Set("to", p.To.UTC().Format(time.RFC3339))
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *MessagesListParams) withCursor(c string) *MessagesListParams {
	cp := MessagesListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *MessagesListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// ConversationsListParams são os filtros de consulta de GET /v1/conversations (conversations_list).
type ConversationsListParams struct {
	ChannelID  string
	BusinessID string
	Limit      int
	Cursor     string
}

func (p *ConversationsListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.ChannelID != "" {
		q.Set("channel_id", p.ChannelID)
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ConversationsListParams) withCursor(c string) *ConversationsListParams {
	cp := ConversationsListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *ConversationsListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// ConversationsTimelineParams são os filtros de consulta de GET /v1/conversations/{channel_id}/{counterparty} (conversations_timeline).
type ConversationsTimelineParams struct {
	Limit  int
	Cursor string
}

func (p *ConversationsTimelineParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ConversationsTimelineParams) withCursor(c string) *ConversationsTimelineParams {
	cp := ConversationsTimelineParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *ConversationsTimelineParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// TemplatesListParams são os filtros de consulta de GET /v1/templates (templates_list).
type TemplatesListParams struct {
	BusinessID string
	// Templates of this channel's WhatsApp account
	ChannelID string
	Status    string
	// Case-insensitive partial match
	Name   string
	Limit  int
	Cursor string
}

func (p *TemplatesListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.ChannelID != "" {
		q.Set("channel_id", p.ChannelID)
	}
	if p.Status != "" {
		q.Set("status", p.Status)
	}
	if p.Name != "" {
		q.Set("name", p.Name)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *TemplatesListParams) withCursor(c string) *TemplatesListParams {
	cp := TemplatesListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *TemplatesListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// TemplatesSyncParams são os filtros de consulta de POST /v1/templates/sync (templates_sync).
type TemplatesSyncParams struct {
	// Obrigatório.
	ChannelID string
}

func (p *TemplatesSyncParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.ChannelID != "" {
		q.Set("channel_id", p.ChannelID)
	}
	return q
}

// WebhookEndpointsListParams são os filtros de consulta de GET /v1/webhook-endpoints (webhook_endpoints_list).
type WebhookEndpointsListParams struct {
	Limit  int
	Cursor string
}

func (p *WebhookEndpointsListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *WebhookEndpointsListParams) withCursor(c string) *WebhookEndpointsListParams {
	cp := WebhookEndpointsListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *WebhookEndpointsListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// WebhookDeliveriesListParams são os filtros de consulta de GET /v1/webhook-deliveries (webhook_deliveries_list).
type WebhookDeliveriesListParams struct {
	EndpointID string
	Status     string
	EventType  string
	EventID    string
	Limit      int
	Cursor     string
}

func (p *WebhookDeliveriesListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.EndpointID != "" {
		q.Set("endpoint_id", p.EndpointID)
	}
	if p.Status != "" {
		q.Set("status", p.Status)
	}
	if p.EventType != "" {
		q.Set("event_type", p.EventType)
	}
	if p.EventID != "" {
		q.Set("event_id", p.EventID)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *WebhookDeliveriesListParams) withCursor(c string) *WebhookDeliveriesListParams {
	cp := WebhookDeliveriesListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *WebhookDeliveriesListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// EventsListParams são os filtros de consulta de GET /v1/events (events_list).
type EventsListParams struct {
	Type   string
	Since  time.Time
	Limit  int
	Cursor string
}

func (p *EventsListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Type != "" {
		q.Set("type", p.Type)
	}
	if !p.Since.IsZero() {
		q.Set("since", p.Since.UTC().Format(time.RFC3339))
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *EventsListParams) withCursor(c string) *EventsListParams {
	cp := EventsListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *EventsListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// UsageGetParams são os filtros de consulta de GET /v1/usage (usage_get).
type UsageGetParams struct {
	BusinessID string
	ChannelID  string
	// YYYY-MM-DD (UTC, inclusive)
	From string
	// YYYY-MM-DD (UTC, inclusive)
	To          string
	Granularity string
}

func (p *UsageGetParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.ChannelID != "" {
		q.Set("channel_id", p.ChannelID)
	}
	if p.From != "" {
		q.Set("from", p.From)
	}
	if p.To != "" {
		q.Set("to", p.To)
	}
	if p.Granularity != "" {
		q.Set("granularity", p.Granularity)
	}
	return q
}

// UsageExportCSVParams são os filtros de consulta de GET /v1/usage/export.csv (usage_export_csv).
type UsageExportCSVParams struct {
	BusinessID string
	ChannelID  string
	From       string
	To         string
}

func (p *UsageExportCSVParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.ChannelID != "" {
		q.Set("channel_id", p.ChannelID)
	}
	if p.From != "" {
		q.Set("from", p.From)
	}
	if p.To != "" {
		q.Set("to", p.To)
	}
	return q
}

// OverviewGetParams são os filtros de consulta de GET /v1/overview (overview_get).
type OverviewGetParams struct {
	BusinessID string
	From       time.Time
	To         time.Time
}

func (p *OverviewGetParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if !p.From.IsZero() {
		q.Set("from", p.From.UTC().Format(time.RFC3339))
	}
	if !p.To.IsZero() {
		q.Set("to", p.To.UTC().Format(time.RFC3339))
	}
	return q
}

// ConsentsListParams são os filtros de consulta de GET /v1/consents (consents_list).
type ConsentsListParams struct {
	BusinessID string
	// Contact phone in E.164.
	Phone  string
	Limit  int
	Cursor string
}

func (p *ConsentsListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.Phone != "" {
		q.Set("phone", p.Phone)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *ConsentsListParams) withCursor(c string) *ConsentsListParams {
	cp := ConsentsListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *ConsentsListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// OptOutsListParams são os filtros de consulta de GET /v1/opt-outs (opt_outs_list).
type OptOutsListParams struct {
	BusinessID string
	// Contact phone in E.164.
	Phone  string
	Limit  int
	Cursor string
}

func (p *OptOutsListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.Phone != "" {
		q.Set("phone", p.Phone)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *OptOutsListParams) withCursor(c string) *OptOutsListParams {
	cp := OptOutsListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *OptOutsListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// OptOutsDeleteParams são os filtros de consulta de DELETE /v1/opt-outs (opt_outs_delete).
type OptOutsDeleteParams struct {
	// Obrigatório.
	BusinessID string
	// Contact phone in E.164.
	// Obrigatório.
	Phone string
}

func (p *OptOutsDeleteParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.Phone != "" {
		q.Set("phone", p.Phone)
	}
	return q
}

// BillingEntriesListParams são os filtros de consulta de GET /v1/billing/entries (billing_entries_list).
type BillingEntriesListParams struct {
	// comma-separated entry kinds
	// Valores: "topup", "topup_refund", "topup_chargeback", "bonus", "bonus_revoked", "bonus_expired", "outbound_message", "inbound_message", "number_fee", "adjustment", "opening_balance", "chargeback_reversal", "closure_refund".
	Kind       string
	BusinessID string
	ChannelID  string
	// YYYY-MM-DD (UTC, inclusive)
	From string
	// YYYY-MM-DD (UTC, inclusive)
	To     string
	Limit  int
	Cursor string
}

func (p *BillingEntriesListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Kind != "" {
		q.Set("kind", p.Kind)
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.ChannelID != "" {
		q.Set("channel_id", p.ChannelID)
	}
	if p.From != "" {
		q.Set("from", p.From)
	}
	if p.To != "" {
		q.Set("to", p.To)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *BillingEntriesListParams) withCursor(c string) *BillingEntriesListParams {
	cp := BillingEntriesListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *BillingEntriesListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// BillingEntriesExportCSVParams são os filtros de consulta de GET /v1/billing/entries/export.csv (billing_entries_export_csv).
type BillingEntriesExportCSVParams struct {
	// comma-separated entry kinds
	// Valores: "topup", "topup_refund", "topup_chargeback", "bonus", "bonus_revoked", "bonus_expired", "outbound_message", "inbound_message", "number_fee", "adjustment", "opening_balance", "chargeback_reversal", "closure_refund".
	Kind       string
	BusinessID string
	ChannelID  string
	// YYYY-MM-DD (UTC, inclusive)
	From string
	// YYYY-MM-DD (UTC, inclusive)
	To string
}

func (p *BillingEntriesExportCSVParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Kind != "" {
		q.Set("kind", p.Kind)
	}
	if p.BusinessID != "" {
		q.Set("business_id", p.BusinessID)
	}
	if p.ChannelID != "" {
		q.Set("channel_id", p.ChannelID)
	}
	if p.From != "" {
		q.Set("from", p.From)
	}
	if p.To != "" {
		q.Set("to", p.To)
	}
	return q
}

// BillingConsumptionParams são os filtros de consulta de GET /v1/billing/consumption (billing_consumption).
type BillingConsumptionParams struct {
	// YYYY-MM-DD (UTC, inclusive)
	From string
	// YYYY-MM-DD (UTC, inclusive)
	To string
	// Valores: "business", "channel", "kind", "day".
	GroupBy string
}

func (p *BillingConsumptionParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.From != "" {
		q.Set("from", p.From)
	}
	if p.To != "" {
		q.Set("to", p.To)
	}
	if p.GroupBy != "" {
		q.Set("group_by", p.GroupBy)
	}
	return q
}

// BillingNumberFeesListParams são os filtros de consulta de GET /v1/billing/number-fees (billing_number_fees_list).
type BillingNumberFeesListParams struct {
	ChannelID string
	// Valores: "charged", "past_due", "exempt", "waived".
	Status string
	Limit  int
	Cursor string
}

func (p *BillingNumberFeesListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.ChannelID != "" {
		q.Set("channel_id", p.ChannelID)
	}
	if p.Status != "" {
		q.Set("status", p.Status)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *BillingNumberFeesListParams) withCursor(c string) *BillingNumberFeesListParams {
	cp := BillingNumberFeesListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *BillingNumberFeesListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}

// TopupsListParams são os filtros de consulta de GET /v1/billing/topups (topups_list).
type TopupsListParams struct {
	// Valores: "pending", "paid", "failed", "canceled", "refund_processing", "refunded", "charged_back".
	Status string
	Limit  int
	Cursor string
}

func (p *TopupsListParams) query() url.Values {
	q := url.Values{}
	if p == nil {
		return q
	}
	if p.Status != "" {
		q.Set("status", p.Status)
	}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}
	return q
}

func (p *TopupsListParams) withCursor(c string) *TopupsListParams {
	cp := TopupsListParams{}
	if p != nil {
		cp = *p
	}
	cp.Cursor = c
	return &cp
}

func (p *TopupsListParams) cursorValue() string {
	if p == nil {
		return ""
	}
	return p.Cursor
}
