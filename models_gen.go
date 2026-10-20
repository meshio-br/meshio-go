// Código gerado por internal/gen/genmodels.py a partir do contrato. NÃO EDITE.

package meshio

import "time"

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

// BalanceStatus corresponde ao schema BalanceStatusOut do contrato.
type BalanceStatus struct {
	FeatureEnabled bool `json:"feature_enabled"`
	// available BRL balance in micros (R$ 1 = 1000000)
	AvailableMicros  *string `json:"available_micros"`
	NumbersInGrace   int64   `json:"numbers_in_grace"`
	NumbersSuspended int64   `json:"numbers_suspended"`
}

// BillingSettingsParams corresponde ao schema BillingSettingsIn do contrato.
type BillingSettingsParams struct {
	// integer string in micros (millionths of BRL): R$ 1 = 1000000, R$ 0.01 = 10000
	LowBalanceThresholdMicros string `json:"low_balance_threshold_micros"`
}

// BonusOfferList corresponde ao schema BonusOfferListOut do contrato.
type BonusOfferList struct {
	Data []BonusOffer `json:"data"`
}

// BonusOffer corresponde ao schema BonusOfferOut do contrato.
type BonusOffer struct {
	CampaignID     string     `json:"campaign_id"`
	Name           string     `json:"name"`
	DisplayText    string     `json:"display_text"`
	StartsAt       time.Time  `json:"starts_at"`
	EndsAt         *time.Time `json:"ends_at"`
	MinTopupMicros string     `json:"min_topup_micros"`
	// Valores: "fixed", "percent".
	BonusType        string  `json:"bonus_type"`
	BonusFixedMicros string  `json:"bonus_fixed_micros"`
	BonusPercentBp   int64   `json:"bonus_percent_bp"`
	BonusCapMicros   *string `json:"bonus_cap_micros"`
	// bonus at the minimum top-up
	ExampleBonusMicros string `json:"example_bonus_micros"`
	// ESTIMATE: sent messages the example bonus pays at the current outbound price
	MessagesEquivalent *int64 `json:"messages_equivalent"`
	RemainingUses      *int64 `json:"remaining_uses"`
	BonusValidityDays  *int64 `json:"bonus_validity_days"`
	FirstTopupOnly     bool   `json:"first_topup_only"`
}

