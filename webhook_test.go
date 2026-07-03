package moyasar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookCreateSendsJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/webhooks" {
			t.Fatalf("path = %s, want /webhooks", r.URL.Path)
		}
		var body CreateWebhookRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.HTTPMethod != "post" {
			t.Fatalf("http_method = %q, want post", body.HTTPMethod)
		}
		if body.SharedSecret != "123" {
			t.Fatalf("shared_secret = %q, want 123", body.SharedSecret)
		}
		if len(body.Events) != 2 || body.Events[0] != WebhookEventPaymentPaid {
			t.Fatalf("events = %#v, want payment events", body.Events)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"wh_123","http_method":"post","url":"https://example.com/updatepayments","created_at":"2022-12-07T08:24:23.097Z","events":["payment_paid","payment_failed"]}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	webhook, err := client.Webhooks.Create(context.Background(), CreateWebhookRequest{
		HTTPMethod:   "post",
		URL:          "https://example.com/updatepayments",
		SharedSecret: "123",
		Events: []WebhookEventType{
			WebhookEventPaymentPaid,
			WebhookEventPaymentFailed,
		},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if webhook.ID != "wh_123" {
		t.Fatalf("webhook.ID = %q, want wh_123", webhook.ID)
	}
}

func TestWebhookAvailableEvents(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/webhooks/available_events" {
			t.Fatalf("path = %s, want /webhooks/available_events", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"events":["payment_paid","card_auth_authenticated"]}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	events, err := client.Webhooks.AvailableEvents(context.Background())
	if err != nil {
		t.Fatalf("AvailableEvents returned error: %v", err)
	}
	if len(events.Events) != 2 || events.Events[1] != WebhookEventCardAuthAuthenticated {
		t.Fatalf("events = %#v, want card_auth_authenticated", events.Events)
	}
}

func TestWebhookDeleteReturnsMessage(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/webhooks/wh_123" {
			t.Fatalf("path = %s, want /webhooks/wh_123", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"Webhook was deleted successfully"}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	response, err := client.Webhooks.Delete(context.Background(), "wh_123")
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if response.Message != "Webhook was deleted successfully" {
		t.Fatalf("message = %q, want deletion message", response.Message)
	}
}

func TestWebhookAttemptListBuildsQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/webhooks/attempts" {
			t.Fatalf("path = %s, want /webhooks/attempts", r.URL.Path)
		}
		query := r.URL.Query()
		if got := query.Get("page"); got != "2" {
			t.Fatalf("page = %q, want 2", got)
		}
		if got := query.Get("webhook_id"); got != "wh_123" {
			t.Fatalf("webhook_id = %q, want wh_123", got)
		}
		if got := query.Get("event_type"); got != "payment_paid" {
			t.Fatalf("event_type = %q, want payment_paid", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"webhook_attempts":[{"id":"att_123","webhook_id":"wh_123","event_id":"evt_123","event_type":"payment_paid","retry_number":1,"result":"success","message":"Webhook message was delivered successfully","response_code":200,"response_headers":"{h_foo:\"h_bar\"}","response_body":"{foo:\"bar\"}","created_at":"2022-12-24T10:34:00.005Z"}],"meta":{"current_page":2,"next_page":null,"prev_page":1,"total_pages":2,"total_count":1}}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	list, err := client.Webhooks.ListAttempts(context.Background(), ListWebhookAttemptsParams{
		Page:      2,
		WebhookID: "wh_123",
		EventType: WebhookEventPaymentPaid,
	})
	if err != nil {
		t.Fatalf("ListAttempts returned error: %v", err)
	}
	if len(list.Attempts) != 1 || list.Attempts[0].ResponseCode != 200 {
		t.Fatalf("attempts = %#v, want one successful attempt", list.Attempts)
	}
}

func TestWebhookEventPreservesRawData(t *testing.T) {
	t.Parallel()

	var event WebhookEvent
	if err := json.Unmarshal([]byte(`{"id":"evt_123","type":"card_auth_authenticated","created_at":"2026-05-20T10:00:00Z","secret_token":"secret","account_name":"My Store","live":true,"data":{"id":"ca_123","status":"authenticated"}}`), &event); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if event.Type != WebhookEventCardAuthAuthenticated {
		t.Fatalf("type = %q, want card_auth_authenticated", event.Type)
	}
	if string(event.Data) != `{"id":"ca_123","status":"authenticated"}` {
		t.Fatalf("data = %s, want raw card auth object", event.Data)
	}
}
