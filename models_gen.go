// Código gerado por internal/gen/genmodels.py a partir do contrato. NÃO EDITE.

package meshio

import "time"

// Allowance corresponde ao schema AllowanceOut do contrato.
type Allowance struct {
	// null when there is no open cycle (see access_reason)
	CycleID *string `json:"cycle_id"`
	// Valores: "trial", "paid".
	Kind        *string    `json:"kind"`
	PeriodStart *time.Time `json:"period_start"`
	PeriodEnd   *time.Time `json:"period_end"`
	// why there is no open cycle; null otherwise
	// Valores: "trial_ended", "no_plan", "subscription_canceled", "subscription_incomplete", "trial_available".
	AccessReason     *string `json:"access_reason"`
	IncludedMessages int64   `json:"included_messages"`
	// messages counted in the cycle (delivered sends + new received)
	UsedMessages int64 `json:"used_messages"`
	// accepted sends waiting for delivery (reserved)
	HeldMessages int64 `json:"held_messages"`
	// included - used - held, never below 0
	RemainingMessages int64 `json:"remaining_messages"`
	// used / included, integer percent (may exceed 100)
	Percent          int64  `json:"percent"`
	OutboundMessages int64  `json:"outbound_messages"`
	InboundMessages  int64  `json:"inbound_messages"`
	OverageMessages  int64  `json:"overage_messages"`
	OverageMicros    string `json:"overage_micros"`
	// snapshot of the cycle; null = overage unavailable
	OveragePriceMicros *string `json:"overage_price_micros"`
	OverageEnabled     bool    `json:"overage_enabled"`
	OverageLimitMicros *string `json:"overage_limit_micros"`
	// received above the allowance without overage (no charge)
	ExcessUnbilledMessages int64         `json:"excess_unbilled_messages"`
	ByNumber               []NumberUsage `json:"by_number"`
}

// ApprovalSummary corresponde ao schema ApprovalSummaryOut do contrato.
type ApprovalSummary struct {
	// Valores: "pending", "approved", "rejected".
	Status string `json:"status"`
	// Only when rejected
	Reason    *string    `json:"reason"`
	DecidedAt *time.Time `json:"decided_at"`
}

// Attention corresponde ao schema AttentionOut do contrato.
type Attention struct {
	// Valores: "number_limited", "number_blocked", "quality_low", "template_paused", "webhook_failing", "billing".
	Kind string `json:"kind"`
	// pt-BR message for the panel
	Message string `json:"message"`
	// Panel path to act on the item
	Href string `json:"href"`
}

// BillingStatus corresponde ao schema BillingStatusOut do contrato.
type BillingStatus struct {
	// none = never had a plan (trial still available); trial = free trial open; trial_ended = trial over without a subscription; incomplete = awaiting the first payment; null = the caller cannot read billing
	// Valores: "none", "trial", "trial_ended", "incomplete", "active", "past_due", "canceled".
	SubscriptionStatus *string `json:"subscription_status"`
	// floor(counted / included * 100) of the open cycle; null without an open cycle
	AllowancePercent *int64 `json:"allowance_percent"`
	// End of the open cycle (renewal or trial end)
	CycleEndsAt *time.Time `json:"cycle_ends_at"`
	// Most relevant billing issue; also listed in `attention` with a pt-BR message
	// Valores: "past_due_blocked", "past_due", "subscription_required", "allowance_exhausted", "allowance_warning", "trial_ending".
	Attention *string `json:"attention"`
}

// BusinessCompliance corresponde ao schema BusinessComplianceOut do contrato.
type BusinessCompliance struct {
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	// Valores: "pending", "verified", "rejected", "manual_review", "unavailable".
	VerificationStatus string     `json:"verification_status"`
	Cnpj               *string    `json:"cnpj"`
	Reasons            []string   `json:"reasons"`
	CheckedAt          *time.Time `json:"checked_at"`
	Suspended          bool       `json:"suspended"`
}

// Business corresponde ao schema BusinessOut do contrato.
type Business struct {
	ID              string    `json:"id"`
	IntegratorID    string    `json:"integrator_id"`
	Name            string    `json:"name"`
	LegalName       *string   `json:"legal_name"`
	TaxID           *string   `json:"tax_id"`
	Country         string    `json:"country"`
	Timezone        string    `json:"timezone"`
	ContactEmail    string    `json:"contact_email"`
	ContactPhone    *string   `json:"contact_phone"`
	Purpose         string    `json:"purpose"`
	TermsAcceptedAt time.Time `json:"terms_accepted_at"`
	TermsAcceptedBy string    `json:"terms_accepted_by"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BusinessPage corresponde ao schema BusinessPage do contrato.
type BusinessPage struct {
	Data       []Business `json:"data"`
	NextCursor *string    `json:"next_cursor"`
}

// CancelSignupParams corresponde ao schema CancelSignupIn do contrato.
type CancelSignupParams struct {
	Reason    string `json:"reason,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}

// Capabilities corresponde ao schema CapabilitiesOut do contrato.
type Capabilities struct {
	SendText          bool `json:"send_text"`
	SendTemplate      bool `json:"send_template"`
	Receive           bool `json:"receive"`
	BusinessAppEchoes bool `json:"business_app_echoes"`
	HistorySync       bool `json:"history_sync"`
}

// ChannelHistoryItem corresponde ao schema ChannelHistoryItem do contrato.
type ChannelHistoryItem struct {
	ID         string  `json:"id"`
	ChannelID  string  `json:"channel_id"`
	FromStatus *string `json:"from_status"`
	ToStatus   string  `json:"to_status"`
	// Valores: "user", "meta", "system", "health_check".
	Source     string    `json:"source"`
	Cause      *string   `json:"cause"`
	OccurredAt time.Time `json:"occurred_at"`
}

// ChannelHistoryPage corresponde ao schema ChannelHistoryPage do contrato.
type ChannelHistoryPage struct {
	Data       []ChannelHistoryItem `json:"data"`
	NextCursor *string              `json:"next_cursor"`
}

// Channel corresponde ao schema ChannelOut do contrato.
type Channel struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	// Valores: "cloud_api", "coexistence".
	Mode string `json:"mode"`
	// Valores: "not_started", "pending_authorization", "pending_setup", "active", "attention_required", "disconnected", "revoked".
	Status             string         `json:"status"`
	StatusReason       *string        `json:"status_reason"`
	DisplayPhoneNumber *string        `json:"display_phone_number"`
	VerifiedName       *string        `json:"verified_name"`
	WABAID             string         `json:"waba_id"`
	PhoneNumberID      string         `json:"phone_number_id"`
	QualityRating      *string        `json:"quality_rating"`
	Capabilities       Capabilities   `json:"capabilities"`
	LastHealthCheckAt  *time.Time     `json:"last_health_check_at"`
	LastHealthStatus   map[string]any `json:"last_health_status"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	// Valores: "simulator", "live".
	Environment string     `json:"environment"`
	MetaHealth  MetaHealth `json:"meta_health"`
}

// ChannelPage corresponde ao schema ChannelPage do contrato.
type ChannelPage struct {
	Data       []Channel `json:"data"`
	NextCursor *string   `json:"next_cursor"`
}

// CompleteSignupParams corresponde ao schema CompleteSignupIn do contrato.
type CompleteSignupParams struct {
	State         string `json:"state"`
	Code          string `json:"code"`
	WABAID        string `json:"waba_id"`
	PhoneNumberID string `json:"phone_number_id"`
	Event         string `json:"event"`
}

// ComplianceStatus corresponde ao schema ComplianceStatusOut do contrato.
type ComplianceStatus struct {
	IntegratorID string               `json:"integrator_id"`
	Approval     ApprovalSummary      `json:"approval"`
	Businesses   []BusinessCompliance `json:"businesses"`
	// Active holds visible to the caller
	Holds []HoldSummary `json:"holds"`
	// Blocked sends/templates by reason since blocked_since
	BlockedCounts map[string]int64 `json:"blocked_counts"`
	BlockedSince  time.Time        `json:"blocked_since"`
}

// ConnectLink corresponde ao schema ConnectLinkOut do contrato.
type ConnectLink struct {
	ID           string `json:"id"`
	BusinessID   string `json:"business_id"`
	BusinessName string `json:"business_name"`
	// Valores: "cloud_api", "coexistence", "choose".
	Mode string `json:"mode"`
	// Valores: "active", "used", "revoked", "expired".
	Status string `json:"status"`
	// Hosted connection URL; only while the link is active.
	URL              *string    `json:"url"`
	TokenPrefix      string     `json:"token_prefix"`
	ExpiresAt        time.Time  `json:"expires_at"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedByDisplay *string    `json:"created_by_display"`
	UsedAt           *time.Time `json:"used_at"`
	ChannelID        *string    `json:"channel_id"`
	RevokedAt        *time.Time `json:"revoked_at"`
}

// ConnectLinkPage corresponde ao schema ConnectLinkPage do contrato.
type ConnectLinkPage struct {
	Data       []ConnectLink `json:"data"`
	NextCursor *string       `json:"next_cursor"`
}

// ConsentBulkParams corresponde ao schema ConsentBulkIn do contrato.
type ConsentBulkParams struct {
	Consents []ConsentParams `json:"consents"`
}

