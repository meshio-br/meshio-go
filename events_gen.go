// Código gerado por internal/gen/genmodels.py a partir do contrato. NÃO EDITE.

package meshio

import (
	"encoding/json"
	"time"
)

// Tipos de evento conhecidos. Novos tipos podem surgir; trate os desconhecidos sem falhar.
const (
	EventTypeBillingAllowanceExhausted  = "billing.allowance_exhausted"
	EventTypeBillingAllowanceWarning    = "billing.allowance_warning"
	EventTypeBillingInvoiceCreated      = "billing.invoice_created"
	EventTypeBillingInvoicePaid         = "billing.invoice_paid"
	EventTypeBillingInvoicePastDue      = "billing.invoice_past_due"
	EventTypeBillingOverageLimitReached = "billing.overage_limit_reached"
	EventTypeBillingSubscriptionUpdated = "billing.subscription_updated"
	EventTypeBillingTrialEnded          = "billing.trial_ended"
	EventTypeBillingTrialEnding         = "billing.trial_ending"
	EventTypeBusinessComplianceHold     = "business.compliance_hold"
	EventTypeChannelAttentionRequired   = "channel.attention_required"
	EventTypeChannelComplianceHold      = "channel.compliance_hold"
	EventTypeChannelConnected           = "channel.connected"
	EventTypeChannelDisconnected        = "channel.disconnected"
	EventTypeChannelHealthChanged       = "channel.health_changed"
	EventTypeContactOptedIn             = "contact.opted_in"
	EventTypeContactOptedOut            = "contact.opted_out"
	EventTypeMessagePolicyBlocked       = "message.policy_blocked"
	EventTypeMessageReceived            = "message.received"
	EventTypeMessageStatusChanged       = "message.status_changed"
	EventTypeTemplateStatusChanged      = "template.status_changed"
	EventTypeWebhookTest                = "webhook.test"
)

// BillingAllowanceExhaustedData é o campo data do evento billing.allowance_exhausted.
// Allowance exhausted (100% of the cycle allowance): new sends are refused with 402 ALLOWANCE_EXHAUSTED unless overage is available and authorized. Same data as billing.allowance.common.schema.json.
type BillingAllowanceExhaustedData struct {
	// Billing cycle id.
	CycleID string `json:"cycle_id"`
	// trial = free trial cycle; paid = subscription cycle.
	// Valores: "trial", "paid".
	CycleKind string `json:"cycle_kind,omitempty"`
	// Messages included in the cycle, shared by all numbers of the integrator (base + prorated seat additions).
	IncludedMessages int64 `json:"included_messages"`
	// Messages counted in the cycle (sent and delivered or read + new inbound).
	UsedMessages int64 `json:"used_messages"`
	// floor(used_messages / included_messages * 100); may exceed 100 (inbound messages are always received and counted).
	Percent int64 `json:"percent"`
	// End of the cycle (renewal or trial end), when the allowance renews.
	CycleEndsAt time.Time `json:"cycle_ends_at"`
}

// BillingAllowanceWarningData é o campo data do evento billing.allowance_warning.
// Allowance warning (80% of the cycle allowance). Same data as billing.allowance.common.schema.json.
type BillingAllowanceWarningData struct {
	// Billing cycle id.
	CycleID string `json:"cycle_id"`
	// trial = free trial cycle; paid = subscription cycle.
	// Valores: "trial", "paid".
	CycleKind string `json:"cycle_kind,omitempty"`
	// Messages included in the cycle, shared by all numbers of the integrator (base + prorated seat additions).
	IncludedMessages int64 `json:"included_messages"`
	// Messages counted in the cycle (sent and delivered or read + new inbound).
	UsedMessages int64 `json:"used_messages"`
	// floor(used_messages / included_messages * 100); may exceed 100 (inbound messages are always received and counted).
	Percent int64 `json:"percent"`
	// End of the cycle (renewal or trial end), when the allowance renews.
	CycleEndsAt time.Time `json:"cycle_ends_at"`
}

