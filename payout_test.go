package moyasar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPayoutCreateAccountSendsPropertiesAndCredentials(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/payout_accounts" {
			t.Fatalf("path = %s, want /payout_accounts", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got := body["account_type"]; got != "bank" {
			t.Fatalf("account_type = %#v, want bank", got)
		}
		properties := body["properties"].(map[string]any)
		if got := properties["iban"]; got != "SA8430400108057386290038" {
			t.Fatalf("properties.iban = %#v, want IBAN", got)
		}
		credentials := body["credentials"].(map[string]any)
		if got := credentials["private_key"]; got != "secret" {
			t.Fatalf("credentials.private_key = %#v, want secret", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"acct_123","account_type":"bank","currency":"SAR","properties":{"iban":"SA8430400108057386290038"},"created_at":"2026-07-03T00:00:00Z"}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	account, err := client.Payouts.CreateAccount(context.Background(), CreatePayoutAccountRequest{
		AccountType: PayoutAccountTypeBank,
		Properties:  PayoutObject{"iban": "SA8430400108057386290038"},
		Credentials: PayoutObject{
			"private_key": "secret",
		},
	})
	if err != nil {
		t.Fatalf("CreateAccount returned error: %v", err)
	}
	if account.ID != "acct_123" {
		t.Fatalf("account.ID = %q, want acct_123", account.ID)
	}
}

func TestPayoutCreateSendsBankDestinationType(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/payouts" {
			t.Fatalf("path = %s, want /payouts", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		destination := body["destination"].(map[string]any)
		if got := destination["type"]; got != "bank" {
			t.Fatalf("destination.type = %#v, want bank", got)
		}
		if got := destination["iban"]; got != "SA8430400108057386290038" {
			t.Fatalf("destination.iban = %#v, want IBAN", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"po_123","source_id":"acct_123","sequence_number":"6244377266243449","channel":"ips","status":"queued","amount":100,"currency":"SAR","purpose":"bills_or_rent","comment":"my comment","destination":{"type":"bank","iban":"SA8430400108057386290038","name":"Beneficiary","mobile":"0500000000","country":"SA","city":"Riyadh"},"message":"queued","failure_reason":"","created_at":"2026-07-03T00:00:00Z","updated_at":"2026-07-03T00:00:00Z","metadata":{"order_id":"1000"}}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	payout, err := client.Payouts.Create(context.Background(), CreatePayoutRequest{
		SourceID:       "acct_123",
		SequenceNumber: "6244377266243449",
		Amount:         100,
		Purpose:        PayoutPurposeBillsOrRent,
		Destination: BankPayoutDestination{
			IBAN:    "SA8430400108057386290038",
			Name:    "Beneficiary",
			Mobile:  "0500000000",
			Country: "SA",
			City:    "Riyadh",
		},
		Comment:  "my comment",
		Metadata: Metadata{"order_id": "1000"},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if payout.ID != "po_123" {
		t.Fatalf("payout.ID = %q, want po_123", payout.ID)
	}
	if payout.Destination == nil || payout.Destination.Type != "bank" {
		t.Fatalf("destination = %#v, want bank", payout.Destination)
	}
}

func TestPayoutListAccountsBuildsQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Fatalf("page = %q, want 2", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"payout_accounts":[{"id":"acct_123","account_type":"bank","currency":"SAR","properties":{"iban":"SA8430400108057386290038"},"created_at":"2026-07-03T00:00:00Z"}],"meta":{"current_page":2,"next_page":null,"prev_page":1,"total_pages":2,"total_count":1}}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	list, err := client.Payouts.ListAccounts(context.Background(), ListPayoutAccountsParams{Page: 2})
	if err != nil {
		t.Fatalf("ListAccounts returned error: %v", err)
	}
	if len(list.PayoutAccounts) != 1 || list.Meta.CurrentPage != 2 {
		t.Fatalf("list = %#v, want one account on page 2", list)
	}
}

func TestPayoutGetUsesDocumentedSingularPath(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/payout/po_123" {
			t.Fatalf("path = %s, want /payout/po_123", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"po_123","source_id":"acct_123","sequence_number":"6244377266243449","channel":"ips","status":"paid","amount":100,"currency":"SAR","purpose":"bills_or_rent","comment":"my comment","destination":{"type":"wallet","mobile":"0500000000"},"message":"paid","failure_reason":"","created_at":"2026-07-03T00:00:00Z","updated_at":"2026-07-03T00:00:00Z","metadata":{}}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	payout, err := client.Payouts.Get(context.Background(), "po_123")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if payout.Status != PayoutStatusPaid {
		t.Fatalf("status = %q, want paid", payout.Status)
	}
}

func TestPayoutBulkCreate(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/payouts/bulk" {
			t.Fatalf("path = %s, want /payouts/bulk", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		payouts := body["payouts"].([]any)
		item := payouts[0].(map[string]any)
		destination := item["destination"].(map[string]any)
		if got := destination["type"]; got != "wallet" {
			t.Fatalf("destination.type = %#v, want wallet", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"payouts":[{"id":"po_123","source_id":"acct_123","sequence_number":"6244377266243449","channel":"internal","status":"queued","amount":100,"currency":"SAR","purpose":"personal","comment":"bulk","destination":{"type":"wallet","mobile":"0500000000"},"message":"queued","failure_reason":"","created_at":"2026-07-03T00:00:00Z","updated_at":"2026-07-03T00:00:00Z","metadata":{}}]}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	response, err := client.Payouts.BulkCreate(context.Background(), CreateBulkPayoutRequest{
		SourceID: "acct_123",
		Payouts: []BulkPayoutItem{
			{
				Amount:      100,
				Purpose:     PayoutPurposePersonal,
				Destination: WalletPayoutDestination{Mobile: "0500000000"},
				Comment:     "bulk",
			},
		},
	})
	if err != nil {
		t.Fatalf("BulkCreate returned error: %v", err)
	}
	if len(response.Payouts) != 1 || response.Payouts[0].Channel != PayoutChannelInternal {
		t.Fatalf("response = %#v, want one internal payout", response)
	}
}