// ConsentBulk corresponde ao schema ConsentBulkOut do contrato.
type ConsentBulk struct {
	Data     []ConsentRecorded `json:"data"`
	Created  int64             `json:"created"`
	Existing int64             `json:"existing"`
}

// ConsentParams corresponde ao schema ConsentIn do contrato.
type ConsentParams struct {
	BusinessID string `json:"business_id"`
	// Contact phone in E.164.
	Phone string `json:"phone"`
	// Valores: "checkout", "website_form", "whatsapp", "in_person", "other".
	Source      string    `json:"source"`
	CapturedAt  time.Time `json:"captured_at"`
	EvidenceURL string    `json:"evidence_url,omitempty"`
}

// Consent corresponde ao schema ConsentOut do contrato.
type Consent struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	Phone      string `json:"phone"`
	// Valores: "checkout", "website_form", "whatsapp", "in_person", "other".
	Source      string    `json:"source"`
	CapturedAt  time.Time `json:"captured_at"`
	EvidenceURL *string   `json:"evidence_url,omitempty"`
	// Valores: "api", "send_inline", "panel", "whatsapp_keyword".
	CreatedVia string     `json:"created_via"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	// Valores: "active", "revoked".
	Status string `json:"status"`
}

// ConsentPage corresponde ao schema ConsentPage do contrato.
type ConsentPage struct {
	Data       []Consent `json:"data"`
	NextCursor *string   `json:"next_cursor,omitempty"`
}

// ConsentRecorded corresponde ao schema ConsentRecordedOut do contrato.
type ConsentRecorded struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	Phone      string `json:"phone"`
	// Valores: "checkout", "website_form", "whatsapp", "in_person", "other".
	Source      string    `json:"source"`
	CapturedAt  time.Time `json:"captured_at"`
	EvidenceURL *string   `json:"evidence_url,omitempty"`
	// Valores: "api", "send_inline", "panel", "whatsapp_keyword".
	CreatedVia string     `json:"created_via"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	// Valores: "active", "revoked".
	Status string `json:"status"`
	// False when a live consent already existed (kept as is).
	Created bool `json:"created"`
	// True when this consent (captured after the opt-out) removed the contact from the opt-out list.
	OptOutLifted bool `json:"opt_out_lifted"`
}

// Consumption corresponde ao schema ConsumptionOut do contrato.
type Consumption struct {
	ID            string  `json:"id"`
	MessageID     string  `json:"message_id"`
	MetaMessageID *string `json:"meta_message_id"`
	CycleID       *string `json:"cycle_id"`
	BusinessID    *string `json:"business_id"`
	ChannelID     *string `json:"channel_id"`
	// Valores: "outbound", "inbound".
	Direction string  `json:"direction"`
	Origin    *string `json:"origin"`
	// Valores: "included", "overage", "excess_unbilled", "no_cycle".
	Classification     string  `json:"classification"`
	OveragePriceMicros *string `json:"overage_price_micros"`
	// Valores: "delivered", "read", "inbound_received".
	Evidence        string  `json:"evidence"`
	StatusHistoryID *string `json:"status_history_id"`
	// when Meshio accepted the send (outbound)
	AcceptedAt *time.Time `json:"accepted_at"`
	// Meta timestamp of the delivery/receipt (fallback: receipt)
	OccurredAt time.Time `json:"occurred_at"`
	CountedAt  time.Time `json:"counted_at"`
	// counted in a different cycle from the one of acceptance/occurrence
	Late        bool  `json:"late"`
	RuleVersion int64 `json:"rule_version"`
}

// ConsumptionPage corresponde ao schema ConsumptionPage do contrato.
type ConsumptionPage struct {
	Data       []Consumption `json:"data"`
	NextCursor *string       `json:"next_cursor"`
}

// Conversation corresponde ao schema ConversationOut do contrato.
type Conversation struct {
	ChannelID         string             `json:"channel_id"`
	BusinessID        string             `json:"business_id"`
	CounterpartyPhone string             `json:"counterparty_phone"`
	ContactName       *string            `json:"contact_name,omitempty"`
	MessageCount      int64              `json:"message_count"`
	LastMessage       Message            `json:"last_message"`
	Window            ConversationWindow `json:"window"`
}

// ConversationPage corresponde ao schema ConversationPage do contrato.
type ConversationPage struct {
	Data       []Conversation `json:"data"`
	NextCursor *string        `json:"next_cursor"`
}

// ConversationWindow corresponde ao schema ConversationWindowOut do contrato.
type ConversationWindow struct {
	// True while free-form text is allowed (24h customer service window).
	Open          bool       `json:"open"`
	LastInboundAt *time.Time `json:"last_inbound_at,omitempty"`
	ClosesAt      *time.Time `json:"closes_at,omitempty"`
}

// CreateConnectLinkParams corresponde ao schema CreateConnectLinkIn do contrato.
type CreateConnectLinkParams struct {
	BusinessID string `json:"business_id"`
	// `choose` lets the business owner pick the mode on the link page.
	// Valores: "cloud_api", "coexistence", "choose".
	Mode           string `json:"mode,omitempty"`
	ExpiresInHours *int64 `json:"expires_in_hours,omitempty"`
}

// CreateSignupSessionParams corresponde ao schema CreateSignupSessionIn do contrato.
type CreateSignupSessionParams struct {
	BusinessID string `json:"business_id"`
	// Valores: "cloud_api", "coexistence".
	Mode string `json:"mode"`
}

// CreateSubscriptionParams corresponde ao schema CreateSubscriptionIn do contrato.
type CreateSubscriptionParams struct {
	Seats int64 `json:"seats"`
	// invoice_link (Pix or card per invoice) | wallet (preserved paid balance; requires accept_wallet_terms). card_auto_renewal and pix_automatic → 422 PAYMENT_METHOD_UNAVAILABLE
	// Valores: "invoice_link", "wallet", "card_auto_renewal", "pix_automatic".
	PaymentMethod     string `json:"payment_method,omitempty"`
	AcceptWalletTerms *bool  `json:"accept_wallet_terms,omitempty"`
}

// CreateTemplateParams corresponde ao schema CreateTemplateIn do contrato.
type CreateTemplateParams struct {
	BusinessID string `json:"business_id"`
	ChannelID  string `json:"channel_id"`
	// lowercase letters, digits and underscore
	Name     string `json:"name"`
	Language string `json:"language"`
	// Valores: "MARKETING", "UTILITY", "AUTHENTICATION".
	Category string `json:"category"`
	// Use {{1}}, {{2}}... for variables
	BodyText      string                 `json:"body_text,omitempty"`
	BodyExamples  []string               `json:"body_examples,omitempty"`
	HeaderText    string                 `json:"header_text,omitempty"`
	HeaderExample string                 `json:"header_example,omitempty"`
	FooterText    string                 `json:"footer_text,omitempty"`
	Buttons       []TemplateButtonParams `json:"buttons,omitempty"`
	// Meta-format components (BODY/HEADER/FOOTER/BUTTONS)
	Components []map[string]any `json:"components,omitempty"`
}

// CreateWebhookEndpointParams corresponde ao schema CreateWebhookEndpointIn do contrato.
type CreateWebhookEndpointParams struct {
	URL         string   `json:"url"`
	EventTypes  []string `json:"event_types"`
	Description string   `json:"description,omitempty"`
	// Valores: "enabled", "disabled".
	Status string `json:"status,omitempty"`
}

// CyclePage corresponde ao schema CyclePage do contrato.
type CyclePage struct {
	Data       []CycleSummary `json:"data"`
	NextCursor *string        `json:"next_cursor"`
}

// CycleSummary corresponde ao schema CycleSummaryOut do contrato.
type CycleSummary struct {
	ID string `json:"id"`
	// Valores: "trial", "paid".
	Kind string `json:"kind"`
	// Valores: "open", "closed".
	Status                 string    `json:"status"`
	PeriodStart            time.Time `json:"period_start"`
	PeriodEnd              time.Time `json:"period_end"`
	Seats                  int64     `json:"seats"`
	MaxNumbers             int64     `json:"max_numbers"`
	IncludedMessages       int64     `json:"included_messages"`
	UsedMessages           int64     `json:"used_messages"`
	HeldMessages           int64     `json:"held_messages"`
	RemainingMessages      int64     `json:"remaining_messages"`
	OverageMessages        int64     `json:"overage_messages"`
	OverageMicros          string    `json:"overage_micros"`
	ExcessUnbilledMessages int64     `json:"excess_unbilled_messages"`
	UnitPriceMicros        string    `json:"unit_price_micros"`
	OveragePriceMicros     *string   `json:"overage_price_micros"`
	Complimentary          bool      `json:"complimentary"`
}

// DeliveryStats corresponde ao schema DeliveryStatsOut do contrato.
type DeliveryStats struct {
	// Deliveries created in the window.
	Total     int64 `json:"total"`
	Succeeded int64 `json:"succeeded"`
	// Deliveries that exhausted retries (final failure).
	Failed int64 `json:"failed"`
	// succeeded / (succeeded + failed), 0..1; null when nothing finished yet.
	SuccessRate *float64 `json:"success_rate"`
}

