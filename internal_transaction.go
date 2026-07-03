package moyasar

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// InternalTransactionService provides access to Moyasar internal transaction APIs.
type InternalTransactionService struct {
	client *Client
}

// InternalTransaction transfers an amount from the current wallet of the caller
// to the current wallet of the recipient.
type InternalTransaction struct {
	// ID is the internal transaction ID.
	ID string `json:"id"`
	// RecipientType is Entity, Platform, or Beneficiary.
	RecipientType string `json:"recipient_type"`
	// RecipientID is the recipient Entity, Platform, or Beneficiary ID.
	RecipientID string `json:"recipient_id"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// Amount is the internal transaction amount in the smallest currency unit.
	Amount int `json:"amount"`
	// TransferID is set when the internal transaction has been settled.
	TransferID string `json:"transfer_id"`
	// Description describes the purpose of this internal transaction.
	Description string `json:"description"`
	// CreatedAt is the time the internal transaction was created.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the time the internal transaction status was updated.
	UpdatedAt string `json:"updated_at"`
	// SettledAt is the time the internal transaction was settled.
	SettledAt string `json:"settled_at"`
	// Metadata is merchant-defined key/value data returned in responses and webhooks.
	Metadata Metadata `json:"metadata"`
}

// CreateInternalTransactionRequest creates an instant internal transaction.
//
// To revert a successful transaction, the recipient has to transfer the amount back.
type CreateInternalTransactionRequest struct {
	// RecipientID is the Entity, Platform, or Beneficiary receiving the amount.
	RecipientID string `json:"recipient_id"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// Amount is the internal transaction amount in the smallest currency unit.
	Amount int `json:"amount"`
	// Description describes the purpose of this internal transaction.
	Description string `json:"description,omitempty"`
	// Metadata is merchant-defined key/value data returned in responses and webhooks.
	Metadata Metadata `json:"metadata,omitempty"`
}

// ListInternalTransactionsParams filters and paginates internal transaction list results.
type ListInternalTransactionsParams struct {
	// Page is the requested page number.
	Page int
	// ID filters by internal transaction ID.
	ID string
	// Currency filters transactions by currency.
	Currency string
	// CreatedAtGT filters transactions created equal to or greater than this timestamp.
	CreatedAtGT string
	// CreatedAtLT filters transactions created equal to or less than this timestamp.
	CreatedAtLT string
	// UpdatedAtGT filters transactions updated equal to or greater than this timestamp.
	UpdatedAtGT string
	// UpdatedAtLT filters transactions updated equal to or less than this timestamp.
	UpdatedAtLT string
	// SettledAtGT filters transactions settled equal to or greater than this timestamp.
	SettledAtGT string
	// SettledAtLT filters transactions settled equal to or less than this timestamp.
	SettledAtLT string
}

// InternalTransactionList is a page of internal transactions.
//
// Moyasar documents this API as infinite scroll; meta.total_count is always null.
type InternalTransactionList struct {
	// InternalTransactions contains the returned internal transaction objects.
	InternalTransactions []InternalTransaction `json:"internal_transactions"`
	// Meta contains pagination metadata.
	Meta PageMeta `json:"meta"`
}

// Create transfers an amount from the current wallet of the caller to a recipient.
func (s *InternalTransactionService) Create(ctx context.Context, req CreateInternalTransactionRequest) (*InternalTransaction, error) {
	var transaction InternalTransaction
	if err := s.client.do(ctx, http.MethodPost, "/internal_transactions", nil, req, &transaction); err != nil {
		return nil, err
	}
	return &transaction, nil
}

// List lists internal transactions performed on the authorized Entity or Platform.
func (s *InternalTransactionService) List(ctx context.Context, params ListInternalTransactionsParams) (*InternalTransactionList, error) {
	var list InternalTransactionList
	if err := s.client.do(ctx, http.MethodGet, "/internal_transactions", params.values(), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (p ListInternalTransactionsParams) values() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.ID != "" {
		values.Set("id", p.ID)
	}
	if p.Currency != "" {
		values.Set("currency", p.Currency)
	}
	if p.CreatedAtGT != "" {
		values.Set("created_at[gt]", p.CreatedAtGT)
	}
	if p.CreatedAtLT != "" {
		values.Set("created_at[lt]", p.CreatedAtLT)
	}
	if p.UpdatedAtGT != "" {
		values.Set("updated_at[gt]", p.UpdatedAtGT)
	}
	if p.UpdatedAtLT != "" {
		values.Set("updated_at[lt]", p.UpdatedAtLT)
	}
	if p.SettledAtGT != "" {
		values.Set("settled_at[gt]", p.SettledAtGT)
	}
	if p.SettledAtLT != "" {
		values.Set("settled_at[lt]", p.SettledAtLT)
	}
	return values
}