// BillingInvoiceCreatedData é o campo data do evento billing.invoice_created.
// An invoice was issued: the monthly fee of a cycle (issued before the renewal), prorated seats added mid-cycle or the overage of a closed cycle. Pay it through checkout_url (payment link) or, with consent, with the preserved paid balance. Envelope business_id and channel_id are null.
type BillingInvoiceCreatedData struct {
	InvoiceID      string `json:"invoice_id"`
	SubscriptionID string `json:"subscription_id,omitempty"`
	// subscription = monthly fee of the cycle; seat_addition = prorated seats added mid-cycle; overage = messages above the allowance (only with a published price and the integrator's authorization).
	// Valores: "subscription", "seat_addition", "overage".
	Kind string `json:"kind"`
	// Invoice total (multiple of 10000: whole cents). Integer string in micros (R$ 1 = "1000000").
	TotalMicros string `json:"total_micros"`
	// Due date.
	DueAt time.Time `json:"due_at"`
	// Payment link (Pix or card); null when the invoice is paid with the preserved balance or the link is not generated yet.
	CheckoutURL *string `json:"checkout_url,omitempty"`
}

// BillingInvoicePaidData é o campo data do evento billing.invoice_paid.
// Payment of an invoice was confirmed (re-read from the payment provider, paid with the preserved balance or registered by Meshio). Envelope business_id and channel_id are null.
type BillingInvoicePaidData struct {
	InvoiceID      string `json:"invoice_id"`
	SubscriptionID string `json:"subscription_id,omitempty"`
	// subscription = monthly fee of the cycle; seat_addition = prorated seats added mid-cycle; overage = messages above the allowance (only with a published price and the integrator's authorization).
	// Valores: "subscription", "seat_addition", "overage".
	Kind string `json:"kind"`
	// Invoice total. Integer string in micros (R$ 1 = "1000000").
	TotalMicros string `json:"total_micros"`
	// Payment confirmation time.
	PaidAt time.Time `json:"paid_at"`
	// wallet = preserved paid balance (opt-in); manual = registered by Meshio.
	// Valores: "pix", "card", "wallet", "manual".
	PaidVia string `json:"paid_via"`
}

// BillingInvoicePastDueData é o campo data do evento billing.invoice_past_due.
// An invoice is past due and unpaid: the subscription becomes past_due. After grace_until, new sends answer 402 SUBSCRIPTION_PAST_DUE until the invoice is paid; inbound messages, delivery events and webhooks never stop. Envelope business_id and channel_id are null.
type BillingInvoicePastDueData struct {
	InvoiceID      string `json:"invoice_id"`
	SubscriptionID string `json:"subscription_id,omitempty"`
	// subscription = monthly fee of the cycle; seat_addition = prorated seats added mid-cycle; overage = messages above the allowance (only with a published price and the integrator's authorization).
	// Valores: "subscription", "seat_addition", "overage".
	Kind string `json:"kind,omitempty"`
	// Invoice total. Integer string in micros (R$ 1 = "1000000").
	TotalMicros string `json:"total_micros"`
	// Due date.
	DueAt time.Time `json:"due_at"`
	// End of the grace period; new sends are blocked from then on.
	GraceUntil time.Time `json:"grace_until"`
}

// BillingOverageLimitReachedData é o campo data do evento billing.overage_limit_reached.
// The authorized overage reached the integrator's BRL limit for the cycle: new sends above the allowance are refused with 402 ALLOWANCE_EXHAUSTED (reason overage_limit_reached). Sent at most once per cycle. Envelope business_id and channel_id are null.
type BillingOverageLimitReachedData struct {
	// Billing cycle id.
	CycleID string `json:"cycle_id"`
	// Overage limit authorized for the cycle. Integer string in micros (R$ 1 = "1000000").
	OverageLimitMicros string `json:"overage_limit_micros"`
	// Overage already counted in the cycle (sends awaiting delivery included). Integer string in micros (R$ 1 = "1000000").
	OverageUsedMicros string `json:"overage_used_micros"`
}