// DisconnectParams corresponde ao schema DisconnectIn do contrato.
type DisconnectParams struct {
	Confirm *bool `json:"confirm,omitempty"`
}

// Entry corresponde ao schema EntryOut do contrato.
type Entry struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	// Valores: "topup", "topup_refund", "topup_chargeback", "bonus", "bonus_revoked", "bonus_expired", "outbound_message", "inbound_message", "number_fee", "adjustment", "opening_balance", "chargeback_reversal", "closure_refund", "invoice_payment".
	Kind string `json:"kind"`
	// signed: + credit to the wallet, - debit
	AmountMicros string `json:"amount_micros"`
	// part from/to the paid balance
	PaidMicros string `json:"paid_micros"`
	// part from/to the bonus balance
	BonusMicros        string  `json:"bonus_micros"`
	BalanceAfterMicros string  `json:"balance_after_micros"`
	Quantity           *int64  `json:"quantity"`
	UnitPriceMicros    *string `json:"unit_price_micros"`
	// Valores: "default", "custom".
	PriceSource            *string        `json:"price_source"`
	PricingDefaultVersion  *int64         `json:"pricing_default_version"`
	PricingOverrideVersion *int64         `json:"pricing_override_version"`
	BillingExempt          bool           `json:"billing_exempt"`
	RatedMicros            *string        `json:"rated_micros"`
	BusinessID             *string        `json:"business_id"`
	BusinessName           *string        `json:"business_name"`
	ChannelID              *string        `json:"channel_id"`
	ChannelDisplayPhone    *string        `json:"channel_display_phone"`
	ReferenceType          *string        `json:"reference_type"`
	ReferenceID            *string        `json:"reference_id"`
	Reason                 *string        `json:"reason"`
	ActorType              *string        `json:"actor_type"`
	ActorDisplay           *string        `json:"actor_display"`
	Details                map[string]any `json:"details"`
}

// EntryPage corresponde ao schema EntryPage do contrato.
type EntryPage struct {
	Data       []Entry `json:"data"`
	NextCursor *string `json:"next_cursor"`
}

// Event corresponde ao schema EventOut do contrato.
type Event struct {
	ID         string         `json:"id"`
	EventType  string         `json:"event_type"`
	OccurredAt time.Time      `json:"occurred_at"`
	CreatedAt  time.Time      `json:"created_at"`
	BusinessID *string        `json:"business_id"`
	ChannelID  *string        `json:"channel_id"`
	Data       map[string]any `json:"data"`
}

// EventPage corresponde ao schema EventPage do contrato.
type EventPage struct {
	Data       []Event `json:"data"`
	NextCursor *string `json:"next_cursor"`
}

// HTTPValidationError corresponde ao schema HTTPValidationError do contrato.
type HTTPValidationError struct {
	Detail []ValidationErrorItem `json:"detail,omitempty"`
}

