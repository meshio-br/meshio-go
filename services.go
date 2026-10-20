package meshio

import (
	"context"
	"net/http"
	"net/url"
)

// Cada método corresponde a um operationId do contrato (indicado no comentário).

func doJSON[T any](ctx context.Context, c *Client, cl call, opts []RequestOption) (*T, error) {
	var out T
	if err := c.do(ctx, cl, &out, opts); err != nil {
		return nil, err
	}
	return &out, nil
}

func getJSON[T any](ctx context.Context, c *Client, path string, q url.Values, opts []RequestOption) (*T, error) {
	return doJSON[T](ctx, c, call{method: http.MethodGet, path: path, query: q}, opts)
}

func postJSON[T any](ctx context.Context, c *Client, path string, body any, opts []RequestOption) (*T, error) {
	return doJSON[T](ctx, c, call{method: http.MethodPost, path: path, body: body}, opts)
}

func getCSV(ctx context.Context, c *Client, path string, q url.Values, opts []RequestOption) ([]byte, error) {
	var out []byte
	if err := c.do(ctx, call{method: http.MethodGet, path: path, query: q, raw: true}, &out, opts); err != nil {
		return nil, err
	}
	return out, nil
}

type cursorParams interface{ cursorValue() string }

func cursorParam(p cursorParams) string { return p.cursorValue() }

// ---------------------------------------------------------------- version

// VersionService consulta a versão da API.
type VersionService struct{ client *Client }

// Get devolve a versão implantada (version_get).
func (s *VersionService) Get(ctx context.Context, opts ...RequestOption) (*Version, error) {
	return getJSON[Version](ctx, s.client, "/v1/version", nil, opts)
}

// ---------------------------------------------------------------- businesses

// BusinessesService gerencia as empresas (clientes finais) do integrador.
type BusinessesService struct{ client *Client }

// List devolve uma página de empresas (businesses_list).
func (s *BusinessesService) List(ctx context.Context, params *BusinessesListParams, opts ...RequestOption) (*BusinessPage, error) {
	return getJSON[BusinessPage](ctx, s.client, "/v1/businesses", params.query(), opts)
}

// ListAll percorre todas as empresas.
func (s *BusinessesService) ListAll(ctx context.Context, params *BusinessesListParams, opts ...RequestOption) *Iter[Business] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]Business, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Get devolve uma empresa (businesses_get).
func (s *BusinessesService) Get(ctx context.Context, businessID string, opts ...RequestOption) (*Business, error) {
	return getJSON[Business](ctx, s.client, ep("/v1/businesses/", businessID), nil, opts)
}

// ---------------------------------------------------------------- channels

// ChannelsService gerencia os canais (números de WhatsApp).
type ChannelsService struct {
	client         *Client
	SignupSessions *SignupSessionsService
	MetaLogs       *MetaLogsService
}

// List devolve uma página de canais (channels_list).
func (s *ChannelsService) List(ctx context.Context, params *ChannelsListParams, opts ...RequestOption) (*ChannelPage, error) {
	return getJSON[ChannelPage](ctx, s.client, "/v1/channels", params.query(), opts)
}

// ListAll percorre todos os canais.
func (s *ChannelsService) ListAll(ctx context.Context, params *ChannelsListParams, opts ...RequestOption) *Iter[Channel] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]Channel, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Get devolve um canal (channels_get).
func (s *ChannelsService) Get(ctx context.Context, channelID string, opts ...RequestOption) (*Channel, error) {
	return getJSON[Channel](ctx, s.client, ep("/v1/channels/", channelID), nil, opts)
}

// History devolve uma página do histórico de status do canal (channels_history).
func (s *ChannelsService) History(ctx context.Context, channelID string, params *ChannelsHistoryParams, opts ...RequestOption) (*ChannelHistoryPage, error) {
	return getJSON[ChannelHistoryPage](ctx, s.client, ep("/v1/channels/", channelID, "/history"), params.query(), opts)
}

