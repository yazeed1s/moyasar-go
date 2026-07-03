package moyasar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSettlementListBuildsFilters(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("page"); got != "2" {
			t.Fatalf("page = %q, want 2", got)
		}
		if got := query.Get("created[gt]"); got != "2026-01-01T00:00:00Z" {
			t.Fatalf("created[gt] = %q, want timestamp", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"settlements":[{"id":"set_123","recipient_type":"Entity","recipient_id":"ent_123","currency":"SAR","source_currency":"SAR","invoicing_currency":"SAR","amount":24209,"fee":46,"tax":6,"invoicing_fee":46,"invoicing_tax":6,"invoicing_ex_rate":1,"reference":null,"settlement_count":1,"invoice_url":"https://example.com/invoice.pdf","csv_list_url":"https://example.com/list.csv","pdf_list_url":"https://example.com/list.pdf","created_at":"2026-07-03T00:00:00Z"}],"meta":{"current_page":2,"next_page":null,"prev_page":1,"total_pages":2,"total_count":1}}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	list, err := client.Settlements.List(context.Background(), ListSettlementsParams{
		Page:      2,
		CreatedGT: "2026-01-01T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(list.Settlements) != 1 || list.Settlements[0].Amount != 24209 {
		t.Fatalf("list = %#v, want settlement amount 24209", list)
	}
}

func TestSettlementListLines(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/settlements/set_123/lines" {
			t.Fatalf("path = %s, want /settlements/set_123/lines", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "3" {
			t.Fatalf("page = %q, want 3", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"lines":[{"payment_id":"pay_123","type":"payment","currency":"SAR","amount":10000,"fee":0,"tax":0,"source":{"type":"creditcard","company":"visa"}}],"meta":{"current_page":3,"next_page":null,"prev_page":2,"total_pages":3,"total_count":1}}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	list, err := client.Settlements.ListLines(context.Background(), "set_123", ListSettlementLinesParams{Page: 3})
	if err != nil {
		t.Fatalf("ListLines returned error: %v", err)
	}
	if len(list.Lines) != 1 || list.Lines[0].Source.Type != "creditcard" {
		t.Fatalf("lines = %#v, want creditcard source", list.Lines)
	}
}

func TestInternalTransactionCreateAndListFilters(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method + " " + r.URL.Path {
		case "POST /internal_transactions":
			var body CreateInternalTransactionRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body.RecipientID != "ent_123" || body.Amount != 24209 {
				t.Fatalf("body = %#v, want recipient and amount", body)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"it_123","recipient_type":"Entity","recipient_id":"ent_123","currency":"SAR","amount":24209,"transfer_id":"tr_123","description":"string","created_at":"2026-02-15T12:59:11Z","updated_at":"2026-02-15T12:59:11Z","settled_at":"2026-02-16T10:13:12Z","metadata":{"cart_id":"cart_123"}}`))
		case "GET /internal_transactions":
			query := r.URL.Query()
			if got := query.Get("created_at[gt]"); got != "2026-02-01T00:00:00Z" {
				t.Fatalf("created_at[gt] = %q, want timestamp", got)
			}
			if got := query.Get("settled_at[lt]"); got != "2026-03-01T00:00:00Z" {
				t.Fatalf("settled_at[lt] = %q, want timestamp", got)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"internal_transactions":[],"meta":{"current_page":1,"next_page":null,"prev_page":null,"total_pages":1,"total_count":null}}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithBaseURL(server.URL))
	transaction, err := client.InternalTransactions.Create(context.Background(), CreateInternalTransactionRequest{
		RecipientID: "ent_123",
		Currency:    "SAR",
		Amount:      24209,
		Description: "string",
		Metadata:    Metadata{"cart_id": "cart_123"},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if transaction.ID != "it_123" {
		t.Fatalf("transaction.ID = %q, want it_123", transaction.ID)
	}

	_, err = client.InternalTransactions.List(context.Background(), ListInternalTransactionsParams{
		CreatedAtGT: "2026-02-01T00:00:00Z",
		SettledAtLT: "2026-03-01T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
}

func TestTransfersUseTransferBaseURL(t *testing.T) {
	t.Parallel()

	mainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("main base URL should not receive transfer request: %s", r.URL.Path)
	}))
	t.Cleanup(mainServer.Close)

	transferServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/transfers" {
			t.Fatalf("path = %s, want /transfers", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Fatalf("page = %q, want 2", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"transfers":[{"id":"tr_123","recipient_type":"Entity","recipient_id":"ent_123","currency":"SAR","amount":120000,"fee":0,"tax":0,"reference":"bank_ref_789","transaction_count":0,"created_at":"2023-02-11T08:06:54.000Z"}],"meta":{"current_page":2,"next_page":null,"prev_page":1,"total_pages":2,"total_count":1}}`))
	}))
	t.Cleanup(transferServer.Close)

	client := NewClient("sk_test_123", WithBaseURL(mainServer.URL), WithTransferURL(transferServer.URL))
	list, err := client.Transfers.List(context.Background(), ListTransfersParams{Page: 2})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(list.Transfers) != 1 || list.Transfers[0].Reference != "bank_ref_789" {
		t.Fatalf("transfers = %#v, want bank_ref_789", list.Transfers)
	}
}

func TestTransferGetAndListLines(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/transfers/tr_123":
			_, _ = w.Write([]byte(`{"id":"tr_123","recipient_type":"Entity","recipient_id":"ent_123","currency":"SAR","amount":120000,"fee":0,"tax":0,"reference":"bank_ref_789","transaction_count":0,"created_at":"2023-02-11T08:06:54.000Z"}`))
		case "/transfers/tr_123/lines":
			if got := r.URL.Query().Get("page"); got != "2" {
				t.Fatalf("page = %q, want 2", got)
			}
			_, _ = w.Write([]byte(`{"lines":[{"payment_id":"pay_123","type":"payment","amount":10000,"fee":0,"tax":0}],"meta":{"current_page":2,"next_page":3,"prev_page":1,"total_pages":3,"total_count":100}}`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := NewClient("sk_test_123", WithTransferURL(server.URL))
	transfer, err := client.Transfers.Get(context.Background(), "tr_123")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if transfer.ID != "tr_123" {
		t.Fatalf("transfer.ID = %q, want tr_123", transfer.ID)
	}

	lines, err := client.Transfers.ListLines(context.Background(), "tr_123", ListTransferLinesParams{Page: 2})
	if err != nil {
		t.Fatalf("ListLines returned error: %v", err)
	}
	if len(lines.Lines) != 1 || lines.Lines[0].Amount != 10000 {
		t.Fatalf("lines = %#v, want amount 10000", lines.Lines)
	}
}