// HoldSummary corresponde ao schema HoldSummaryOut do contrato.
type HoldSummary struct {
	ID string `json:"id"`
	// Valores: "channel", "business", "integrator".
	Scope      string  `json:"scope"`
	BusinessID *string `json:"business_id"`
	ChannelID  *string `json:"channel_id"`
	// Valores: "meta_account_violation", "meta_account_restriction", "meta_account_banned", "meta_account_disabled", "quality_red", "templates_paused", "meta_policy_errors", "manual".
	Reason      string `json:"reason"`
	ReasonLabel string `json:"reason_label"`
	// Valores: "meta_signal", "operator", "policy".
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

// IntegrationStatus corresponde ao schema IntegrationStatusOut do contrato.
type IntegrationStatus struct {
	Numbers      []StatusNumber `json:"numbers"`
	NumbersTotal int64          `json:"numbers_total"`
	// Active numbers that Meta does not block
	NumbersCanSend int64           `json:"numbers_can_send"`
	Messages24h    Messages24h     `json:"messages_24h"`
	Webhooks24h    Webhooks24h     `json:"webhooks_24h"`
	Templates      TemplatesStatus `json:"templates"`
	Billing        BillingStatus   `json:"billing"`
	Attention      []Attention     `json:"attention"`
}

// InvoiceLine corresponde ao schema InvoiceLineOut do contrato.
type InvoiceLine struct {
	Description     string     `json:"description"`
	Quantity        int64      `json:"quantity"`
	UnitPriceMicros string     `json:"unit_price_micros"`
	AmountMicros    string     `json:"amount_micros"`
	PeriodStart     *time.Time `json:"period_start"`
	PeriodEnd       *time.Time `json:"period_end"`
}

// Invoice corresponde ao schema InvoiceOut do contrato.
type Invoice struct {
	ID             string  `json:"id"`
	SubscriptionID string  `json:"subscription_id"`
	CycleID        *string `json:"cycle_id"`
	// Valores: "subscription", "seat_addition", "overage".
	Kind string `json:"kind"`
	// Valores: "open", "paid", "void", "uncollectible".
	Status string `json:"status"`
	// Valor: "BRL".
	Currency    string    `json:"currency,omitempty"`
	TotalMicros string    `json:"total_micros"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	IssuedAt    time.Time `json:"issued_at"`
	DueAt       time.Time `json:"due_at"`
	// sends keep working until this moment if unpaid
	GraceUntil *time.Time `json:"grace_until"`
	PaidAt     *time.Time `json:"paid_at"`
	// Valores: "pix", "card", "wallet", "manual".
	PaidVia *string `json:"paid_via"`
	// payment link (Pix or card) while the invoice is open
	CheckoutURL       *string    `json:"checkout_url"`
	CheckoutExpiresAt *time.Time `json:"checkout_expires_at"`
	// true = simulated provider (dev/test), no real money
	Simulated bool          `json:"simulated,omitempty"`
	Lines     []InvoiceLine `json:"lines"`
}

// InvoicePage corresponde ao schema InvoicePage do contrato.
type InvoicePage struct {
	Data       []Invoice `json:"data"`
	NextCursor *string   `json:"next_cursor"`
}

// Launch corresponde ao schema LaunchOut do contrato.
type Launch struct {
	// Valores: "facebook_sdk", "simulator".
	Kind   string         `json:"kind"`
	URL    *string        `json:"url,omitempty"`
	Extras map[string]any `json:"extras,omitempty"`
}

// MessageContext corresponde ao schema MessageContextOut do contrato.
type MessageContext struct {
	// wamid of the message the contact replied to (Meta `context.id`).
	MetaMessageID string `json:"meta_message_id"`
	// Meshio id of the replied message, when it belongs to this integrator.
	MessageID *string `json:"message_id,omitempty"`
	// Type of the replied message, when found.
	Type *string `json:"type,omitempty"`
	// Template of the replied message, when found.
	TemplateName *string `json:"template_name,omitempty"`
}

// MessageError corresponde ao schema MessageErrorOut do contrato.
type MessageError struct {
	// Safe Meshio error code (UPPER_SNAKE_CASE).
	Code    string  `json:"code"`
	Message *string `json:"message,omitempty"`
	Kind    *string `json:"kind,omitempty"`
	// Original Meta error code, when available.
	MetaCode  *int64  `json:"meta_code,omitempty"`
	FbtraceID *string `json:"fbtrace_id,omitempty"`
}

// Message corresponde ao schema MessageOut do contrato.
type Message struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	ChannelID  string `json:"channel_id"`
	// Valores: "outbound", "inbound".
	Direction string `json:"direction"`
	// Valores: "api", "panel", "business_app", "customer", "unknown".
	Origin string `json:"origin"`
	// `interactive_reply` (inbound only): the contact tapped a template quick reply button or answered an interactive message; see `reply`. New types may be added.
	// Valores: "text", "template", "interactive_reply", "unsupported".
	Type              string `json:"type"`
	CounterpartyPhone string `json:"counterparty_phone"`
	// Text content; for `interactive_reply`, the tapped button/row title.
	BodyText *string          `json:"body_text,omitempty"`
	Template *MessageTemplate `json:"template,omitempty"`
	// Click data for `interactive_reply`.
	Reply *MessageReply `json:"reply,omitempty"`
	// Message the contact replied to (quoted reply or button click).
	Context         *MessageContext `json:"context,omitempty"`
	ClientReference *string         `json:"client_reference,omitempty"`
	// Valores: "queued", "sending", "accepted_by_meta", "delivered", "read", "failed", "canceled", "unknown", "received".
	Status        string        `json:"status"`
	MetaMessageID *string       `json:"meta_message_id,omitempty"`
	Error         *MessageError `json:"error,omitempty"`
	Attempts      int64         `json:"attempts"`
	// True for messages imported from the coexistence history sync.
	Imported    bool    `json:"imported"`
	ContactName *string `json:"contact_name,omitempty"`
	// User (`usr_…`) or API key (`key_…`) that sent it (outbound only).
	CreatedBy *string `json:"created_by,omitempty"`
	// Human-readable sender; filled in `GET /v1/messages/{id}` only.
	CreatedByDisplay *string         `json:"created_by_display,omitempty"`
	MetaTimestamp    *time.Time      `json:"meta_timestamp,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	StatusHistory    []StatusHistory `json:"status_history,omitempty"`
}

// MessagePage corresponde ao schema MessagePage do contrato.
type MessagePage struct {
	Data       []Message `json:"data"`
	NextCursor *string   `json:"next_cursor"`
}

// MessageReply corresponde ao schema MessageReplyOut do contrato.
type MessageReply struct {
	// `template_quick_reply`: quick reply button of a template; `button_reply`/`list_reply`: reply to an interactive message.
	// Valores: "template_quick_reply", "button_reply", "list_reply".
	Kind string `json:"kind"`
	// Button payload (template quick reply) or the reply/row id (interactive).
	ID *string `json:"id,omitempty"`
	// Text the contact saw and tapped (also in `body_text`).
	Title string `json:"title"`
	// List row description (`list_reply` only).
	Description *string `json:"description,omitempty"`
}

// MessageTemplateButton corresponde ao schema MessageTemplateButtonOut do contrato.
type MessageTemplateButton struct {
	// Valores: "quick_reply", "url", "phone", "other".
	Type string `json:"type"`
	Text string `json:"text"`
}

// MessageTemplateContent corresponde ao schema MessageTemplateContentOut do contrato.
type MessageTemplateContent struct {
	// TEXT, IMAGE, VIDEO, DOCUMENT or LOCATION; null without header.
	HeaderFormat *string                 `json:"header_format,omitempty"`
	HeaderText   *string                 `json:"header_text,omitempty"`
	BodyText     string                  `json:"body_text"`
	FooterText   *string                 `json:"footer_text,omitempty"`
	Buttons      []MessageTemplateButton `json:"buttons,omitempty"`
	// `snapshot`: frozen when the message was sent; `template`: rebuilt from the current template (messages sent before snapshots existed; may differ if the template changed).
	// Valores: "snapshot", "template".
	Source string `json:"source"`
}

// MessageTemplate corresponde ao schema MessageTemplateOut do contrato.
type MessageTemplate struct {
	Name      string            `json:"name"`
	Language  string            `json:"language"`
	Variables TemplateVariables `json:"variables"`
	// Structure of the template as the contact sees it; null when the template no longer exists.
	Content *MessageTemplateContent `json:"content,omitempty"`
}

// Messages24h corresponde ao schema Messages24hOut do contrato.
type Messages24h struct {
	// Outbound messages accepted by Meta in the last 24 hours
	Sent   int64 `json:"sent"`
	Failed int64 `json:"failed"`
	// delivered or read / sent; null when nothing was sent
	DeliveredRatio *float64 `json:"delivered_ratio"`
	// Up to 5 most frequent failure codes
	TopErrors []TopError `json:"top_errors"`
}

// MetaHealthCapability corresponde ao schema MetaHealthCapabilityOut do contrato.
type MetaHealthCapability struct {
	// reply: answer who wrote in the last 24 hours; initiate: start conversations (templates)
	// Valores: "reply", "initiate".
	Key string `json:"key"`
	// pt-BR label
	Label string `json:"label"`
	// Valores: "allowed", "limited", "blocked", "unknown".
	Status string `json:"status"`
	// pt-BR detail (may be empty)
	Detail string `json:"detail"`
}

// MetaHealthExplanation corresponde ao schema MetaHealthExplanationOut do contrato.
type MetaHealthExplanation struct {
	// Valores: "available", "limited", "initiate_blocked", "blocked", "unknown".
	State string `json:"state"`
	// pt-BR sentence
	Headline string `json:"headline"`
	// exactly two: reply and initiate
	Capabilities []MetaHealthCapability `json:"capabilities"`
	// relevant issues: blocks, limits, info; catalogued before uncatalogued
	Issues []MetaHealthIssue `json:"issues"`
	// pt-BR one-liners for issues that do not affect messages
	Ignored []string `json:"ignored"`
}

// MetaHealthIssue corresponde ao schema MetaHealthIssueOut do contrato.
type MetaHealthIssue struct {
	// Meta error_code (null for additional_info)
	Code *int64 `json:"code"`
	// PHONE_NUMBER, WABA, BUSINESS or APP
	EntityType string  `json:"entity_type"`
	EntityID   *string `json:"entity_id"`
	// Valores: "blocks", "limits", "info".
	Severity string `json:"severity"`
	// pt-BR title
	Title string `json:"title"`
	// pt-BR: where the issue is (number, WhatsApp account, portfolio, app)
	ScopeLabel string `json:"scope_label"`
	// pt-BR effect on sending
	Effect string `json:"effect"`
	// pt-BR instructions
	HowToFix string `json:"how_to_fix"`
	// Meta page where the issue is fixed
	FixURL *string `json:"fix_url"`
	// Meta's original text (verbatim, usually English)
	MetaText string `json:"meta_text"`
	// false: code/text not catalogued by Meshio yet (operator is alerted)
	Known bool `json:"known"`
}

// MetaHealth corresponde ao schema MetaHealthOut do contrato.
type MetaHealth struct {
	// Meta health_status.can_send_message; null until the first reading
	// Valores: "available", "limited", "blocked".
	CanSend *string `json:"can_send"`
	// Meta's first error/additional info when limited or blocked (verbatim)
	Reason *string `json:"reason"`
	// Display name review status (APPROVED, PENDING_REVIEW, ...)
	NameStatus *string `json:"name_status"`
	// Messaging limit tier (TIER_250, TIER_2K, ... TIER_UNLIMITED)
	MessagingLimit *string `json:"messaging_limit"`
	QualityRating  *string `json:"quality_rating"`
	// Last time Meta returned the health status
	CheckedAt *time.Time `json:"checked_at"`
	// Explained health; null until Meta returns the health status
	Explanation *MetaHealthExplanation `json:"explanation"`
}

// MetaLogEntryDetail corresponde ao schema MetaLogEntryDetailOut do contrato.
type MetaLogEntryDetail struct {
	// mcl_… (Meshio → Meta call) or iev_… (Meta → Meshio webhook)
	ID         string    `json:"id"`
	OccurredAt time.Time `json:"occurred_at"`
	// Valores: "outbound", "inbound".
	Direction string `json:"direction"`
	// Graph call operation (e.g. get_phone_number) or webhook field
	Operation string `json:"operation"`
	// pt-BR title, e.g. 'Registrar número na Cloud API'
	Title string `json:"title"`
	// pt-BR one-liner (origin · METHOD path)
	Summary string        `json:"summary"`
	Result  MetaLogResult `json:"result"`
	// HTTP method (outbound only)
	Method *string `json:"method,omitempty"`
	// Graph path without version or query (outbound only)
	Path                *string `json:"path,omitempty"`
	GraphAPIVersion     *string `json:"graph_api_version,omitempty"`
	HTTPStatus          *int64  `json:"http_status"`
	DurationMs          *int64  `json:"duration_ms"`
	MetaErrorCode       *int64  `json:"meta_error_code"`
	FbtraceID           *string `json:"fbtrace_id"`
	IntegratorID        *string `json:"integrator_id"`
	IntegratorName      *string `json:"integrator_name"`
	ChannelID           *string `json:"channel_id"`
	ChannelDisplayPhone *string `json:"channel_display_phone"`
	WABAID              *string `json:"waba_id"`
	PhoneNumberID       *string `json:"phone_number_id"`
	// job:<kind> | api:<METHOD> <path_template> | cli:<name>
	Trigger   *string `json:"trigger"`
	RequestID *string `json:"request_id"`
	// permanent | until:<YYYY-MM-DD> | expired | none
	ContentState  string `json:"content_state"`
	AccountNotice bool   `json:"account_notice"`
	// Outbound: {query, body} already scrubbed (tokens never stored)
	Request any `json:"request,omitempty"`
	// Outbound: Meta response body (scrubbed)
	Response any `json:"response,omitempty"`
	// Inbound: raw webhook change
	Payload any `json:"payload,omitempty"`
	// pt-BR note on how long the content is kept
	RetentionNote string `json:"retention_note"`
}

// MetaLogEntry corresponde ao schema MetaLogEntryOut do contrato.
type MetaLogEntry struct {
	// mcl_… (Meshio → Meta call) or iev_… (Meta → Meshio webhook)
	ID         string    `json:"id"`
	OccurredAt time.Time `json:"occurred_at"`
	// Valores: "outbound", "inbound".
	Direction string `json:"direction"`
	// Graph call operation (e.g. get_phone_number) or webhook field
	Operation string `json:"operation"`
	// pt-BR title, e.g. 'Registrar número na Cloud API'
	Title string `json:"title"`
	// pt-BR one-liner (origin · METHOD path)
	Summary string        `json:"summary"`
	Result  MetaLogResult `json:"result"`
	// HTTP method (outbound only)
	Method *string `json:"method,omitempty"`
	// Graph path without version or query (outbound only)
	Path                *string `json:"path,omitempty"`
	GraphAPIVersion     *string `json:"graph_api_version,omitempty"`
	HTTPStatus          *int64  `json:"http_status"`
	DurationMs          *int64  `json:"duration_ms"`
	MetaErrorCode       *int64  `json:"meta_error_code"`
	FbtraceID           *string `json:"fbtrace_id"`
	IntegratorID        *string `json:"integrator_id"`
	IntegratorName      *string `json:"integrator_name"`
	ChannelID           *string `json:"channel_id"`
	ChannelDisplayPhone *string `json:"channel_display_phone"`
	WABAID              *string `json:"waba_id"`
	PhoneNumberID       *string `json:"phone_number_id"`
	// job:<kind> | api:<METHOD> <path_template> | cli:<name>
	Trigger   *string `json:"trigger"`
	RequestID *string `json:"request_id"`
	// permanent | until:<YYYY-MM-DD> | expired | none
	ContentState  string `json:"content_state"`
	AccountNotice bool   `json:"account_notice"`
}

// MetaLogPage corresponde ao schema MetaLogPage do contrato.
type MetaLogPage struct {
	Data       []MetaLogEntry `json:"data"`
	NextCursor *string        `json:"next_cursor"`
}

// MetaLogResult corresponde ao schema MetaLogResultOut do contrato.
type MetaLogResult struct {
	// Valores: "ok", "warning", "error", "info".
	Tone string `json:"tone"`
	// Short pt-BR result, e.g. '400 · erro 131031' or '200 · envio bloqueado'
	Label string `json:"label"`
}

// NumberLimit corresponde ao schema NumberLimitOut do contrato.
type NumberLimit struct {
	Limit int64 `json:"limit"`
	// Valores: "trial", "seats", "no_active_plan", "trial_available".
	Reason string `json:"reason"`
}

// NumberUsage corresponde ao schema NumberUsageOut do contrato.
type NumberUsage struct {
	ChannelID          string  `json:"channel_id"`
	DisplayPhoneNumber *string `json:"display_phone_number"`
	OutboundMessages   int64   `json:"outbound_messages"`
	InboundMessages    int64   `json:"inbound_messages"`
}

// OptInParams corresponde ao schema OptInIn do contrato.
type OptInParams struct {
	// Where the contact gave consent.
	// Valores: "checkout", "website_form", "whatsapp", "in_person", "other".
	Source string `json:"source"`
	// When the consent was given (ISO 8601 with timezone).
	CapturedAt time.Time `json:"captured_at"`
	// Optional link to the consent evidence (http/https).
	EvidenceURL string `json:"evidence_url,omitempty"`
}

// OptOutParams corresponde ao schema OptOutIn do contrato.
type OptOutParams struct {
	BusinessID string `json:"business_id"`
	// Contact phone in E.164.
	Phone string `json:"phone"`
}

// OptOut corresponde ao schema OptOutOut do contrato.
type OptOut struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	// Last 4 digits (Meshio keeps only an HMAC of the number).
	PhoneSuffix string `json:"phone_suffix"`
	// Valores: "keyword", "api", "panel".
	Source string `json:"source"`
	// Keyword received (e.g. SAIR) when source=keyword.
	Keyword   *string   `json:"keyword,omitempty"`
	ChannelID *string   `json:"channel_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// OptOutPage corresponde ao schema OptOutPage do contrato.
type OptOutPage struct {
	Data       []OptOut `json:"data"`
	NextCursor *string  `json:"next_cursor,omitempty"`
}

// OverageParams corresponde ao schema OverageIn do contrato.
type OverageParams struct {
	Enabled bool `json:"enabled"`
	// limit per cycle, in micros
	LimitMicros string `json:"limit_micros,omitempty"`
}

// Overage corresponde ao schema OverageOut do contrato.
type Overage struct {
	// null = overage unavailable
	PriceMicros *string `json:"price_micros"`
	Available   bool    `json:"available"`
	// Valores: "price_undefined", "trial", "feature_disabled", "no_subscription".
	UnavailableReason *string `json:"unavailable_reason"`
	Enabled           bool    `json:"enabled"`
	LimitMicros       *string `json:"limit_micros"`
	UsedMicros        string  `json:"used_micros"`
}

// OverviewChannels corresponde ao schema OverviewChannelsOut do contrato.
type OverviewChannels struct {
	Total int64 `json:"total"`
	// Every channel status, zero-filled.
	ByStatus map[string]int64 `json:"by_status"`
	Active   int64            `json:"active"`
	// Channels in `attention_required` or `revoked`.
	WithProblems int64 `json:"with_problems"`
}

// OverviewMessages corresponde ao schema OverviewMessagesOut do contrato.
type OverviewMessages struct {
	// Outbound (api/panel) accepted by Meta or later (delivered/read).
	Sent     int64 `json:"sent"`
	Received int64 `json:"received"`
	Failed   int64 `json:"failed"`
	Unknown  int64 `json:"unknown"`
	Pending  int64 `json:"pending"`
}

// Overview corresponde ao schema OverviewOut do contrato.
type Overview struct {
	BusinessID *string        `json:"business_id"`
	Period     OverviewPeriod `json:"period"`
	// Valor: "meshio".
	DataSource string           `json:"data_source,omitempty"`
	Channels   OverviewChannels `json:"channels"`
	Messages   OverviewMessages `json:"messages"`
	Webhooks   OverviewWebhooks `json:"webhooks"`
	// Same as `webhooks` (alias).
	WebhookDeliveries OverviewWebhooks `json:"webhook_deliveries"`
}

// OverviewPeriod corresponde ao schema OverviewPeriodOut do contrato.
type OverviewPeriod struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// OverviewWebhooks corresponde ao schema OverviewWebhooksOut do contrato.
type OverviewWebhooks struct {
	DeliveriesFailed int64 `json:"deliveries_failed"`
	// Alias of `deliveries_failed`.
	Failed int64 `json:"failed"`
}

// PastDue corresponde ao schema PastDueOut do contrato.
type PastDue struct {
	InvoiceID      string    `json:"invoice_id"`
	AmountMicros   string    `json:"amount_micros"`
	Since          time.Time `json:"since"`
	GraceUntil     time.Time `json:"grace_until"`
	SendingBlocked bool      `json:"sending_blocked"`
}

// PayInvoiceParams corresponde ao schema PayInvoiceIn do contrato.
type PayInvoiceParams struct {
	// Valores: "link", "wallet".
	Method string `json:"method,omitempty"`
	// required to pay with the preserved paid balance without a prior opt-in
	AcceptWalletTerms *bool `json:"accept_wallet_terms,omitempty"`
}

// PaymentMethodParams corresponde ao schema PaymentMethodIn do contrato.
type PaymentMethodParams struct {
	// Valores: "invoice_link", "wallet", "card_auto_renewal", "pix_automatic".
	PaymentMethod     string `json:"payment_method"`
	AcceptWalletTerms *bool  `json:"accept_wallet_terms,omitempty"`
}

// Pricing corresponde ao schema PricingOut do contrato.
type Pricing struct {
	// Valor: "BRL".
	Currency           string     `json:"currency,omitempty"`
	PlanVersion        *int64     `json:"plan_version"`
	EffectiveFrom      *time.Time `json:"effective_from"`
	Tiers              []Tier     `json:"tiers"`
	MessagesPerNumber  int64      `json:"messages_per_number"`
	OveragePriceMicros *string    `json:"overage_price_micros"`
	Complimentary      bool       `json:"complimentary"`
	// a custom condition (override) applies to this integrator
	CustomConditions bool       `json:"custom_conditions"`
	Trial            TrialTerms `json:"trial"`
	GraceDays        int64      `json:"grace_days"`
	RenewalLeadDays  int64      `json:"renewal_lead_days"`
	// Valor: "True".
	MetaFeesSeparate bool `json:"meta_fees_separate,omitempty"`
}

// PromotionalLot corresponde ao schema PromotionalLotOut do contrato.
type PromotionalLot struct {
	RemainingMicros string     `json:"remaining_micros"`
	ExpiresAt       *time.Time `json:"expires_at"`
	Usable          bool       `json:"usable"`
	Label           string     `json:"label"`
}

// RefundRequest corresponde ao schema RefundRequestOut do contrato.
type RefundRequest struct {
	ID      string `json:"id"`
	TopupID string `json:"topup_id"`
	// Valores: "requested", "executed", "rejected".
	Status           string     `json:"status"`
	Reason           string     `json:"reason"`
	CreatedAt        time.Time  `json:"created_at"`
	ResolvedAt       *time.Time `json:"resolved_at"`
	ResolutionReason *string    `json:"resolution_reason"`
}

// RejectedTemplate corresponde ao schema RejectedTemplateOut do contrato.
type RejectedTemplate struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Reason *string `json:"reason"`
}

// Renewal corresponde ao schema RenewalOut do contrato.
type Renewal struct {
	At              time.Time `json:"at"`
	Seats           int64     `json:"seats"`
	UnitPriceMicros string    `json:"unit_price_micros"`
	TotalMicros     string    `json:"total_micros"`
	// renewal invoice already issued (renewal_lead_days before)
	InvoiceID *string `json:"invoice_id"`
}

// SeatAddition corresponde ao schema SeatAdditionOut do contrato.
type SeatAddition struct {
	Add      int64 `json:"add"`
	NewSeats int64 `json:"new_seats"`
	Tier     Tier  `json:"tier"`
	// price of the new seats (tier of the new total)
	UnitPriceMicros        string `json:"unit_price_micros"`
	RemainingDays          int64  `json:"remaining_days"`
	CycleDays              int64  `json:"cycle_days"`
	ProratedAmountMicros   string `json:"prorated_amount_micros"`
	AddedMessages          int64  `json:"added_messages"`
	IncludedMessagesAfter  int64  `json:"included_messages_after"`
	RenewalSeats           int64  `json:"renewal_seats"`
	RenewalUnitPriceMicros string `json:"renewal_unit_price_micros"`
	RenewalTotalMicros     string `json:"renewal_total_micros"`
}

// SeatReduction corresponde ao schema SeatReductionOut do contrato.
type SeatReduction struct {
	Target           int64 `json:"target"`
	CurrentSeats     int64 `json:"current_seats"`
	ConnectedNumbers int64 `json:"connected_numbers"`
	Allowed          bool  `json:"allowed"`
	// Valores: "target_not_lower", "too_many_connected_numbers", "seats_invalid".
	BlockReason            *string `json:"block_reason"`
	Tier                   Tier    `json:"tier"`
	RenewalUnitPriceMicros string  `json:"renewal_unit_price_micros"`
	RenewalTotalMicros     string  `json:"renewal_total_micros"`
	// the reduction applies at the next renewal (no credit)
	EffectiveAt *time.Time `json:"effective_at"`
}

// SeatsChange corresponde ao schema SeatsChangeOut do contrato.
type SeatsChange struct {
	Subscription Subscription `json:"subscription"`
	Invoice      *Invoice     `json:"invoice"`
	// the payment link could not be created now; retry with POST /v1/billing/invoices/{id}/pay
	// Valor: "provider_unavailable".
	PaymentLinkError *string        `json:"payment_link_error,omitempty"`
	Addition         *SeatAddition  `json:"addition,omitempty"`
	Reduction        *SeatReduction `json:"reduction,omitempty"`
}

// SeatsParams corresponde ao schema SeatsIn do contrato.
type SeatsParams struct {
	// immediate, prorated invoice
	Add *int64 `json:"add,omitempty"`
	// reduction at the next renewal
	Target                   *int64 `json:"target,omitempty"`
	CancelScheduledReduction *bool  `json:"cancel_scheduled_reduction,omitempty"`
}

// SeatsPreviewParams corresponde ao schema SeatsPreviewIn do contrato.
type SeatsPreviewParams struct {
	Add    *int64 `json:"add,omitempty"`
	Target *int64 `json:"target,omitempty"`
}

// SeatsPreview corresponde ao schema SeatsPreviewOut do contrato.
type SeatsPreview struct {
	Addition  *SeatAddition  `json:"addition"`
	Reduction *SeatReduction `json:"reduction"`
}

// SendMessageParams corresponde ao schema SendMessageIn do contrato.
type SendMessageParams struct {
	BusinessID string `json:"business_id,omitempty"`
	ChannelID  string `json:"channel_id"`
	// Recipient in E.164 (e.g. +5511988887777).
	To string `json:"to"`
	// Valores: "text", "template".
	Type            string          `json:"type"`
	Text            *TextParams     `json:"text,omitempty"`
	Template        *TemplateParams `json:"template,omitempty"`
	ClientReference string          `json:"client_reference,omitempty"`
	// Contact consent for business-initiated messages (templates). Required for templates when no consent is registered (`POST /v1/consents`) for this business and phone; creates the consent record when the send is accepted. Ignored for free text.
	OptIn *OptInParams `json:"opt_in,omitempty"`
}

// SendMessage corresponde ao schema SendMessageOut do contrato.
type SendMessage struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	ChannelID  string `json:"channel_id"`
	// Valores: "outbound", "inbound".
	Direction string `json:"direction"`
	// Valores: "api", "panel", "business_app", "customer", "unknown".
	Origin string `json:"origin"`
	// `interactive_reply` (inbound only): the contact tapped a template quick reply button or answered an interactive message; see `reply`. New types may be added.
	// Valores: "text", "template", "interactive_reply", "unsupported".
	Type              string `json:"type"`
	CounterpartyPhone string `json:"counterparty_phone"`
	// Text content; for `interactive_reply`, the tapped button/row title.
	BodyText *string          `json:"body_text,omitempty"`
	Template *MessageTemplate `json:"template,omitempty"`
	// Click data for `interactive_reply`.
	Reply *MessageReply `json:"reply,omitempty"`
	// Message the contact replied to (quoted reply or button click).
	Context         *MessageContext `json:"context,omitempty"`
	ClientReference *string         `json:"client_reference,omitempty"`
	// Valores: "queued", "sending", "accepted_by_meta", "delivered", "read", "failed", "canceled", "unknown", "received".
	Status        string        `json:"status"`
	MetaMessageID *string       `json:"meta_message_id,omitempty"`
	Error         *MessageError `json:"error,omitempty"`
	Attempts      int64         `json:"attempts"`
	// True for messages imported from the coexistence history sync.
	Imported    bool    `json:"imported"`
	ContactName *string `json:"contact_name,omitempty"`
	// User (`usr_…`) or API key (`key_…`) that sent it (outbound only).
	CreatedBy *string `json:"created_by,omitempty"`
	// Human-readable sender; filled in `GET /v1/messages/{id}` only.
	CreatedByDisplay *string         `json:"created_by_display,omitempty"`
	MetaTimestamp    *time.Time      `json:"meta_timestamp,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	StatusHistory    []StatusHistory `json:"status_history,omitempty"`
	RequestID        *string         `json:"request_id,omitempty"`
	IdempotentReplay bool            `json:"idempotent_replay,omitempty"`
}