// HistoryAll percorre todo o histórico do canal.
func (s *ChannelsService) HistoryAll(ctx context.Context, channelID string, params *ChannelsHistoryParams, opts ...RequestOption) *Iter[ChannelHistoryItem] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]ChannelHistoryItem, *string, error) {
		pg, err := s.History(ctx, channelID, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// HealthCheck pede uma verificação de saúde do canal (channels_health_check).
func (s *ChannelsService) HealthCheck(ctx context.Context, channelID string, opts ...RequestOption) (*Channel, error) {
	return postJSON[Channel](ctx, s.client, ep("/v1/channels/", channelID, "/health-check"), nil, opts)
}

// Disconnect desconecta o canal (channels_disconnect).
func (s *ChannelsService) Disconnect(ctx context.Context, channelID string, params *DisconnectParams, opts ...RequestOption) (*Channel, error) {
	if params == nil {
		params = &DisconnectParams{}
	}
	return postJSON[Channel](ctx, s.client, ep("/v1/channels/", channelID, "/disconnect"), params, opts)
}

// SignupSessionsService gerencia sessões de cadastro embutido (Embedded Signup).
type SignupSessionsService struct{ client *Client }

// Create inicia uma sessão de cadastro (signup_sessions_create).
func (s *SignupSessionsService) Create(ctx context.Context, params *CreateSignupSessionParams, opts ...RequestOption) (*SignupSessionCreated, error) {
	return postJSON[SignupSessionCreated](ctx, s.client, "/v1/channels/signup-sessions", params, opts)
}

// Get devolve uma sessão (signup_sessions_get).
func (s *SignupSessionsService) Get(ctx context.Context, sessionID string, opts ...RequestOption) (*SignupSession, error) {
	return getJSON[SignupSession](ctx, s.client, ep("/v1/channels/signup-sessions/", sessionID), nil, opts)
}

// Complete conclui a sessão e cria o canal (signup_sessions_complete).
func (s *SignupSessionsService) Complete(ctx context.Context, sessionID string, params *CompleteSignupParams, opts ...RequestOption) (*Channel, error) {
	return postJSON[Channel](ctx, s.client, ep("/v1/channels/signup-sessions/", sessionID, "/complete"), params, opts)
}

// Cancel cancela a sessão (signup_sessions_cancel).
func (s *SignupSessionsService) Cancel(ctx context.Context, sessionID string, params *CancelSignupParams, opts ...RequestOption) (*SignupSession, error) {
	if params == nil {
		params = &CancelSignupParams{}
	}
	return postJSON[SignupSession](ctx, s.client, ep("/v1/channels/signup-sessions/", sessionID, "/cancel"), params, opts)
}

// MetaLogsService consulta o registro de interações com a Meta por canal.
type MetaLogsService struct{ client *Client }

// List devolve uma página de registros (meta_logs_list).
func (s *MetaLogsService) List(ctx context.Context, channelID string, params *MetaLogsListParams, opts ...RequestOption) (*MetaLogPage, error) {
	return getJSON[MetaLogPage](ctx, s.client, ep("/v1/channels/", channelID, "/meta-logs"), params.query(), opts)
}

// ListAll percorre todos os registros.
func (s *MetaLogsService) ListAll(ctx context.Context, channelID string, params *MetaLogsListParams, opts ...RequestOption) *Iter[MetaLogEntry] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]MetaLogEntry, *string, error) {
		pg, err := s.List(ctx, channelID, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Get devolve um registro detalhado (meta_logs_get).
func (s *MetaLogsService) Get(ctx context.Context, channelID, entryID string, opts ...RequestOption) (*MetaLogEntryDetail, error) {
	return getJSON[MetaLogEntryDetail](ctx, s.client, ep("/v1/channels/", channelID, "/meta-logs/", entryID), nil, opts)
}

// ---------------------------------------------------------------- connect links

// ConnectLinksService gerencia links de conexão (o dono da empresa conecta o número por uma página hospedada).
type ConnectLinksService struct{ client *Client }

// Create cria um link de conexão (connect_links_create).
func (s *ConnectLinksService) Create(ctx context.Context, params *CreateConnectLinkParams, opts ...RequestOption) (*ConnectLink, error) {
	return postJSON[ConnectLink](ctx, s.client, "/v1/connect-links", params, opts)
}

// List devolve uma página de links (connect_links_list).
func (s *ConnectLinksService) List(ctx context.Context, params *ConnectLinksListParams, opts ...RequestOption) (*ConnectLinkPage, error) {
	return getJSON[ConnectLinkPage](ctx, s.client, "/v1/connect-links", params.query(), opts)
}

// ListAll percorre todos os links.
func (s *ConnectLinksService) ListAll(ctx context.Context, params *ConnectLinksListParams, opts ...RequestOption) *Iter[ConnectLink] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]ConnectLink, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Get devolve um link (connect_links_get).
func (s *ConnectLinksService) Get(ctx context.Context, linkID string, opts ...RequestOption) (*ConnectLink, error) {
	return getJSON[ConnectLink](ctx, s.client, ep("/v1/connect-links/", linkID), nil, opts)
}

// Revoke revoga um link (connect_links_revoke).
func (s *ConnectLinksService) Revoke(ctx context.Context, linkID string, opts ...RequestOption) (*ConnectLink, error) {
	return postJSON[ConnectLink](ctx, s.client, ep("/v1/connect-links/", linkID, "/revoke"), nil, opts)
}

// ---------------------------------------------------------------- messages

// MessagesService envia e consulta mensagens.
type MessagesService struct{ client *Client }

// Send envia uma mensagem de texto ou template (messages_send).
// Sempre envia Idempotency-Key: use WithIdempotencyKey ou o SDK gera um UUID v4,
// reutilizado nas novas tentativas. IdempotentReplay indica resposta repetida.
func (s *MessagesService) Send(ctx context.Context, params *SendMessageParams, opts ...RequestOption) (*SendMessage, error) {
	var meta ResponseMeta
	opts = append(opts[:len(opts):len(opts)], WithResponseMeta(&meta))
	out, err := doJSON[SendMessage](ctx, s.client, call{method: http.MethodPost, path: "/v1/messages", body: params, idempotent: true}, opts)
	if err != nil {
		return nil, err
	}
	if meta.IdempotentReplay {
		out.IdempotentReplay = true
	}
	return out, nil
}

// List devolve uma página de mensagens (messages_list).
func (s *MessagesService) List(ctx context.Context, params *MessagesListParams, opts ...RequestOption) (*MessagePage, error) {
	return getJSON[MessagePage](ctx, s.client, "/v1/messages", params.query(), opts)
}

// ListAll percorre todas as mensagens.
func (s *MessagesService) ListAll(ctx context.Context, params *MessagesListParams, opts ...RequestOption) *Iter[Message] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]Message, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Get devolve uma mensagem (messages_get).
func (s *MessagesService) Get(ctx context.Context, messageID string, opts ...RequestOption) (*Message, error) {
	return getJSON[Message](ctx, s.client, ep("/v1/messages/", messageID), nil, opts)
}