// BillingSubscriptionUpdatedData é o campo data do evento billing.subscription_updated.
// The integrator's subscription changed (activated, seats added, reduction scheduled or canceled, cancellation scheduled or reverted, canceled, overage or payment method updated, past due, reactivated). action names the change; the other fields are the subscription after it. Envelope business_id and channel_id are null.
type BillingSubscriptionUpdatedData struct {
	SubscriptionID string `json:"subscription_id"`
	// Valores: "incomplete", "active", "past_due", "canceled".
	Status string `json:"status"`
	// Contracted seats (billable numbers).
	Seats int64 `json:"seats"`
	// Seat reduction scheduled for the next renewal; null when none.
	ScheduledSeats *int64 `json:"scheduled_seats,omitempty"`
	// End of the current paid cycle; null before the first payment.
	CurrentPeriodEnd  *time.Time `json:"current_period_end,omitempty"`
	CancelAtPeriodEnd bool       `json:"cancel_at_period_end"`
	// Valores: "created", "activated", "seats_added", "seats_reduction_scheduled", "seats_reduction_canceled", "cancel_scheduled", "cancel_reverted", "canceled", "overage_updated", "payment_method_updated", "past_due", "reactivated".
	Action string `json:"action"`
}

// BillingTrialEndedData é o campo data do evento billing.trial_ended.
// The free trial ended. There is no automatic charge: new sends answer 402 SUBSCRIPTION_REQUIRED (reason trial_ended) until the integrator subscribes; inbound messages keep being received. A trial never restarts. Envelope business_id and channel_id are null.
type BillingTrialEndedData struct {
	// Billing cycle id.
	CycleID string `json:"cycle_id,omitempty"`
	// End of the trial.
	EndedAt time.Time `json:"ended_at"`
}

// BillingTrialEndingData é o campo data do evento billing.trial_ending.
// The integrator's free trial ends soon (sent once, 3 days before the end). There is no automatic charge: to keep sending after the trial, subscribe. Envelope business_id and channel_id are null.
type BillingTrialEndingData struct {
	// Billing cycle id.
	CycleID string `json:"cycle_id,omitempty"`
	// End of the trial.
	EndsAt   time.Time `json:"ends_at"`
	DaysLeft int64     `json:"days_left"`
	// Messages included in the trial.
	MessageLimit int64 `json:"message_limit"`
	// Messages counted in the trial so far.
	UsedMessages int64 `json:"used_messages"`
}

// BusinessComplianceHoldData é o campo data do evento business.compliance_hold.
// Meshio suspended (state=active) or released (state=lifted) all sends of a business (customer) for compliance. Automatic holds come from Meta signals: account_update violation/restriction/ban/disable, 3+ templates paused in 7 days or 20+ Meta policy errors (131048, 131049, 368, 470) in 24 h; operators may also hold manually. While active, sends and new templates for the business are refused with 422 POLICY_BLOCKED (reason business_suspended) without calling Meta. Only a Meshio operator can lift a hold.
type BusinessComplianceHoldData struct {
	HoldID string `json:"hold_id"`
	// Valores: "active", "lifted".
	State string `json:"state"`
	// Valor: "business".
	Scope      string `json:"scope"`
	BusinessID string `json:"business_id"`
	// Valores: "meta_account_violation", "meta_account_restriction", "meta_account_banned", "meta_account_disabled", "templates_paused", "meta_policy_errors", "manual".
	Reason string `json:"reason"`
	// Valores: "meta_signal", "operator", "policy".
	Source string `json:"source"`
	// Short explanation in pt-BR.
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	LiftedAt  *time.Time `json:"lifted_at"`
}