// Sending corresponde ao schema SendingOut do contrato.
type Sending struct {
	Allowed bool `json:"allowed"`
	// why new sends are blocked (null = allowed): trial_ended | no_plan | subscription_canceled | subscription_incomplete | trial_available | subscription_past_due
	// Valores: "trial_ended", "no_plan", "subscription_canceled", "subscription_incomplete", "trial_available", "subscription_past_due".
	Reason *string `json:"reason"`
}

// SignupSessionCreated corresponde ao schema SignupSessionCreatedOut do contrato.
type SignupSessionCreated struct {
	SignupSessionID string `json:"signup_session_id"`
	BusinessID      string `json:"business_id"`
	// Anti-CSRF nonce; shown only on creation. Send it back on complete.
	State string `json:"state"`
	// Valores: "cloud_api", "coexistence".
	Mode            string    `json:"mode"`
	ConfigID        string    `json:"config_id"`
	AppID           string    `json:"app_id"`
	GraphAPIVersion string    `json:"graph_api_version"`
	Launch          Launch    `json:"launch"`
	ExpiresAt       time.Time `json:"expires_at"`
	// Valores: "simulator", "live".
	Environment string `json:"environment"`
}

// SignupSession corresponde ao schema SignupSessionOut do contrato.
type SignupSession struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	// Valores: "cloud_api", "coexistence".
	Mode string `json:"mode"`
	// Valores: "started", "completed", "failed", "canceled", "expired".
	Status      string     `json:"status"`
	ErrorCode   *string    `json:"error_code"`
	Reason      *string    `json:"reason"`
	ChannelID   *string    `json:"channel_id"`
	Channel     *Channel   `json:"channel,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// StatusEndpoint corresponde ao schema StatusEndpointOut do contrato.
type StatusEndpoint struct {
	ID  string `json:"id"`
	URL string `json:"url"`
	// Valores: "ok", "failing".
	State       string     `json:"state"`
	LastErrorAt *time.Time `json:"last_error_at"`
}

// StatusHistory corresponde ao schema StatusHistoryOut do contrato.
type StatusHistory struct {
	Status         string         `json:"status"`
	Source         string         `json:"source"`
	ReceivedAt     time.Time      `json:"received_at"`
	MetaOccurredAt *time.Time     `json:"meta_occurred_at,omitempty"`
	Detail         map[string]any `json:"detail,omitempty"`
}

// StatusNumber corresponde ao schema StatusNumberOut do contrato.
type StatusNumber struct {
	ChannelID          string  `json:"channel_id"`
	BusinessID         string  `json:"business_id"`
	BusinessName       string  `json:"business_name"`
	DisplayPhoneNumber *string `json:"display_phone_number"`
	// Valores: "cloud_api", "coexistence".
	Mode string `json:"mode"`
	// Meta health_status.can_send_message; null until the first reading
	// Valores: "available", "limited", "blocked".
	CanSend *string `json:"can_send"`
	// Meta's reason when limited or blocked (verbatim)
	Reason         *string    `json:"reason"`
	QualityRating  *string    `json:"quality_rating"`
	NameStatus     *string    `json:"name_status"`
	MessagingLimit *string    `json:"messaging_limit"`
	CheckedAt      *time.Time `json:"checked_at"`
}

// SubscriptionChangeResult corresponde ao schema SubscriptionChangeResultOut do contrato.
type SubscriptionChangeResult struct {
	Subscription Subscription `json:"subscription"`
	Invoice      *Invoice     `json:"invoice"`
	// the payment link could not be created now; retry with POST /v1/billing/invoices/{id}/pay
	// Valor: "provider_unavailable".
	PaymentLinkError *string `json:"payment_link_error,omitempty"`
}

// Subscription corresponde ao schema SubscriptionOut do contrato.
type Subscription struct {
	// none = never had a plan (trial still available); trial = free trial open; trial_ended = trial over without a subscription; incomplete = awaiting the first payment; same values as `billing.subscription_status` in GET /v1/status/integration
	// Valores: "none", "trial", "trial_ended", "incomplete", "active", "past_due", "canceled".
	Status         string  `json:"status"`
	SubscriptionID *string `json:"subscription_id"`
	Seats          *int64  `json:"seats"`
	// seat reduction scheduled for the next renewal
	ScheduledSeats     *int64        `json:"scheduled_seats"`
	ConnectedNumbers   int64         `json:"connected_numbers"`
	NumberLimit        NumberLimit   `json:"number_limit"`
	Tier               *Tier         `json:"tier"`
	UnitPriceMicros    *string       `json:"unit_price_micros"`
	MonthlyTotalMicros *string       `json:"monthly_total_micros"`
	MessagesPerNumber  *int64        `json:"messages_per_number"`
	CurrentCycle       *CycleSummary `json:"current_cycle"`
	NextRenewalAt      *time.Time    `json:"next_renewal_at"`
	Renewal            *Renewal      `json:"renewal"`
	CancelAtPeriodEnd  bool          `json:"cancel_at_period_end"`
	CanceledAt         *time.Time    `json:"canceled_at"`
	// Valores: "invoice_link", "wallet".
	PaymentMethod   *string        `json:"payment_method"`
	WalletConsent   *WalletConsent `json:"wallet_consent"`
	Overage         Overage        `json:"overage"`
	PastDue         *PastDue       `json:"past_due"`
	Trial           *Trial         `json:"trial"`
	TrialAvailable  bool           `json:"trial_available"`
	Sending         Sending        `json:"sending"`
	OpenInvoices    []Invoice      `json:"open_invoices"`
	BillingEnabled  bool           `json:"billing_enabled"`
	PaymentsEnabled bool           `json:"payments_enabled"`
	// false = no published plan: contracting is unavailable (503)
	PlanAvailable bool `json:"plan_available"`
}

// SubscriptionPreviewParams corresponde ao schema SubscriptionPreviewIn do contrato.
type SubscriptionPreviewParams struct {
	Seats int64 `json:"seats"`
}

// SubscriptionPreview corresponde ao schema SubscriptionPreviewOut do contrato.
type SubscriptionPreview struct {
	Seats              int64  `json:"seats"`
	Tier               Tier   `json:"tier"`
	UnitPriceMicros    string `json:"unit_price_micros"`
	MonthlyTotalMicros string `json:"monthly_total_micros"`
	// seats × messages_per_number, shared by all numbers
	IncludedMessages   int64   `json:"included_messages"`
	MessagesPerNumber  int64   `json:"messages_per_number"`
	OveragePriceMicros *string `json:"overage_price_micros"`
	ConnectedNumbers   int64   `json:"connected_numbers"`
	PlanVersion        *int64  `json:"plan_version"`
	Complimentary      bool    `json:"complimentary"`
}

// TemplateButtonParams corresponde ao schema TemplateButtonIn do contrato.
type TemplateButtonParams struct {
	// Valores: "quick_reply", "url".
	Type string `json:"type"`
	Text string `json:"text"`
	URL  string `json:"url,omitempty"`
}

// TemplateButton corresponde ao schema TemplateButtonOut do contrato.
type TemplateButton struct {
	// Valores: "quick_reply", "url".
	Type string  `json:"type"`
	Text string  `json:"text"`
	URL  *string `json:"url,omitempty"`
}

// TemplateDelete corresponde ao schema TemplateDeleteOut do contrato.
type TemplateDelete struct {
	ID string `json:"id"`
	// Valores: "hard", "soft".
	Deleted string `json:"deleted"`
	// Final status for soft deletes; null when removed
	// Valores: "draft", "pending_review", "pending", "approved", "rejected", "paused", "disabled", "deleted".
	Status *string `json:"status"`
	// Whether Meta confirmed the deletion (best-effort)
	MetaDeleted bool `json:"meta_deleted"`
}

// TemplateParams corresponde ao schema TemplateIn do contrato.
type TemplateParams struct {
	Name      string                   `json:"name"`
	Language  string                   `json:"language"`
	Variables *TemplateVariablesParams `json:"variables,omitempty"`
}

// Template corresponde ao schema TemplateOut do contrato.
type Template struct {
	ID             string  `json:"id"`
	BusinessID     string  `json:"business_id"`
	ChannelID      *string `json:"channel_id"`
	WABAID         string  `json:"waba_id"`
	MetaTemplateID *string `json:"meta_template_id"`
	Name           string  `json:"name"`
	Language       string  `json:"language"`
	Category       string  `json:"category"`
	// Valores: "draft", "pending_review", "pending", "approved", "rejected", "paused", "disabled", "deleted".
	Status          string  `json:"status"`
	RejectionReason *string `json:"rejection_reason"`
	// True only when Meta approved the template (status=approved)
	Usable        bool               `json:"usable"`
	BodyText      *string            `json:"body_text"`
	BodyExamples  []string           `json:"body_examples"`
	HeaderText    *string            `json:"header_text"`
	HeaderExample *string            `json:"header_example"`
	FooterText    *string            `json:"footer_text"`
	Buttons       []TemplateButton   `json:"buttons"`
	Parameters    TemplateParameters `json:"parameters"`
	Components    []map[string]any   `json:"components"`
	CreatedBy     *string            `json:"created_by"`
	// Human-readable author (user name or API key name + prefix)
	CreatedByDisplay *string    `json:"created_by_display,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	LastSyncedAt     *time.Time `json:"last_synced_at"`
}