// ---------------------------------------------------------------- conversations

// ConversationsService consulta conversas por contato.
type ConversationsService struct{ client *Client }

// List devolve uma página de conversas (conversations_list).
func (s *ConversationsService) List(ctx context.Context, params *ConversationsListParams, opts ...RequestOption) (*ConversationPage, error) {
	return getJSON[ConversationPage](ctx, s.client, "/v1/conversations", params.query(), opts)
}

// ListAll percorre todas as conversas.
func (s *ConversationsService) ListAll(ctx context.Context, params *ConversationsListParams, opts ...RequestOption) *Iter[Conversation] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]Conversation, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Timeline devolve uma página de mensagens com um contato (conversations_timeline).
// counterparty é o telefone em E.164.
func (s *ConversationsService) Timeline(ctx context.Context, channelID, counterparty string, params *ConversationsTimelineParams, opts ...RequestOption) (*MessagePage, error) {
	return getJSON[MessagePage](ctx, s.client, ep("/v1/conversations/", channelID, "/", counterparty), params.query(), opts)
}

// TimelineAll percorre toda a conversa com um contato.
func (s *ConversationsService) TimelineAll(ctx context.Context, channelID, counterparty string, params *ConversationsTimelineParams, opts ...RequestOption) *Iter[Message] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]Message, *string, error) {
		pg, err := s.Timeline(ctx, channelID, counterparty, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// ---------------------------------------------------------------- templates

// TemplatesService gerencia templates de mensagem.
type TemplatesService struct{ client *Client }

// List devolve uma página de templates (templates_list).
func (s *TemplatesService) List(ctx context.Context, params *TemplatesListParams, opts ...RequestOption) (*TemplatePage, error) {
	return getJSON[TemplatePage](ctx, s.client, "/v1/templates", params.query(), opts)
}

// ListAll percorre todos os templates.
func (s *TemplatesService) ListAll(ctx context.Context, params *TemplatesListParams, opts ...RequestOption) *Iter[Template] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]Template, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Create cria um template (templates_create).
func (s *TemplatesService) Create(ctx context.Context, params *CreateTemplateParams, opts ...RequestOption) (*Template, error) {
	return postJSON[Template](ctx, s.client, "/v1/templates", params, opts)
}

// Get devolve um template (templates_get).
func (s *TemplatesService) Get(ctx context.Context, templateID string, opts ...RequestOption) (*Template, error) {
	return getJSON[Template](ctx, s.client, ep("/v1/templates/", templateID), nil, opts)
}

// Delete apaga um template (templates_delete).
func (s *TemplatesService) Delete(ctx context.Context, templateID string, opts ...RequestOption) (*TemplateDelete, error) {
	return doJSON[TemplateDelete](ctx, s.client, call{method: http.MethodDelete, path: ep("/v1/templates/", templateID)}, opts)
}

// Sync sincroniza os templates do canal com a Meta (templates_sync).
func (s *TemplatesService) Sync(ctx context.Context, params *TemplatesSyncParams, opts ...RequestOption) (*TemplateSync, error) {
	return doJSON[TemplateSync](ctx, s.client, call{method: http.MethodPost, path: "/v1/templates/sync", query: params.query()}, opts)
}

// Submit envia o template para aprovação da Meta (templates_submit).
func (s *TemplatesService) Submit(ctx context.Context, templateID string, opts ...RequestOption) (*Template, error) {
	return postJSON[Template](ctx, s.client, ep("/v1/templates/", templateID, "/submit"), nil, opts)
}

// ---------------------------------------------------------------- webhook endpoints

// WebhookEndpointsService gerencia os endpoints que recebem eventos.
type WebhookEndpointsService struct{ client *Client }

// List devolve uma página de endpoints (webhook_endpoints_list).
func (s *WebhookEndpointsService) List(ctx context.Context, params *WebhookEndpointsListParams, opts ...RequestOption) (*WebhookEndpointPage, error) {
	return getJSON[WebhookEndpointPage](ctx, s.client, "/v1/webhook-endpoints", params.query(), opts)
}

// ListAll percorre todos os endpoints.
func (s *WebhookEndpointsService) ListAll(ctx context.Context, params *WebhookEndpointsListParams, opts ...RequestOption) *Iter[WebhookEndpoint] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]WebhookEndpoint, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Create cria um endpoint; o segredo (whsec_…) só é devolvido aqui e na rotação (webhook_endpoints_create).
func (s *WebhookEndpointsService) Create(ctx context.Context, params *CreateWebhookEndpointParams, opts ...RequestOption) (*WebhookEndpointWithSecret, error) {
	return postJSON[WebhookEndpointWithSecret](ctx, s.client, "/v1/webhook-endpoints", params, opts)
}