// ChannelAttentionRequiredData é o campo data do evento channel.attention_required.
// Channel needs action (revoked token, lost scope, failed health check, quality problem). Sends are blocked until the channel returns to active.
type ChannelAttentionRequiredData struct {
	ChannelID  string `json:"channel_id"`
	BusinessID string `json:"business_id"`
	// Valores: "cloud_api", "coexistence".
	Mode string `json:"mode"`
	// Valores: "not_started", "pending_authorization", "pending_setup", "active", "attention_required", "disconnected", "revoked".
	// Valor: "attention_required".
	Status string `json:"status"`
	// Valores: "not_started", "pending_authorization", "pending_setup", "active", "attention_required", "disconnected", "revoked".
	PreviousStatus *string `json:"previous_status"`
	StatusReason   string  `json:"status_reason"`
	// Who caused the transition.
	// Valores: "user", "meta", "system", "health_check".
	Source string `json:"source,omitempty"`
	// Business phone number as reported by Meta.
	DisplayPhoneNumber *string                                  `json:"display_phone_number,omitempty"`
	VerifiedName       *string                                  `json:"verified_name,omitempty"`
	WABAID             *string                                  `json:"waba_id,omitempty"`
	PhoneNumberID      *string                                  `json:"phone_number_id,omitempty"`
	QualityRating      *string                                  `json:"quality_rating,omitempty"`
	Capabilities       ChannelAttentionRequiredDataCapabilities `json:"capabilities"`
	// Valores: "simulator", "live".
	Environment string `json:"environment"`
}

// ChannelAttentionRequiredDataCapabilities corresponde ao schema do contrato.
type ChannelAttentionRequiredDataCapabilities struct {
	SendText          bool `json:"send_text"`
	SendTemplate      bool `json:"send_template"`
	Receive           bool `json:"receive"`
	BusinessAppEchoes bool `json:"business_app_echoes"`
	HistorySync       bool `json:"history_sync"`
}

// ChannelComplianceHoldData é o campo data do evento channel.compliance_hold.
// Meshio suspended (state=active) or released (state=lifted) sends of one WhatsApp number for compliance. Automatic holds happen when Meta's quality rating of the number turns RED; operators may also hold manually. While active, sends and new templates on the number are refused with 422 POLICY_BLOCKED (reason channel_compliance_hold) without calling Meta. Only a Meshio operator can lift a hold.
type ChannelComplianceHoldData struct {
	HoldID string `json:"hold_id"`
	// Valores: "active", "lifted".
	State string `json:"state"`
	// Valor: "channel".
	Scope              string  `json:"scope"`
	BusinessID         string  `json:"business_id"`
	ChannelID          string  `json:"channel_id"`
	DisplayPhoneNumber *string `json:"display_phone_number"`
	// Valores: "quality_red", "manual".
	Reason string `json:"reason"`
	// Valores: "meta_signal", "operator", "policy".
	Source string `json:"source"`
	// Short explanation in pt-BR.
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	LiftedAt  *time.Time `json:"lifted_at"`
}

// ChannelConnectedData é o campo data do evento channel.connected.
// Emitted only after the backend validated the Meta authorization, subscribed the app to the WABA and verified the phone number.
type ChannelConnectedData struct {
	ChannelID  string `json:"channel_id"`
	BusinessID string `json:"business_id"`
	// Valores: "cloud_api", "coexistence".
	Mode string `json:"mode"`
	// Valores: "not_started", "pending_authorization", "pending_setup", "active", "attention_required", "disconnected", "revoked".
	// Valor: "active".
	Status string `json:"status"`
	// Valores: "not_started", "pending_authorization", "pending_setup", "active", "attention_required", "disconnected", "revoked".
	PreviousStatus *string `json:"previous_status"`
	// Machine-readable reason (e.g. meta_auth_revoked, health_check_failed); null when not applicable.
	StatusReason *string `json:"status_reason"`
	// Who caused the transition.
	// Valores: "user", "meta", "system", "health_check".
	Source             string                           `json:"source,omitempty"`
	DisplayPhoneNumber string                           `json:"display_phone_number"`
	VerifiedName       *string                          `json:"verified_name,omitempty"`
	WABAID             string                           `json:"waba_id"`
	PhoneNumberID      string                           `json:"phone_number_id"`
	QualityRating      *string                          `json:"quality_rating,omitempty"`
	Capabilities       ChannelConnectedDataCapabilities `json:"capabilities"`
	// Valores: "simulator", "live".
	Environment string `json:"environment"`
	// Present when the business owner connected the number through a hosted connection link (POST /v1/connect-links).
	ConnectLinkID string `json:"connect_link_id,omitempty"`
}