// TemplatePage corresponde ao schema TemplatePage do contrato.
type TemplatePage struct {
	Data       []Template `json:"data"`
	NextCursor *string    `json:"next_cursor"`
}

// TemplateParameters corresponde ao schema TemplateParametersOut do contrato.
type TemplateParameters struct {
	BodyCount   int64 `json:"body_count"`
	HeaderCount int64 `json:"header_count"`
	// False when the template needs parameters Meshio cannot send yet
	Sendable bool `json:"sendable"`
}

// TemplateRef corresponde ao schema TemplateRefOut do contrato.
type TemplateRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TemplateSync corresponde ao schema TemplateSyncOut do contrato.
type TemplateSync struct {
	JobID     string `json:"job_id"`
	ChannelID string `json:"channel_id"`
	// Valor: "queued".
	Status string `json:"status,omitempty"`
}

// TemplateVariablesParams corresponde ao schema TemplateVariablesIn do contrato.
type TemplateVariablesParams struct {
	Body   []string `json:"body,omitempty"`
	Header []string `json:"header,omitempty"`
}

// TemplateVariables corresponde ao schema TemplateVariablesOut do contrato.
type TemplateVariables struct {
	Body   []string `json:"body"`
	Header []string `json:"header,omitempty"`
}

// TemplatesStatus corresponde ao schema TemplatesStatusOut do contrato.
type TemplatesStatus struct {
	Approved int64              `json:"approved"`
	Paused   []TemplateRef      `json:"paused"`
	Rejected []RejectedTemplate `json:"rejected"`
}