// BusinessAllowance corresponde ao schema BusinessAllowanceOut do contrato.
type BusinessAllowance struct {
	// Allowance in micros (integer string); null = no limit
	MonthlyMicros *string `json:"monthly_micros"`
	// Usage in the current month (integer string in micros): charged messages plus sends still awaiting delivery
	UsedMicros string `json:"used_micros"`
	// Current period YYYY-MM (America/Sao_Paulo)
	Period string `json:"period"`
	// floor(used / monthly * 100); null without allowance; 100 when the allowance is 0
	Percent *int64 `json:"percent"`
	// none = no allowance; warning >= 80%; reached >= 100% (outbound sends refused)
	// Valores: "none", "ok", "warning", "reached".
	Status string `json:"status"`
	// First day of the next month (YYYY-MM-DD), when usage resets
	ResetsAt string `json:"resets_at"`
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
	ID              string            `json:"id"`
	IntegratorID    string            `json:"integrator_id"`
	Name            string            `json:"name"`
	LegalName       *string           `json:"legal_name"`
	TaxID           *string           `json:"tax_id"`
	Country         string            `json:"country"`
	Timezone        string            `json:"timezone"`
	ContactEmail    string            `json:"contact_email"`
	ContactPhone    *string           `json:"contact_phone"`
	Purpose         string            `json:"purpose"`
	TermsAcceptedAt time.Time         `json:"terms_accepted_at"`
	TermsAcceptedBy string            `json:"terms_accepted_by"`
	Status          string            `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Allowance       BusinessAllowance `json:"allowance"`
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

// ChannelBilling corresponde ao schema ChannelBillingOut do contrato.
type ChannelBilling struct {
	// grace: fee overdue, still sending until grace_until; suspended: outbound sends blocked
	// Valores: "ok", "grace", "suspended".
	State string `json:"state"`
	// end of the paid/open cycle = next charge
	CurrentPeriodEnd *time.Time `json:"current_period_end"`
	GraceUntil       *time.Time `json:"grace_until"`
	// effective monthly fee for this number in micros (0 = free)
	MonthlyFeeMicros string `json:"monthly_fee_micros"`
	// always null: the price origin is internal to Meshio operators
	// Valores: "default", "custom".
	PriceSource *string `json:"price_source"`
	// overdue monthly fees (past_due) for this number, in micros
	AmountDueMicros string `json:"amount_due_micros"`
}

// ChannelFee corresponde ao schema ChannelFeeOut do contrato.
type ChannelFee struct {
	ID                 string  `json:"id"`
	ChannelID          string  `json:"channel_id"`
	BusinessID         string  `json:"business_id"`
	BusinessName       *string `json:"business_name"`
	DisplayPhoneNumber *string `json:"display_phone_number"`
	// Valores: "cloud_api", "coexistence".
	Mode string `json:"mode"`
	// Valores: "cycle", "mode_change".
	Kind            string    `json:"kind"`
	PeriodStart     time.Time `json:"period_start"`
	PeriodEnd       time.Time `json:"period_end"`
	AmountMicros    string    `json:"amount_micros"`
	UnitPriceMicros string    `json:"unit_price_micros"`
	// Valores: "default", "custom".
	PriceSource *string `json:"price_source"`
	// Valores: "charged", "past_due", "exempt", "waived".
	Status      string     `json:"status"`
	ChargedAt   *time.Time `json:"charged_at"`
	Attempts    int64      `json:"attempts"`
	WaivedAt    *time.Time `json:"waived_at"`
	WaiveReason *string    `json:"waive_reason"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ChannelFeePage corresponde ao schema ChannelFeePage do contrato.
type ChannelFeePage struct {
	Data       []ChannelFee `json:"data"`
	NextCursor *string      `json:"next_cursor"`
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
	Environment string         `json:"environment"`
	Billing     ChannelBilling `json:"billing"`
	MetaHealth  MetaHealth     `json:"meta_health"`
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
	Period Period `json:"period"`
	// Valores: "business", "channel", "kind", "day".
	GroupBy       string           `json:"group_by"`
	Rows          []ConsumptionRow `json:"rows"`
	TotalMicros   string           `json:"total_micros"`
	TotalQuantity int64            `json:"total_quantity"`
}

// ConsumptionRow corresponde ao schema ConsumptionRowOut do contrato.
type ConsumptionRow struct {
	Key   string  `json:"key"`
	Label *string `json:"label"`
	// Valores: "outbound_message", "inbound_message", "number_fee".
	Kind         string `json:"kind"`
	Quantity     int64  `json:"quantity"`
	AmountMicros string `json:"amount_micros"`
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

// EffectivePrice corresponde ao schema EffectivePriceOut do contrato.
type EffectivePrice struct {
	// Valores: "outbound_message", "inbound_message", "number_fee_cloud_api", "number_fee_coexistence".
	Field string `json:"field"`
	// integer string in micros (millionths of BRL): R$ 1 = 1000000, R$ 0.01 = 10000
	Micros string `json:"micros"`
	// exact decimal BRL, e.g. 0.005
	Brl string `json:"brl"`
	// operator views only (null in the integrator API)
	// Valores: "default", "custom".
	Source *string `json:"source"`
	// operator views only (null in the integrator API)
	DefaultMicros *string `json:"default_micros"`
}

// EffectivePricing corresponde ao schema EffectivePricingOut do contrato.
type EffectivePricing struct {
	// false = no default table published (nothing is charged)
	Available bool `json:"available"`
	// Valor: "BRL".
	Currency        string           `json:"currency,omitempty"`
	DefaultVersion  *int64           `json:"default_version"`
	OverrideVersion *int64           `json:"override_version"`
	Prices          []EffectivePrice `json:"prices"`
}

// Entry corresponde ao schema EntryOut do contrato.
type Entry struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	// Valores: "topup", "topup_refund", "topup_chargeback", "bonus", "bonus_revoked", "bonus_expired", "outbound_message", "inbound_message", "number_fee", "adjustment", "opening_balance", "chargeback_reversal", "closure_refund".
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
	Balance        BalanceStatus   `json:"balance"`
	Attention      []Attention     `json:"attention"`
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

// ModeTotal corresponde ao schema ModeTotalOut do contrato.
type ModeTotal struct {
	Count     int64  `json:"count"`
	FeeMicros string `json:"fee_micros"`
	// Valores: "default", "custom".
	Source string `json:"source"`
}

// NumberFeeSummary corresponde ao schema NumberFeeSummaryOut do contrato.
type NumberFeeSummary struct {
	FeatureEnabled bool `json:"feature_enabled"`
	BillingExempt  bool `json:"billing_exempt"`
	// Valor: "monthly".
	Cycle              string               `json:"cycle,omitempty"`
	ActiveNumbers      int64                `json:"active_numbers"`
	MonthlyTotalMicros string               `json:"monthly_total_micros"`
	ByMode             map[string]ModeTotal `json:"by_mode"`
	NextChargeAt       *time.Time           `json:"next_charge_at"`
	NextChargeMicros   string               `json:"next_charge_micros"`
	PastDueMicros      string               `json:"past_due_micros"`
	SuspendedChannels  int64                `json:"suspended_channels"`
	ChannelsInGrace    int64                `json:"channels_in_grace"`
	RunwayDays         *int64               `json:"runway_days"`
	Numbers            []NumberLine         `json:"numbers"`
}

// NumberLine corresponde ao schema NumberLineOut do contrato.
type NumberLine struct {
	ChannelID          string  `json:"channel_id"`
	BusinessID         string  `json:"business_id"`
	BusinessName       *string `json:"business_name"`
	DisplayPhoneNumber *string `json:"display_phone_number"`
	VerifiedName       *string `json:"verified_name"`
	// Valores: "cloud_api", "coexistence".
	Mode             string `json:"mode"`
	MonthlyFeeMicros string `json:"monthly_fee_micros"`
	// Valores: "ok", "grace", "suspended".
	State           string     `json:"state"`
	NextChargeAt    *time.Time `json:"next_charge_at"`
	GraceUntil      *time.Time `json:"grace_until"`
	AmountDueMicros string     `json:"amount_due_micros"`
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

// Period corresponde ao schema PeriodOut do contrato.
type Period struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// RefundRequestParams corresponde ao schema RefundRequestIn do contrato.
type RefundRequestParams struct {
	Reason string `json:"reason"`
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

// Settings corresponde ao schema SettingsOut do contrato.
type Settings struct {
	ID                            string    `json:"id"`
	Version                       int64     `json:"version"`
	MinTopupCents                 int64     `json:"min_topup_cents"`
	TrialDays                     int64     `json:"trial_days"`
	TrialBonusMicros              string    `json:"trial_bonus_micros"`
	TrialMaxNumbers               int64     `json:"trial_max_numbers"`
	TrialBonusExpires             bool      `json:"trial_bonus_expires"`
	TrialWarningDays              int64     `json:"trial_warning_days"`
	RefundWindowDays              int64     `json:"refund_window_days"`
	ChargebackSuspensionThreshold int64     `json:"chargeback_suspension_threshold"`
	Reason                        string    `json:"reason"`
	CreatedBy                     *string   `json:"created_by"`
	CreatedAt                     time.Time `json:"created_at"`
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

// TopError corresponde ao schema TopErrorOut do contrato.
type TopError struct {
	// Meta error code when known (e.g. 131047), otherwise the Meshio error code
	Code  string `json:"code"`
	Title string `json:"title"`
	Count int64  `json:"count"`
}

// TopupParams corresponde ao schema TopupIn do contrato.
type TopupParams struct {
	// BRL cents (R$ 100,00 = 10000); minimum = setting min_topup_cents
	AmountCents int64 `json:"amount_cents"`
	// Optional; omit it. The payment page accepts Pix and credit card (single payment), or only Pix when card payments are blocked for the account. The top-up method stays null until paid and is then the method actually used at the checkout.
	// Valores: "pix", "card".
	Method string `json:"method,omitempty"`
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
	Exists bool `json:"exists"`
	Active bool `json:"active"`
	// Valores: "active", "ended", "canceled".
	Status                *string    `json:"status"`
	StartedAt             *time.Time `json:"started_at"`
	EndsAt                *time.Time `json:"ends_at"`
	Days                  *int64     `json:"days"`
	DaysLeft              int64      `json:"days_left"`
	BonusGrantedMicros    string     `json:"bonus_granted_micros"`
	BonusRemainingMicros  string     `json:"bonus_remaining_micros"`
	MaxNumbers            *int64     `json:"max_numbers"`
	NumbersUsed           int64      `json:"numbers_used"`
	NumberFeesDuringTrial bool       `json:"number_fees_during_trial,omitempty"`
	SettingsVersion       *int64     `json:"settings_version"`
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

// Wallet corresponde ao schema WalletOut do contrato.
type Wallet struct {
	// Valor: "BRL".
	Currency string `json:"currency,omitempty"`
	// total = paid + bonus
	BalanceMicros      string `json:"balance_micros"`
	PaidBalanceMicros  string `json:"paid_balance_micros"`
	BonusBalanceMicros string `json:"bonus_balance_micros"`
	// reserved for accepted sends awaiting delivery
	ReservedMicros string `json:"reserved_micros"`
	// charges owed awaiting balance (e.g. inbound without balance)
	PendingMicros string `json:"pending_micros"`
	// balance - reserved - pending + credit_limit
	AvailableMicros           string `json:"available_micros"`
	CreditLimitMicros         string `json:"credit_limit_micros"`
	LowBalanceThresholdMicros string `json:"low_balance_threshold_micros"`
	BillingExempt             bool   `json:"billing_exempt"`
	// Valores: "ok", "low", "depleted", "negative", "exempt".
	Status               string `json:"status"`
	FeatureEnabled       bool   `json:"feature_enabled"`
	PaymentsEnabled      bool   `json:"payments_enabled"`
	MonthlyNumbersMicros string `json:"monthly_numbers_micros"`
	PastDueMicros        string `json:"past_due_micros"`
	// false after chargebacks/fraud: top-ups only by Pix (card disabled for this integrator)
	CardPaymentsEnabled bool `json:"card_payments_enabled"`
	// payment disputes (chargebacks) under review
	OpenDisputes int64 `json:"open_disputes"`
	// negative paid balance (e.g. after a chargeback) settled first by the next top-up
	NegativeDueMicros string `json:"negative_due_micros"`
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