// Get devolve um endpoint (webhook_endpoints_get).
func (s *WebhookEndpointsService) Get(ctx context.Context, endpointID string, opts ...RequestOption) (*WebhookEndpoint, error) {
	return getJSON[WebhookEndpoint](ctx, s.client, ep("/v1/webhook-endpoints/", endpointID), nil, opts)
}

// Update altera um endpoint (webhook_endpoints_update).
func (s *WebhookEndpointsService) Update(ctx context.Context, endpointID string, params *UpdateWebhookEndpointParams, opts ...RequestOption) (*WebhookEndpoint, error) {
	return doJSON[WebhookEndpoint](ctx, s.client, call{method: http.MethodPatch, path: ep("/v1/webhook-endpoints/", endpointID), body: params}, opts)
}

// Delete apaga um endpoint (webhook_endpoints_delete).
func (s *WebhookEndpointsService) Delete(ctx context.Context, endpointID string, opts ...RequestOption) error {
	return s.client.do(ctx, call{method: http.MethodDelete, path: ep("/v1/webhook-endpoints/", endpointID)}, nil, opts)
}

// RotateSecret gera um novo segredo (webhook_endpoints_rotate_secret).
func (s *WebhookEndpointsService) RotateSecret(ctx context.Context, endpointID string, opts ...RequestOption) (*WebhookEndpointWithSecret, error) {
	return postJSON[WebhookEndpointWithSecret](ctx, s.client, ep("/v1/webhook-endpoints/", endpointID, "/rotate-secret"), nil, opts)
}

// Test envia um evento webhook.test ao endpoint (webhook_endpoints_test).
func (s *WebhookEndpointsService) Test(ctx context.Context, endpointID string, opts ...RequestOption) (*TestEvent, error) {
	return postJSON[TestEvent](ctx, s.client, ep("/v1/webhook-endpoints/", endpointID, "/test"), nil, opts)
}