// TestEvent corresponde ao schema TestEventOut do contrato.
type TestEvent struct {
	EventID    string `json:"event_id"`
	DeliveryID string `json:"delivery_id"`
	// Valor: "queued".
	Status string `json:"status,omitempty"`
}

// TextParams corresponde ao schema TextIn do contrato.
type TextParams struct {
	Body string `json:"body"`
}

// Tier corresponde ao schema TierOut do contrato.
type Tier struct {
	MinNumbers int64 `json:"min_numbers"`
	// null = open-ended (last tier)
	MaxNumbers *int64 `json:"max_numbers"`
	// price per contracted number per month, for ALL numbers
	UnitPriceMicros string `json:"unit_price_micros"`
}

// TopError corresponde ao schema TopErrorOut do contrato.
type TopError struct {
	// Meta error code when known (e.g. 131047), otherwise the Meshio error code
	Code  string `json:"code"`
	Title string `json:"title"`
	Count int64  `json:"count"`
}

// Topup corresponde ao schema TopupOut do contrato.
type Topup struct {
	ID string `json:"id"`
	// refund_processing = refund sent to the provider, awaiting its confirmation
	// Valores: "pending", "paid", "failed", "canceled", "refund_processing", "refunded", "charged_back".
	Status       string `json:"status"`
	AmountCents  int64  `json:"amount_cents"`
	AmountMicros string `json:"amount_micros"`
	// Valor: "BRL".
	Currency string `json:"currency"`
	// Valores: "pix", "card".
	Method  *string `json:"method"`
	Gateway string  `json:"gateway"`
	// true = simulated provider (dev/test), no real money
	Simulated       bool    `json:"simulated"`
	CheckoutURL     *string `json:"checkout_url"`
	BonusCampaignID *string `json:"bonus_campaign_id"`
	// bonus quoted at creation; granted only after confirmed payment
	BonusQuotedMicros string     `json:"bonus_quoted_micros"`
	CreatedAt         time.Time  `json:"created_at"`
	PaidAt            *time.Time `json:"paid_at"`
	FailedAt          *time.Time `json:"failed_at"`
	CanceledAt        *time.Time `json:"canceled_at"`
	RefundedAt        *time.Time `json:"refunded_at"`
	ChargedBackAt     *time.Time `json:"charged_back_at"`
	FailureReason     *string    `json:"failure_reason"`
	// amount refunded (unused paid balance of this top-up)
	RefundedCents int64 `json:"refunded_cents"`
	// Valores: "provider", "manual_pix".
	RefundRoute *string `json:"refund_route"`
	// refund window and request (paid top-ups; integrator view)
	Refund *TopupRefund `json:"refund,omitempty"`
}