// ChannelConnectedDataCapabilities corresponde ao schema do contrato.
type ChannelConnectedDataCapabilities struct {
	SendText          bool `json:"send_text"`
	SendTemplate      bool `json:"send_template"`
	Receive           bool `json:"receive"`
	BusinessAppEchoes bool `json:"business_app_echoes"`
	HistorySync       bool `json:"history_sync"`
}

// ChannelDisconnectedData é o campo data do evento channel.disconnected.
// Channel disconnected by a user/system or revoked on Meta. New sends are rejected with CHANNEL_NOT_ACTIVE. History is not deleted by disconnection.
type ChannelDisconnectedData struct {
	ChannelID  string `json:"channel_id"`
	BusinessID string `json:"business_id"`
	// Valores: "cloud_api", "coexistence".
	Mode string `json:"mode"`
	// Valores: "disconnected", "revoked".
	Status any `json:"status"`
	// Valores: "not_started", "pending_authorization", "pending_setup", "active", "attention_required", "disconnected", "revoked".
	PreviousStatus *string `json:"previous_status"`
	// Machine-readable reason (e.g. meta_auth_revoked, health_check_failed); null when not applicable.
	StatusReason *string `json:"status_reason"`
	// Who caused the transition.
	// Valores: "user", "meta", "system", "health_check".
	Source string `json:"source,omitempty"`
	// Business phone number as reported by Meta.
	DisplayPhoneNumber *string                             `json:"display_phone_number,omitempty"`
	VerifiedName       *string                             `json:"verified_name,omitempty"`
	WABAID             *string                             `json:"waba_id,omitempty"`
	PhoneNumberID      *string                             `json:"phone_number_id,omitempty"`
	QualityRating      *string                             `json:"quality_rating,omitempty"`
	Capabilities       ChannelDisconnectedDataCapabilities `json:"capabilities"`
	// Valores: "simulator", "live".
	Environment string `json:"environment"`
}

// ChannelDisconnectedDataCapabilities corresponde ao schema do contrato.
type ChannelDisconnectedDataCapabilities struct {
	SendText          bool `json:"send_text"`
	SendTemplate      bool `json:"send_template"`
	Receive           bool `json:"receive"`
	BusinessAppEchoes bool `json:"business_app_echoes"`
	HistorySync       bool `json:"history_sync"`
}

// ChannelHealthChangedData é o campo data do evento channel.health_changed.
// Meta health of a WhatsApp number changed: the sending status reported by Meta (health_status.can_send_message) changed, or the quality rating dropped (GREEN to YELLOW/RED, YELLOW to RED). The first reading only emits when the number is already limited or blocked. Read by the channel health check (every 30 minutes and after Meta quality webhooks).
type ChannelHealthChangedData struct {
	ChannelID          string  `json:"channel_id"`
	BusinessID         string  `json:"business_id"`
	DisplayPhoneNumber *string `json:"display_phone_number"`
	// available: meets all requirements; limited: sends with limitations (see reason); blocked: Meta rejects sends from this number. null = not read yet.
	// Valores: "available", "limited", "blocked".
	CanSend *string `json:"can_send"`
	// Valores: "available", "limited", "blocked".
	PreviousCanSend *string `json:"previous_can_send"`
	// First error description (blocked) or additional info (limited) sent by Meta, verbatim (usually English). null when available.
	Reason *string `json:"reason"`
	// Meta quality rating: GREEN, YELLOW, RED, NA or UNKNOWN.
	QualityRating         *string `json:"quality_rating"`
	PreviousQualityRating *string `json:"previous_quality_rating"`
	// Display name review: APPROVED, AVAILABLE_WITHOUT_REVIEW, DECLINED, EXPIRED, PENDING_REVIEW or NONE.
	NameStatus *string `json:"name_status"`
	// Meta messaging limit tier (whatsapp_business_manager_messaging_limit), e.g. TIER_250, TIER_2K, TIER_10K, TIER_100K, TIER_UNLIMITED.
	MessagingLimit *string `json:"messaging_limit"`
	// Since v0.13.0. Explained health: available (no Meta restriction); limited (sends normally up to the messaging limit for business-initiated conversations); initiate_blocked (replies within 24 hours work, templates are refused, e.g. payment method issue); blocked (Meta blocks sending); unknown (health not read yet).
	// Valores: "available", "limited", "initiate_blocked", "blocked", "unknown".
	State string `json:"state,omitempty"`
	// Since v0.13.0. One pt-BR sentence explaining what the number can do now (same text as the panel).
	Headline string `json:"headline,omitempty"`
}

