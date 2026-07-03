package moyasar

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPaymentCreateSendsAuthAndSourceType(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/payments" {
			t.Fatalf("path = %s, want /payments", r.URL.Path)
		}
		wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("sk_test_123:"))
		if got := r.Header.Get("Authorization"); got != wantAuth {
			t.Fatalf("Authorization = %q, want %q", got, wantAuth)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		source, ok := body["source"].(map[string]any)
		if !ok {
			t.Fatalf("source missing or not object: %#v", body["source"])
		}
		if got := source["type"]; got != "creditcard" {
			t.Fatalf("source.type = %#v, want creditcard", got)
		}
		if got := source["name"]; got != "John Doe" {
			t.Fatalf("source.name = %#v, want John Doe", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id":"pay_123",
			"status":"initiated",
			"amount":100,
			"currency":"SAR",
			"source":{"type":"creditcard","transaction_url":"https://example.com/3ds"}
		}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	payment, err := client.Payments.Create(context.Background(), CreatePaymentRequest{
		GivenID:     "a1168bd1-47a4-4b97-8a50-dd5caaccacf2",
		Amount:      100,
		Currency:    "SAR",
		Description: "Order #123",
		CallbackURL: "https://example.com/callback",
		Source: CreditCardSource{
			Name:   "John Doe",
			Number: "4111111111111111",
			Month:  9,
			Year:   2030,
			CVC:    "123",
		},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if payment.ID != "pay_123" {
		t.Fatalf("payment.ID = %q, want pay_123", payment.ID)
	}
	if payment.Source == nil || payment.Source.Type != "creditcard" {
		t.Fatalf("payment.Source = %#v, want creditcard source", payment.Source)
	}
	if !strings.Contains(string(payment.Source.Raw), "transaction_url") {
		t.Fatalf("payment.Source.Raw = %s, want transaction_url", payment.Source.Raw)
	}
}

func TestPaymentListBuildsQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Fatalf("page = %q, want 2", got)
		}
		if got := r.URL.Query().Get("status"); got != "paid" {
			t.Fatalf("status = %q, want paid", got)
		}
		if got := r.URL.Query().Get("metadata[order_id]"); got != "1000" {
			t.Fatalf("metadata[order_id] = %q, want 1000", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"payments":[],"meta":{"current_page":2,"next_page":null,"prev_page":1,"total_pages":3,"total_count":80}}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	list, err := client.Payments.List(context.Background(), ListPaymentsParams{
		Page:   2,
		Status: PaymentStatusPaid,
		Metadata: Metadata{
			"order_id": "1000",
		},
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if list.Meta.CurrentPage != 2 || list.Meta.TotalPages != 3 {
		t.Fatalf("meta = %#v, want current page 2 and total pages 3", list.Meta)
	}
}

func TestAPIErrorDecoding(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"type":"invalid_request_error","message":"Validation Failed","errors":{"amount":["must be an integer"]}}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	_, err := client.Payments.Get(context.Background(), "missing")
	if err == nil {
		t.Fatal("Get returned nil error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode = %d, want 400", apiErr.StatusCode)
	}
	if apiErr.Type != "invalid_request_error" {
		t.Fatalf("Type = %q, want invalid_request_error", apiErr.Type)
	}
	if len(apiErr.Body) == 0 {
		t.Fatal("Body is empty")
	}
}