// TopupPage corresponde ao schema TopupPage do contrato.
type TopupPage struct {
	Data       []Topup `json:"data"`
	NextCursor *string `json:"next_cursor"`
}

// TopupRefund corresponde ao schema TopupRefundOut do contrato.
type TopupRefund struct {
	// refund may be requested until this moment
	Deadline *time.Time `json:"deadline"`
	Eligible bool       `json:"eligible"`
	// Valores: "topup_not_paid", "refund_processing", "refund_window_expired", "nothing_to_refund".
	Block *string `json:"block"`
	// unused paid balance of this top-up, in whole cents
	EstimatedRefundMicros string         `json:"estimated_refund_micros"`
	Request               *RefundRequest `json:"request"`
}

// Trial corresponde ao schema TrialOut do contrato.
type Trial struct {
	CycleID      string    `json:"cycle_id"`
	Active       bool      `json:"active"`
	PeriodStart  time.Time `json:"period_start"`
	PeriodEnd    time.Time `json:"period_end"`
	DaysLeft     int64     `json:"days_left"`
	MaxNumbers   int64     `json:"max_numbers"`
	MessageLimit int64     `json:"message_limit"`
	UsedMessages int64     `json:"used_messages"`
}

// TrialTerms corresponde ao schema TrialTermsOut do contrato.
type TrialTerms struct {
	Days         int64 `json:"days"`
	MaxNumbers   int64 `json:"max_numbers"`
	MessageLimit int64 `json:"message_limit"`
}

// UpdateWebhookEndpointParams corresponde ao schema UpdateWebhookEndpointIn do contrato.
type UpdateWebhookEndpointParams struct {
	URL         string   `json:"url,omitempty"`
	EventTypes  []string `json:"event_types,omitempty"`
	Description string   `json:"description,omitempty"`
	// Valores: "enabled", "disabled".
	Status string `json:"status,omitempty"`
}

// UsageFilters corresponde ao schema UsageFiltersOut do contrato.
type UsageFilters struct {
	BusinessID *string `json:"business_id"`
	ChannelID  *string `json:"channel_id"`
}

// UsageNote corresponde ao schema UsageNoteOut do contrato.
type UsageNote struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Usage corresponde ao schema UsageOut do contrato.
type Usage struct {
	Period      map[string]any              `json:"period"`
	Filters     UsageFilters                `json:"filters"`
	Series      []UsagePoint                `json:"series"`
	Totals      map[string]int64            `json:"totals"`
	Estimate    bool                        `json:"estimate,omitempty"`
	Sources     map[string]map[string]int64 `json:"sources"`
	Notes       []UsageNote                 `json:"notes"`
	GeneratedAt time.Time                   `json:"generated_at"`
}

// UsagePoint corresponde ao schema UsagePointOut do contrato.
type UsagePoint struct {
	Date    string           `json:"date"`
	Metrics map[string]int64 `json:"metrics"`
}

// ValidationErrorItem corresponde ao schema ValidationError do contrato.
type ValidationErrorItem struct {
	Loc   []any          `json:"loc"`
	Msg   string         `json:"msg"`
	Type  string         `json:"type"`
	Input any            `json:"input,omitempty"`
	Ctx   map[string]any `json:"ctx,omitempty"`
}

// Version corresponde ao schema VersionOut do contrato.
type Version struct {
	// API package version, e.g. 0.13.0
	Version *string `json:"version"`
	// Deployed git ref (git describe), e.g. v0.13.0; null in development
	Ref *string `json:"ref"`
	// Deployed commit (12-char sha = image tag); null in development
	Commit *string `json:"commit"`
	// Deploy time (UTC); null in development
	ReleasedAt *time.Time `json:"released_at"`
}

// WalletConsent corresponde ao schema WalletConsentOut do contrato.
type WalletConsent struct {
	AcceptedAt   time.Time `json:"accepted_at"`
	AcceptedBy   *string   `json:"accepted_by"`
	TermsVersion *string   `json:"terms_version"`
}

// Wallet corresponde ao schema WalletOut do contrato.
type Wallet struct {
	// Valor: "BRL".
	Currency string `json:"currency,omitempty"`
	// paid balance: may pay invoices (opt-in) and is refundable
	PaidBalanceMicros string `json:"paid_balance_micros"`
	// valid bonus/promotional balance: never pays invoices
	PromotionalBalanceMicros string           `json:"promotional_balance_micros"`
	PromotionalLots          []PromotionalLot `json:"promotional_lots"`
	// negative paid balance (e.g. after a chargeback)
	NegativeDueMicros   string `json:"negative_due_micros"`
	CardPaymentsEnabled bool   `json:"card_payments_enabled"`
	OpenDisputes        int64  `json:"open_disputes"`
	// the subscription is set to pay invoices with the paid balance
	PaysInvoices  bool           `json:"pays_invoices"`
	WalletConsent *WalletConsent `json:"wallet_consent"`
	// current version of the pay-with-balance terms
	TermsVersion string `json:"terms_version"`
}

// WebhookDelivery corresponde ao schema WebhookDeliveryOut do contrato.
type WebhookDelivery struct {
	ID         string `json:"id"`
	EndpointID string `json:"endpoint_id"`
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	// Valores: "pending", "retrying", "succeeded", "failed".
	Status string `json:"status"`
	// Total de tentativas, somando ciclos de replay.
	AttemptCount int64 `json:"attempt_count"`
	// Tentativas no ciclo atual (replay reinicia).
	CycleAttemptCount int64 `json:"cycle_attempt_count"`
	// Máximo de tentativas por ciclo.
	MaxAttempts        int64      `json:"max_attempts"`
	ReplayCount        int64      `json:"replay_count"`
	LastResponseStatus *int64     `json:"last_response_status"`
	LastError          *string    `json:"last_error"`
	LastDurationMs     *int64     `json:"last_duration_ms"`
	FirstAttemptAt     *time.Time `json:"first_attempt_at"`
	LastAttemptAt      *time.Time `json:"last_attempt_at"`
	NextAttemptAt      *time.Time `json:"next_attempt_at"`
	CreatedAt          time.Time  `json:"created_at"`
}

// WebhookDeliveryPage corresponde ao schema WebhookDeliveryPage do contrato.
type WebhookDeliveryPage struct {
	Data       []WebhookDelivery `json:"data"`
	NextCursor *string           `json:"next_cursor"`
}

// WebhookEndpoint corresponde ao schema WebhookEndpointOut do contrato.
type WebhookEndpoint struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Description *string  `json:"description"`
	EventTypes  []string `json:"event_types"`
	// Valores: "enabled", "disabled".
	Status    string    `json:"status"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Última rotação do segredo (o valor do segredo nunca é retornado após criação/rotação).
	SecretRotatedAt *time.Time `json:"secret_rotated_at"`
	// Durante a sobreposição de rotação, entregas levam 2 assinaturas v1= até este instante.
	PreviousSecretExpiresAt *time.Time `json:"previous_secret_expires_at"`
	// Delivery aggregate of the last 24 hours (list/get responses).
	Stats24h *DeliveryStats `json:"stats_24h,omitempty"`
}

// WebhookEndpointPage corresponde ao schema WebhookEndpointPage do contrato.
type WebhookEndpointPage struct {
	Data       []WebhookEndpoint `json:"data"`
	NextCursor *string           `json:"next_cursor"`
}

// WebhookEndpointWithSecret corresponde ao schema WebhookEndpointWithSecretOut do contrato.
type WebhookEndpointWithSecret struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Description *string  `json:"description"`
	EventTypes  []string `json:"event_types"`
	// Valores: "enabled", "disabled".
	Status    string    `json:"status"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Última rotação do segredo (o valor do segredo nunca é retornado após criação/rotação).
	SecretRotatedAt *time.Time `json:"secret_rotated_at"`
	// Durante a sobreposição de rotação, entregas levam 2 assinaturas v1= até este instante.
	PreviousSecretExpiresAt *time.Time `json:"previous_secret_expires_at"`
	// Delivery aggregate of the last 24 hours (list/get responses).
	Stats24h *DeliveryStats `json:"stats_24h,omitempty"`
	// Segredo `whsec_...` exibido UMA única vez. Guarde-o agora.
	Secret string `json:"secret"`
}

// Webhooks24h corresponde ao schema Webhooks24hOut do contrato.
type Webhooks24h struct {
	Delivered int64 `json:"delivered"`
	Failed    int64 `json:"failed"`
	// delivered / (delivered + failed); null without deliveries
	SuccessRatio *float64         `json:"success_ratio"`
	Retrying     int64            `json:"retrying"`
	Endpoints    []StatusEndpoint `json:"endpoints"`
}
