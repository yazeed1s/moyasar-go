package moyasar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInvoiceListBuildsQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("page"); got != "3" {
			t.Fatalf("page = %q, want 3", got)
		}
		if got := query.Get("status"); got != "paid" {
			t.Fatalf("status = %q, want paid", got)
		}
		if got := query.Get("created[gt]"); got != "2026-01-01T00:00:00Z" {
			t.Fatalf("created[gt] = %q, want timestamp", got)
		}
		if got := query.Get("metadata[order_id]"); got != "1000" {
			t.Fatalf("metadata[order_id] = %q, want 1000", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"invoices":[],"meta":{"current_page":3,"next_page":null,"prev_page":2,"total_pages":3,"total_count":90}}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	list, err := client.Invoices.List(context.Background(), ListInvoicesParams{
		Page:      3,
		Status:    InvoiceStatusPaid,
		CreatedGT: "2026-01-01T00:00:00Z",
		Metadata: Metadata{
			"order_id": "1000",
		},
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if list.Meta.CurrentPage != 3 {
		t.Fatalf("current page = %d, want 3", list.Meta.CurrentPage)
	}
}

func TestTokenCreateUsesFormEncoding(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("Content-Type = %q, want application/x-www-form-urlencoded", got)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm: %v", err)
		}
		if got := r.Form.Get("name"); got != "Mohammed Ali" {
			t.Fatalf("name = %q, want Mohammed Ali", got)
		}
		if got := r.Form.Get("callback_url"); got != "https://mystore.com/thanks" {
			t.Fatalf("callback_url = %q, want callback URL", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"token_123","status":"initiated","brand":"visa","funding":"credit","country":"US","month":"09","year":"2027","name":"Mohammed Ali","last_four":"1111"}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("pk_test_123", WithBaseURL(server.URL))
	token, err := client.Tokens.Create(context.Background(), CreateTokenRequest{
		Name:        "Mohammed Ali",
		Number:      "4111111111111111",
		Month:       "09",
		Year:        "27",
		CVC:         "911",
		CallbackURL: "https://mystore.com/thanks",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if token.ID != "token_123" {
		t.Fatalf("token.ID = %q, want token_123", token.ID)
	}
}

func TestRetrieveIssuerSendsSourceType(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		source := body["source"]
		if got := source["type"]; got != "creditcard" {
			t.Fatalf("source.type = %#v, want creditcard", got)
		}
		if got := source["number"]; got != "4111111111111111" {
			t.Fatalf("source.number = %#v, want card number", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"issuer_name":"Moyasar Sandbox Bank","issuer_country":"SA","issuer_card_type":"debit","issuer_card_category":"SIGNATURE","company":"mada","first_digits":"41111111","last_digits":"1111"}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	issuer, err := client.Sources.RetrieveIssuer(context.Background(), RetrieveIssuerRequest{
		Source: IssuerCreditCardSource{Number: "4111111111111111"},
	})
	if err != nil {
		t.Fatalf("RetrieveIssuer returned error: %v", err)
	}
	if issuer.IssuerName != "Moyasar Sandbox Bank" {
		t.Fatalf("issuer name = %q, want Moyasar Sandbox Bank", issuer.IssuerName)
	}
}

func TestCardAuthCreateSendsCreditCardSourceType(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"ca_123","status":"available","amount":100,"currency":"SAR","callback_url":"https://merchant.example/3ds/return","transaction_url":"https://api.moyasar.com/v1/card_auth/ca_123/prepare","card":{"company":"mada","last_digits":"1111"},"result":null}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	auth, err := client.CardAuths.Create(context.Background(), CreateCardAuthRequest{
		Amount:      100,
		Currency:    "SAR",
		CallbackURL: "https://merchant.example/3ds/return",
		Source: CardAuthSource{
			Name:   "John Doe",
			Number: "4111111111111111",
			Month:  "09",
			Year:   "2030",
			CVC:    "123",
		},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if auth.Status != CardAuthStatusAvailable {
		t.Fatalf("status = %q, want available", auth.Status)
	}
}