// ContactOptedInData é o campo data do evento contact.opted_in.
// A contact left the business opt-out list: they sent VOLTAR/START on WhatsApp (source=keyword; a consent with source whatsapp is recorded), a consent captured after the opt-out was registered (source=consent), or the integrator removed them via DELETE /v1/opt-outs or the panel (source=api|panel).
type ContactOptedInData struct {
	BusinessID string  `json:"business_id"`
	ChannelID  *string `json:"channel_id,omitempty"`
	// Contact phone in E.164.
	Phone string `json:"phone"`
	// Valores: "keyword", "consent", "api", "panel".
	Source string `json:"source"`
	// VOLTAR or START when source=keyword.
	Keyword *string `json:"keyword,omitempty"`
	// Consent that lifted the opt-out (source=consent).
	ConsentID *string   `json:"consent_id,omitempty"`
	OptedInAt time.Time `json:"opted_in_at"`
}

// ContactOptedOutData é o campo data do evento contact.opted_out.
// A contact joined the business opt-out (suppression) list: they sent SAIR/PARAR/STOP/CANCELAR/DESCADASTRAR/REMOVER on WhatsApp (source=keyword), or the integrator added them via POST /v1/opt-outs or the panel. From now on Meshio refuses templates to this contact from this business (POLICY_BLOCKED recipient_opted_out) and revokes the live consent. Emitted once per entry.
type ContactOptedOutData struct {
	BusinessID string `json:"business_id"`
	// Channel that received the keyword (null for API/panel).
	ChannelID *string `json:"channel_id,omitempty"`
	// Contact phone in E.164.
	Phone string `json:"phone"`
	// Valores: "keyword", "api", "panel".
	Source string `json:"source"`
	// Keyword received, upper case (e.g. SAIR), when source=keyword.
	Keyword    *string   `json:"keyword,omitempty"`
	OptOutID   string    `json:"opt_out_id"`
	OptedOutAt time.Time `json:"opted_out_at"`
}

// MessagePolicyBlockedData é o campo data do evento message.policy_blocked.
// Meshio refused a send before calling Meta (the API answered 422 POLICY_BLOCKED; no message was created and no credits were reserved). Use `reason` (stable code) and `message` (pt-BR) to fix the cause. The content is never included.
type MessagePolicyBlockedData struct {
	DecisionID string `json:"decision_id"`
	BusinessID string `json:"business_id"`
	ChannelID  string `json:"channel_id"`
	// Contact phone in E.164.
	To string `json:"to"`
	// Valores: "text", "template".
	Type     string                            `json:"type"`
	Template *MessagePolicyBlockedDataTemplate `json:"template,omitempty"`
	// New reasons may be added; treat unknown values as a generic policy refusal.
	// Valores: "integrator_not_approved", "business_suspended", "business_not_verified", "channel_compliance_hold", "recipient_opted_out", "opt_in_required", "link_not_allowed", "link_malicious", "content_suspected_scam", "volume_limit_reached".
	Reason string `json:"reason"`
	// Valores: "trust", "opt_out", "opt_in", "links", "scam", "volume".
	Rule string `json:"rule"`
	// Explanation in pt-BR (same as the API error).
	Message string `json:"message"`
	// Structured details, e.g. link_issue and host (links), score and signals (scam), limit/daily_cap/retry_after_seconds (volume).
	Details         map[string]any `json:"details,omitempty"`
	ClientReference *string        `json:"client_reference,omitempty"`
	BlockedAt       time.Time      `json:"blocked_at"`
}

