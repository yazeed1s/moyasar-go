package moyasar

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// TransferService provides access to Moyasar aggregation merchant transfer APIs.
//
// These APIs use the transfer base URL, https://apimig.moyasar.com/v1 by default,
// and are only available for Moyasar aggregation merchants.
type TransferService struct {
	client *Client
}

// Transfer is a settlement transfer made for an aggregation merchant account.
type Transfer struct {
	// ID is the transfer ID.
	ID string `json:"id"`
	// RecipientType is Entity, Platform, or Beneficiary.
	RecipientType string `json:"recipient_type"`
	// RecipientID is the transfer recipient ID.
	RecipientID string `json:"recipient_id"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// Amount is the transfer amount.
	Amount int `json:"amount"`
	// Fee is the transfer fee.
	Fee int `json:"fee"`
	// Tax is VAT.
	Tax int `json:"tax"`
	// Reference is the bank reference.
	Reference string `json:"reference"`
	// TransactionCount is the number of transactions included in the transfer.
	TransactionCount int `json:"transaction_count"`
	// CreatedAt is the time the transfer was created.
	CreatedAt string `json:"created_at"`
}

// ListTransfersParams paginates transfer list results.
type ListTransfersParams struct {
	// Page is the requested page number.
	Page int
}

// TransferList is a page of transfers.
type TransferList struct {
	// Transfers contains the returned transfer objects.
	Transfers []Transfer `json:"transfers"`
	// Meta contains pagination metadata.
	Meta PageMeta `json:"meta"`
}

// TransferLine is a line for a given aggregation merchant transfer.
type TransferLine struct {
	// PaymentID is the payment ID this transfer line is related to.
	PaymentID string `json:"payment_id"`
	// Type is the operation type, such as payment.
	Type string `json:"type"`
	// Amount is the transfer line amount.
	Amount int `json:"amount"`
	// Fee is the transfer line fee.
	Fee int `json:"fee"`
	// Tax is VAT.
	Tax int `json:"tax"`
}

// ListTransferLinesParams paginates transfer line results.
type ListTransferLinesParams struct {
	// Page is the requested page number.
	Page int
}

// TransferLineList is a page of transfer lines.
type TransferLineList struct {
	// Lines contains the returned transfer line objects.
	Lines []TransferLine `json:"lines"`
	// Meta contains pagination metadata.
	Meta PageMeta `json:"meta"`
}

// List lists all transfers that have been performed for the account.
func (s *TransferService) List(ctx context.Context, params ListTransfersParams) (*TransferList, error) {
	var list TransferList
	if err := s.client.doTransfer(ctx, http.MethodGet, "/transfers", params.values(), &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Get fetches a single transfer by ID.
func (s *TransferService) Get(ctx context.Context, id string) (*Transfer, error) {
	var transfer Transfer
	if err := s.client.doTransfer(ctx, http.MethodGet, "/transfers/"+url.PathEscape(id), nil, &transfer); err != nil {
		return nil, err
	}
	return &transfer, nil
}

// ListLines lists all lines for a given transfer.
func (s *TransferService) ListLines(ctx context.Context, id string, params ListTransferLinesParams) (*TransferLineList, error) {
	var list TransferLineList
	if err := s.client.doTransfer(ctx, http.MethodGet, "/transfers/"+url.PathEscape(id)+"/lines", params.values(), &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (p ListTransfersParams) values() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	return values
}

func (p ListTransferLinesParams) values() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	return values
}