// ---------------------------------------------------------------- webhook deliveries

// WebhookDeliveriesService consulta e reenvia entregas de webhook.
type WebhookDeliveriesService struct{ client *Client }

// List devolve uma página de entregas (webhook_deliveries_list).
func (s *WebhookDeliveriesService) List(ctx context.Context, params *WebhookDeliveriesListParams, opts ...RequestOption) (*WebhookDeliveryPage, error) {
	return getJSON[WebhookDeliveryPage](ctx, s.client, "/v1/webhook-deliveries", params.query(), opts)
}

// ListAll percorre todas as entregas.
func (s *WebhookDeliveriesService) ListAll(ctx context.Context, params *WebhookDeliveriesListParams, opts ...RequestOption) *Iter[WebhookDelivery] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]WebhookDelivery, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Get devolve uma entrega (webhook_deliveries_get).
func (s *WebhookDeliveriesService) Get(ctx context.Context, deliveryID string, opts ...RequestOption) (*WebhookDelivery, error) {
	return getJSON[WebhookDelivery](ctx, s.client, ep("/v1/webhook-deliveries/", deliveryID), nil, opts)
}

// Replay reenvia uma entrega (webhook_deliveries_replay).
func (s *WebhookDeliveriesService) Replay(ctx context.Context, deliveryID string, opts ...RequestOption) (*WebhookDelivery, error) {
	return postJSON[WebhookDelivery](ctx, s.client, ep("/v1/webhook-deliveries/", deliveryID, "/replay"), nil, opts)
}

// ---------------------------------------------------------------- events

// EventsService consulta os eventos emitidos.
type EventsService struct{ client *Client }

// List devolve uma página de eventos (events_list).
func (s *EventsService) List(ctx context.Context, params *EventsListParams, opts ...RequestOption) (*EventPage, error) {
	return getJSON[EventPage](ctx, s.client, "/v1/events", params.query(), opts)
}

// ListAll percorre todos os eventos.
func (s *EventsService) ListAll(ctx context.Context, params *EventsListParams, opts ...RequestOption) *Iter[Event] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]Event, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Get devolve um evento (events_get).
func (s *EventsService) Get(ctx context.Context, eventID string, opts ...RequestOption) (*Event, error) {
	return getJSON[Event](ctx, s.client, ep("/v1/events/", eventID), nil, opts)
}

// ---------------------------------------------------------------- usage

// UsageService consulta uso, visão geral e status da integração.
type UsageService struct{ client *Client }

// Get devolve o uso agregado (usage_get).
func (s *UsageService) Get(ctx context.Context, params *UsageGetParams, opts ...RequestOption) (*Usage, error) {
	return getJSON[Usage](ctx, s.client, "/v1/usage", params.query(), opts)
}

// ExportCSV devolve o uso em CSV (usage_export_csv).
func (s *UsageService) ExportCSV(ctx context.Context, params *UsageExportCSVParams, opts ...RequestOption) ([]byte, error) {
	return getCSV(ctx, s.client, "/v1/usage/export.csv", params.query(), opts)
}

// Overview devolve a visão geral (overview_get).
func (s *UsageService) Overview(ctx context.Context, params *OverviewGetParams, opts ...RequestOption) (*Overview, error) {
	return getJSON[Overview](ctx, s.client, "/v1/overview", params.query(), opts)
}

// IntegrationStatus devolve o status da integração (status_integration).
func (s *UsageService) IntegrationStatus(ctx context.Context, opts ...RequestOption) (*IntegrationStatus, error) {
	return getJSON[IntegrationStatus](ctx, s.client, "/v1/status/integration", nil, opts)
}

// ---------------------------------------------------------------- compliance

// ComplianceService consulta a situação de conformidade.
type ComplianceService struct{ client *Client }

// Status devolve a situação de conformidade (compliance_status).
func (s *ComplianceService) Status(ctx context.Context, opts ...RequestOption) (*ComplianceStatus, error) {
	return getJSON[ComplianceStatus](ctx, s.client, "/v1/compliance/status", nil, opts)
}

// ---------------------------------------------------------------- consents

// ConsentsService registra e consulta consentimentos (opt-in) de contatos.
type ConsentsService struct{ client *Client }