// MessagePolicyBlockedDataTemplate corresponde ao schema do contrato.
type MessagePolicyBlockedDataTemplate struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

// MessageReceivedData é o campo data do evento message.received.
// A message from a contact (direction=inbound, origin=customer) or, on coexistence channels, an echo of a message sent from the WhatsApp Business App (direction=outbound, origin=business_app). Never auto-reply unless direction=inbound AND origin=customer.
type MessageReceivedData struct {
	MessageID string `json:"message_id"`
	// Valores: "inbound", "outbound".
	Direction string `json:"direction"`
	// Valores: "customer", "business_app", "unknown".
	Origin string `json:"origin"`
	// `interactive_reply` (since v0.14.0): the contact tapped a template quick reply button or answered an interactive message; see `reply`. New types may be added; treat unknown values as unsupported.
	// Valores: "text", "template", "interactive_reply", "unsupported".
	Type string `json:"type"`
	// E.164 sender.
	From string `json:"from"`
	// E.164 recipient.
	To   string                   `json:"to"`
	Text *MessageReceivedDataText `json:"text,omitempty"`
	// Click data for type=interactive_reply (null otherwise).
	Reply *MessageReceivedDataReply `json:"reply,omitempty"`
	// Message the contact replied to (Meta context.id): quoted replies and button clicks.
	Context *MessageReceivedDataContext `json:"context,omitempty"`
	// WhatsApp profile name of the contact when provided by Meta.
	ContactName   *string    `json:"contact_name,omitempty"`
	MetaMessageID *string    `json:"meta_message_id,omitempty"`
	MetaTimestamp *time.Time `json:"meta_timestamp,omitempty"`
	// True for messages from coexistence history sync (best effort); do not treat as new.
	Imported bool `json:"imported,omitempty"`
}

// MessageReceivedDataText corresponde ao schema do contrato.
type MessageReceivedDataText struct {
	Body string `json:"body"`
}

// Click data for type=interactive_reply (null otherwise).
type MessageReceivedDataReply struct {
	// template_quick_reply: quick reply button of a template (Meta type=button); button_reply/list_reply: reply to an interactive message.
	// Valores: "template_quick_reply", "button_reply", "list_reply".
	Kind string `json:"kind"`
	// Button payload (template quick reply) or the reply/row id (interactive).
	ID *string `json:"id,omitempty"`
	// Text the contact saw and tapped.
	Title string `json:"title"`
	// List row description (list_reply only).
	Description *string `json:"description,omitempty"`
}

// Message the contact replied to (Meta context.id): quoted replies and button clicks.
type MessageReceivedDataContext struct {
	MetaMessageID string `json:"meta_message_id"`
	// Meshio id when the replied message belongs to this integrator.
	MessageID *string `json:"message_id,omitempty"`
	Type      *string `json:"type,omitempty"`
	// Template of the replied message, when it was a template.
	TemplateName *string `json:"template_name,omitempty"`
}

// MessageStatusChangedData é o campo data do evento message.status_changed.
// Forward-only status change of an outbound message. Ranks: queued 0, sending 1, unknown 1, accepted_by_meta 2, delivered 3, read 4; failed is terminal and may follow any rank <= 3. Events may arrive out of order: never regress.
type MessageStatusChangedData struct {
	MessageID string `json:"message_id"`
	// Valor: "outbound".
	Direction string `json:"direction"`
	// Valores: "api", "panel", "business_app", "unknown".
	Origin string `json:"origin"`
	// Valores: "text", "template", "unsupported".
	Type string `json:"type,omitempty"`
	// Valores: "queued", "sending", "accepted_by_meta", "delivered", "read", "failed", "canceled", "unknown".
	Status string `json:"status"`
	// Valores: "queued", "sending", "accepted_by_meta", "delivered", "read", "failed", "canceled", "unknown".
	PreviousStatus  *string `json:"previous_status"`
	To              string  `json:"to"`
	ClientReference *string `json:"client_reference,omitempty"`
	MetaMessageID   *string `json:"meta_message_id,omitempty"`
	// Time reported by Meta for this status, when available.
	MetaOccurredAt *time.Time                     `json:"meta_occurred_at,omitempty"`
	Error          *MessageStatusChangedDataError `json:"error,omitempty"`
}

