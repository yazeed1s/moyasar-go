package moyasar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// WebhookService provides access to Moyasar webhook registration and delivery
// attempt APIs.
type WebhookService struct {
	client *Client
}

// WebhookEventType is the type of event delivered by Moyasar webhooks.
type WebhookEventType string

const (
	WebhookEventPaymentPaid           WebhookEventType = "payment_paid"
	WebhookEventPaymentFailed         WebhookEventType = "payment_failed"
	WebhookEventPaymentVoided         WebhookEventType = "payment_voided"
	WebhookEventPaymentAuthorized     WebhookEventType = "payment_authorized"
	WebhookEventPaymentCaptured       WebhookEventType = "payment_captured"
	WebhookEventPaymentRefunded       WebhookEventType = "payment_refunded"
	WebhookEventPaymentAbandoned      WebhookEventType = "payment_abandoned"
	WebhookEventPaymentVerified       WebhookEventType = "payment_verified"
	WebhookEventCardAuthAuthenticated WebhookEventType = "card_auth_authenticated"
	WebhookEventCardAuthFailed        WebhookEventType = "card_auth_failed"
)

// Webhook is a registered endpoint that receives Moyasar events.
type Webhook struct {
	// ID is the webhook's unique ID.
	ID string `json:"id"`
	// HTTPMethod is the HTTP method used to deliver events.
	HTTPMethod string `json:"http_method"`
	// URL is the notification URL where Moyasar sends event requests.
	URL string `json:"url"`
	// CreatedAt is the time the webhook registration was created.
	CreatedAt string `json:"created_at"`
	// Events contains the event types this webhook listens to.
	Events []WebhookEventType `json:"events"`
}

// CreateWebhookRequest registers a webhook and the events it should listen to.
//
// Omitting Events creates a global event listener for currently available
// events and events added in the future.
type CreateWebhookRequest struct {
	// HTTPMethod is the HTTP method used to deliver events, such as post.
	HTTPMethod string `json:"http_method"`
	// URL is the notification URL where Moyasar sends event requests.
	URL string `json:"url"`
	// SharedSecret is assigned by the consumer to secure webhook requests.
	SharedSecret string `json:"shared_secret,omitempty"`
	// Events contains the event types this webhook listens to.
	Events []WebhookEventType `json:"events,omitempty"`
}

// WebhookList contains registered webhooks.
type WebhookList struct {
	// Webhooks contains the registered webhook objects.
	Webhooks []Webhook `json:"webhooks"`
}

// AvailableWebhookEvents contains event types that can be configured on a webhook.
type AvailableWebhookEvents struct {
	// Events contains all available webhook event names.
	Events []WebhookEventType `json:"events"`
}

// DeleteWebhookResponse is returned after deleting a webhook.
type DeleteWebhookResponse struct {
	// Message is Moyasar's deletion result message.
	Message string `json:"message"`
}

// WebhookEvent is the object Moyasar sends to webhook endpoints.
type WebhookEvent struct {
	// ID is the event's unique ID.
	ID string `json:"id"`
	// Type is the type of event, such as payment_paid or card_auth_authenticated.
	Type WebhookEventType `json:"type"`
	// CreatedAt is the time the webhook object was created.
	CreatedAt string `json:"created_at"`
	// SecretToken is the endpoint secret assigned by the consumer.
	SecretToken string `json:"secret_token"`
	// AccountName is the name of the account in which the event occurred.
	AccountName string `json:"account_name"`
	// Live is true for live mode events and false for test mode events.
	Live bool `json:"live"`
	// Data is the event payload. Payment events contain a Payment object, and
	// card_auth events contain a CardAuth object.
	Data json.RawMessage `json:"data"`
}