// List devolve uma página de consentimentos (consents_list).
func (s *ConsentsService) List(ctx context.Context, params *ConsentsListParams, opts ...RequestOption) (*ConsentPage, error) {
	return getJSON[ConsentPage](ctx, s.client, "/v1/consents", params.query(), opts)
}

// ListAll percorre todos os consentimentos.
func (s *ConsentsService) ListAll(ctx context.Context, params *ConsentsListParams, opts ...RequestOption) *Iter[Consent] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]Consent, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Create registra consentimentos em lote (consents_create).
func (s *ConsentsService) Create(ctx context.Context, params *ConsentBulkParams, opts ...RequestOption) (*ConsentBulk, error) {
	return postJSON[ConsentBulk](ctx, s.client, "/v1/consents", params, opts)
}

// Delete revoga um consentimento (consents_delete).
func (s *ConsentsService) Delete(ctx context.Context, consentID string, opts ...RequestOption) (*Consent, error) {
	return doJSON[Consent](ctx, s.client, call{method: http.MethodDelete, path: ep("/v1/consents/", consentID)}, opts)
}

// ---------------------------------------------------------------- opt-outs

// OptOutsService gerencia descadastros (opt-out) de contatos.
type OptOutsService struct{ client *Client }

// List devolve uma página de opt-outs (opt_outs_list).
func (s *OptOutsService) List(ctx context.Context, params *OptOutsListParams, opts ...RequestOption) (*OptOutPage, error) {
	return getJSON[OptOutPage](ctx, s.client, "/v1/opt-outs", params.query(), opts)
}

// ListAll percorre todos os opt-outs.
func (s *OptOutsService) ListAll(ctx context.Context, params *OptOutsListParams, opts ...RequestOption) *Iter[OptOut] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]OptOut, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Create registra um opt-out (opt_outs_create).
func (s *OptOutsService) Create(ctx context.Context, params *OptOutParams, opts ...RequestOption) (*OptOut, error) {
	return postJSON[OptOut](ctx, s.client, "/v1/opt-outs", params, opts)
}

// Delete remove o opt-out de um contato (opt_outs_delete).
func (s *OptOutsService) Delete(ctx context.Context, params *OptOutsDeleteParams, opts ...RequestOption) error {
	return s.client.do(ctx, call{method: http.MethodDelete, path: "/v1/opt-outs", query: params.query()}, nil, opts)
}

// ---------------------------------------------------------------- billing

// BillingService consulta carteira, preços, lançamentos e recargas.
type BillingService struct {
	client     *Client
	Entries    *BillingEntriesService
	NumberFees *BillingNumberFeesService
	Settings   *BillingSettingsService
	Topups     *TopupsService
}

// Wallet devolve a carteira (billing_wallet).
func (s *BillingService) Wallet(ctx context.Context, opts ...RequestOption) (*Wallet, error) {
	return getJSON[Wallet](ctx, s.client, "/v1/billing/wallet", nil, opts)
}

// Pricing devolve os preços efetivos (billing_pricing).
func (s *BillingService) Pricing(ctx context.Context, opts ...RequestOption) (*EffectivePricing, error) {
	return getJSON[EffectivePricing](ctx, s.client, "/v1/billing/pricing", nil, opts)
}

// Consumption devolve o consumo agregado (billing_consumption).
func (s *BillingService) Consumption(ctx context.Context, params *BillingConsumptionParams, opts ...RequestOption) (*Consumption, error) {
	return getJSON[Consumption](ctx, s.client, "/v1/billing/consumption", params.query(), opts)
}

// Trial devolve a situação do período de teste (billing_trial).
func (s *BillingService) Trial(ctx context.Context, opts ...RequestOption) (*Trial, error) {
	return getJSON[Trial](ctx, s.client, "/v1/billing/trial", nil, opts)
}

// BonusOffers lista as ofertas de bônus (billing_bonus_offers).
func (s *BillingService) BonusOffers(ctx context.Context, opts ...RequestOption) (*BonusOfferList, error) {
	return getJSON[BonusOfferList](ctx, s.client, "/v1/billing/bonus-offers", nil, opts)
}

// BillingEntriesService consulta o extrato.
type BillingEntriesService struct{ client *Client }

// List devolve uma página de lançamentos (billing_entries_list).
func (s *BillingEntriesService) List(ctx context.Context, params *BillingEntriesListParams, opts ...RequestOption) (*EntryPage, error) {
	return getJSON[EntryPage](ctx, s.client, "/v1/billing/entries", params.query(), opts)
}