// MessageStatusChangedDataError corresponde ao schema do contrato.
type MessageStatusChangedDataError struct {
	// Safe Meshio error code (UPPER_SNAKE_CASE). Values in examples are illustrative; original Meta codes stay in restricted diagnostics.
	Code    string `json:"code"`
	Message string `json:"message"`
}

// TemplateStatusChangedData é o campo data do evento template.status_changed.
// Template status changed on Meta (source of truth). Only status=approved templates can be sent. Since v0.12.0 a template held for Meshio review (pending_review) that a Meshio operator refuses changes to rejected with a rejection_reason starting with 'Recusado pela Meshio: '.
type TemplateStatusChangedData struct {
	TemplateID     string  `json:"template_id"`
	MetaTemplateID *string `json:"meta_template_id,omitempty"`
	Name           string  `json:"name"`
	Language       string  `json:"language"`
	// Valores: "MARKETING", "UTILITY", "AUTHENTICATION".
	Category string `json:"category"`
	// Valores: "draft", "pending_review", "pending", "approved", "rejected", "paused", "disabled", "deleted".
	Status string `json:"status"`
	// Valores: "draft", "pending_review", "pending", "approved", "rejected", "paused", "disabled", "deleted".
	PreviousStatus  *string `json:"previous_status"`
	RejectionReason *string `json:"rejection_reason,omitempty"`
}

// WebhookTestData é o campo data do evento webhook.test.
// Sent by POST /v1/webhook-endpoints/{id}/test to validate reachability and signature verification. Envelope business_id and channel_id are null.
type WebhookTestData struct {
	EndpointID string `json:"endpoint_id"`
	Message    string `json:"message,omitempty"`
}

// decodeEventData decodifica data no tipo do evento; ok=false para tipos desconhecidos.
func decodeEventData(eventType string, raw json.RawMessage) (v any, ok bool, err error) {
	switch eventType {
	case "billing.allowance_exhausted":
		var d BillingAllowanceExhaustedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "billing.allowance_warning":
		var d BillingAllowanceWarningData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "billing.invoice_created":
		var d BillingInvoiceCreatedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "billing.invoice_paid":
		var d BillingInvoicePaidData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "billing.invoice_past_due":
		var d BillingInvoicePastDueData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "billing.overage_limit_reached":
		var d BillingOverageLimitReachedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "billing.subscription_updated":
		var d BillingSubscriptionUpdatedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "billing.trial_ended":
		var d BillingTrialEndedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "billing.trial_ending":
		var d BillingTrialEndingData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "business.compliance_hold":
		var d BusinessComplianceHoldData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "channel.attention_required":
		var d ChannelAttentionRequiredData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "channel.compliance_hold":
		var d ChannelComplianceHoldData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "channel.connected":
		var d ChannelConnectedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "channel.disconnected":
		var d ChannelDisconnectedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "channel.health_changed":
		var d ChannelHealthChangedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "contact.opted_in":
		var d ContactOptedInData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "contact.opted_out":
		var d ContactOptedOutData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "message.policy_blocked":
		var d MessagePolicyBlockedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "message.received":
		var d MessageReceivedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "message.status_changed":
		var d MessageStatusChangedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "template.status_changed":
		var d TemplateStatusChangedData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	case "webhook.test":
		var d WebhookTestData
		err = json.Unmarshal(raw, &d)
		return &d, true, err
	}
	return nil, false, nil
}