// WebhookAttempt is a delivery attempt for a Moyasar webhook event.
type WebhookAttempt struct {
	// ID is the webhook attempt's unique ID.
	ID string `json:"id"`
	// WebhookID is the webhook registration ID.
	WebhookID string `json:"webhook_id"`
	// EventID is the delivered event ID.
	EventID string `json:"event_id"`
	// EventType is the delivered event type.
	EventType WebhookEventType `json:"event_type"`
	// RetryNumber is the delivery attempt number.
	RetryNumber int `json:"retry_number"`
	// Result is the delivery result, such as success.
	Result string `json:"result"`
	// Message is Moyasar's human-readable delivery result.
	Message string `json:"message"`
	// ResponseCode is the HTTP status returned by the webhook endpoint.
	ResponseCode int `json:"response_code"`
	// ResponseHeaders contains the response headers returned by the endpoint.
	ResponseHeaders string `json:"response_headers"`
	// ResponseBody contains the response body returned by the endpoint.
	ResponseBody string `json:"response_body"`
	// CreatedAt is the time the delivery attempt was created.
	CreatedAt string `json:"created_at"`
}

// ListWebhookAttemptsParams filters webhook delivery attempts.
type ListWebhookAttemptsParams struct {
	// Page is the requested page number.
	Page int
	// WebhookID filters attempts by webhook registration ID.
	WebhookID string
	// EventID filters attempts by event ID.
	EventID string
	// EventType filters attempts by event type.
	EventType WebhookEventType
}

// WebhookAttemptList contains webhook delivery attempts.
type WebhookAttemptList struct {
	// Attempts contains the returned delivery attempts.
	Attempts []WebhookAttempt `json:"webhook_attempts"`
	// Meta contains pagination metadata when Moyasar includes it.
	Meta PageMeta `json:"meta"`
}

// Create registers a webhook endpoint.
func (s *WebhookService) Create(ctx context.Context, req CreateWebhookRequest) (*Webhook, error) {
	var webhook Webhook
	if err := s.client.do(ctx, http.MethodPost, "/webhooks", nil, req, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// List lists all registered webhooks.
func (s *WebhookService) List(ctx context.Context) (*WebhookList, error) {
	var list WebhookList
	if err := s.client.do(ctx, http.MethodGet, "/webhooks", nil, nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Get fetches a webhook by its ID.
func (s *WebhookService) Get(ctx context.Context, id string) (*Webhook, error) {
	var webhook Webhook
	if err := s.client.do(ctx, http.MethodGet, "/webhooks/"+url.PathEscape(id), nil, nil, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// Delete deletes a webhook by its ID.
func (s *WebhookService) Delete(ctx context.Context, id string) (*DeleteWebhookResponse, error) {
	var response DeleteWebhookResponse
	if err := s.client.do(ctx, http.MethodDelete, "/webhooks/"+url.PathEscape(id), nil, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// AvailableEvents lists all available webhook event types.
func (s *WebhookService) AvailableEvents(ctx context.Context) (*AvailableWebhookEvents, error) {
	var events AvailableWebhookEvents
	if err := s.client.do(ctx, http.MethodGet, "/webhooks/available_events", nil, nil, &events); err != nil {
		return nil, err
	}
	return &events, nil
}

// ListAttempts lists webhook delivery attempts.
//
// Moyasar retries webhook delivery five more times when the endpoint does not
// return a 2xx response, then drops the message.
func (s *WebhookService) ListAttempts(ctx context.Context, params ListWebhookAttemptsParams) (*WebhookAttemptList, error) {
	var list WebhookAttemptList
	if err := s.client.do(ctx, http.MethodGet, "/webhooks/attempts", params.values(), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// GetAttempt fetches a webhook delivery attempt by its ID.
func (s *WebhookService) GetAttempt(ctx context.Context, id string) (*WebhookAttempt, error) {
	var attempt WebhookAttempt
	if err := s.client.do(ctx, http.MethodGet, "/webhooks/attempts/"+url.PathEscape(id), nil, nil, &attempt); err != nil {
		return nil, err
	}
	return &attempt, nil
}

func (p ListWebhookAttemptsParams) values() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.WebhookID != "" {
		values.Set("webhook_id", p.WebhookID)
	}
	if p.EventID != "" {
		values.Set("event_id", p.EventID)
	}
	if p.EventType != "" {
		values.Set("event_type", string(p.EventType))
	}
	return values
}