// ListAll percorre todos os lançamentos.
func (s *BillingEntriesService) ListAll(ctx context.Context, params *BillingEntriesListParams, opts ...RequestOption) *Iter[Entry] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]Entry, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// ExportCSV devolve os lançamentos em CSV (billing_entries_export_csv).
func (s *BillingEntriesService) ExportCSV(ctx context.Context, params *BillingEntriesExportCSVParams, opts ...RequestOption) ([]byte, error) {
	return getCSV(ctx, s.client, "/v1/billing/entries/export.csv", params.query(), opts)
}

// BillingNumberFeesService consulta as mensalidades por número.
type BillingNumberFeesService struct{ client *Client }

// List devolve uma página de mensalidades (billing_number_fees_list).
func (s *BillingNumberFeesService) List(ctx context.Context, params *BillingNumberFeesListParams, opts ...RequestOption) (*ChannelFeePage, error) {
	return getJSON[ChannelFeePage](ctx, s.client, "/v1/billing/number-fees", params.query(), opts)
}

// ListAll percorre todas as mensalidades.
func (s *BillingNumberFeesService) ListAll(ctx context.Context, params *BillingNumberFeesListParams, opts ...RequestOption) *Iter[ChannelFee] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]ChannelFee, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Summary devolve o resumo das mensalidades (billing_number_fees_summary).
func (s *BillingNumberFeesService) Summary(ctx context.Context, opts ...RequestOption) (*NumberFeeSummary, error) {
	return getJSON[NumberFeeSummary](ctx, s.client, "/v1/billing/number-fees/summary", nil, opts)
}

// BillingSettingsService lê e altera as configurações de cobrança.
type BillingSettingsService struct{ client *Client }

// Get devolve as configurações públicas (billing_settings_get).
func (s *BillingSettingsService) Get(ctx context.Context, opts ...RequestOption) (*Settings, error) {
	return getJSON[Settings](ctx, s.client, "/v1/billing/settings/public", nil, opts)
}

// Update altera as configurações e devolve a carteira (billing_settings_update).
func (s *BillingSettingsService) Update(ctx context.Context, params *BillingSettingsParams, opts ...RequestOption) (*Wallet, error) {
	return doJSON[Wallet](ctx, s.client, call{method: http.MethodPatch, path: "/v1/billing/settings", body: params}, opts)
}

// TopupsService gerencia recargas da carteira.
type TopupsService struct{ client *Client }

// List devolve uma página de recargas (topups_list).
func (s *TopupsService) List(ctx context.Context, params *TopupsListParams, opts ...RequestOption) (*TopupPage, error) {
	return getJSON[TopupPage](ctx, s.client, "/v1/billing/topups", params.query(), opts)
}

// ListAll percorre todas as recargas.
func (s *TopupsService) ListAll(ctx context.Context, params *TopupsListParams, opts ...RequestOption) *Iter[Topup] {
	return newIter(ctx, cursorParam(params), func(ctx context.Context, cur string) ([]Topup, *string, error) {
		pg, err := s.List(ctx, params.withCursor(cur), opts...)
		if err != nil {
			return nil, nil, err
		}
		return pg.Data, pg.NextCursor, nil
	})
}

// Create cria uma recarga (topups_create). Envia Idempotency-Key (gerada se não informada).
func (s *TopupsService) Create(ctx context.Context, params *TopupParams, opts ...RequestOption) (*Topup, error) {
	return doJSON[Topup](ctx, s.client, call{method: http.MethodPost, path: "/v1/billing/topups", body: params, idempotent: true}, opts)
}

// Get devolve uma recarga (topups_get).
func (s *TopupsService) Get(ctx context.Context, topupID string, opts ...RequestOption) (*Topup, error) {
	return getJSON[Topup](ctx, s.client, ep("/v1/billing/topups/", topupID), nil, opts)
}

// RequestRefund pede o reembolso de uma recarga (topups_request_refund).
func (s *TopupsService) RequestRefund(ctx context.Context, topupID string, params *RefundRequestParams, opts ...RequestOption) (*RefundRequest, error) {
	if params == nil {
		params = &RefundRequestParams{}
	}
	return postJSON[RefundRequest](ctx, s.client, ep("/v1/billing/topups/", topupID, "/refund-request"), params, opts)
}
