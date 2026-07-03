package moyasar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// SettlementService provides access to Moyasar settlement APIs.
type SettlementService struct {
	client *Client
}

// Settlement is a transfer to a Moyasar merchant bank account.
type Settlement struct {
	// ID is the settlement ID.
	ID string `json:"id"`
	// RecipientType is Entity, Platform, or Beneficiary.
	RecipientType string `json:"recipient_type"`
	// RecipientID is the settlement recipient ID.
	RecipientID string `json:"recipient_id"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// SourceCurrency is the original currency of transactions in this settlement.
	SourceCurrency string `json:"source_currency"`
	// InvoicingCurrency is the currency used for invoicing.
	InvoicingCurrency string `json:"invoicing_currency"`
	// Amount is the full settlement amount before subtracting the bank transfer fee.
	Amount int `json:"amount"`
	// Fee is the settlement bank transfer fee, including VAT.
	Fee int `json:"fee"`
	// Tax is VAT.
	Tax int `json:"tax"`
	// InvoicingFee is the settlement fee in invoicing currency.
	InvoicingFee int `json:"invoicing_fee"`
	// InvoicingTax is VAT in invoicing currency.
	InvoicingTax int `json:"invoicing_tax"`
	// InvoicingExchangeRate is the exchange rate used to convert invoicing_fee to fee.
	InvoicingExchangeRate float64 `json:"invoicing_ex_rate"`
	// Reference is the settlement bank transfer sequence number.
	Reference *string `json:"reference"`
	// SettlementCount is the number of transactions included in the settlement.
	SettlementCount int `json:"settlement_count"`
	// InvoiceURL is the settlement invoice PDF file URL.
	InvoiceURL string `json:"invoice_url"`
	// CSVListURL is the settlement transaction CSV list URL.
	CSVListURL string `json:"csv_list_url"`
	// PDFListURL is the settlement transaction PDF list URL. This is not always available.
	PDFListURL string `json:"pdf_list_url"`
	// CreatedAt is the time the settlement was created.
	CreatedAt string `json:"created_at"`
}

// ListSettlementsParams filters and paginates settlement list results.
type ListSettlementsParams struct {
	// Page is the requested page number.
	Page int
	// ID filters by settlement ID.
	ID string
	// CreatedGT filters settlements created equal to or greater than this timestamp.
	CreatedGT string
	// CreatedLT filters settlements created equal to or less than this timestamp.
	CreatedLT string
}

// SettlementList is a page of settlements.
type SettlementList struct {
	// Settlements contains the returned settlement objects.
	Settlements []Settlement `json:"settlements"`
	// Meta contains pagination metadata.
	Meta PageMeta `json:"meta"`
}

// SettlementLineType is the type of operation settled.
type SettlementLineType string

const (
	SettlementLineTypePayment           SettlementLineType = "payment"
	SettlementLineTypeRefund            SettlementLineType = "refund"
	SettlementLineTypeVoid              SettlementLineType = "void"
	SettlementLineTypeFee               SettlementLineType = "fee"
	SettlementLineTypePlatformDuties    SettlementLineType = "platform_duties"
	SettlementLineTypeOtherDuties       SettlementLineType = "other_duties"
	SettlementLineTypeChargeback        SettlementLineType = "chargeback"
	SettlementLineTypeChargebackPenalty SettlementLineType = "chargeback_penalty"
	SettlementLineTypeInstallment       SettlementLineType = "installment"
)

// SettlementLine is an individual transaction line that makes up a settlement.
type SettlementLine struct {
	// PaymentID is the payment ID this line is related to.
	PaymentID string `json:"payment_id"`
	// Type is the operation settled, such as payment, refund, void, fee, or chargeback.
	Type SettlementLineType `json:"type"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// SourceCurrency is the currency used to settle this transaction.
	SourceCurrency string `json:"source_currency"`
	// InvoicingCurrency is the currency used for invoicing.
	InvoicingCurrency string `json:"invoicing_currency"`
	// PaymentAmount is the total payment amount.
	PaymentAmount int `json:"payment_amount"`
	// Amount is the net amount settled to the merchant. It can be negative.
	Amount int `json:"amount"`
	// SettlementAmount is the net amount settled in the settlement currency.
	SettlementAmount int `json:"settlement_amount"`
	// Fee is the transaction fee including VAT.
	Fee int `json:"fee"`
	// Tax is VAT.
	Tax int `json:"tax"`
	// InvoicingFee is the transaction fee in invoicing_currency.
	InvoicingFee int `json:"invoicing_fee"`
	// InvoicingTax is VAT in invoicing_currency.
	InvoicingTax int `json:"invoicing_tax"`
	// InvoicingExchangeRate converts fee to invoicing_fee.
	InvoicingExchangeRate float64 `json:"i_ex_rate"`
	// SettlementExchangeRate converts amount to settlement_amount.
	SettlementExchangeRate float64 `json:"r_ex_rate"`
	// ReferenceNumber is the retrieval reference number generated by the acquirer.
	ReferenceNumber string `json:"reference_number"`
	// AuthorizationCode is a six-digit issuer authorization code.
	AuthorizationCode string `json:"authorization_code"`
	// IP is the payer IPv4 address collected by Moyasar.
	IP string `json:"ip"`
	// TransactedAt is the time when the transaction being settled occurred.
	TransactedAt string `json:"transacted_at"`
	// Splits is ignored by Moyasar for now and always returns null.
	Splits json.RawMessage `json:"splits"`
	// CustomSplits lists custom splits applied during payment authorization.
	CustomSplits []PaymentSplit `json:"custom_splits"`
	// IsCustomSplit indicates if this is a custom split created by the API user.
	IsCustomSplit bool `json:"is_custom_split"`
	// SplitReference is the reference added during payment creation.
	SplitReference string `json:"split_reference"`
	// SplitDescription is the human-readable split description added during payment creation.
	SplitDescription string `json:"split_description"`
	// Source is the payment source response object for this settlement line.
	Source *RawPaymentSource `json:"source"`
	// Metadata is merchant-defined key/value data returned in responses and webhooks.
	Metadata Metadata `json:"metadata"`
}

// ListSettlementLinesParams paginates settlement line results.
type ListSettlementLinesParams struct {
	// Page is the requested page number.
	Page int
}

// SettlementLineList is a page of settlement lines.
type SettlementLineList struct {
	// Lines contains the returned settlement line objects.
	Lines []SettlementLine `json:"lines"`
	// Meta contains pagination metadata.
	Meta PageMeta `json:"meta"`
}

// List returns settlements transferred to the merchant bank account.
func (s *SettlementService) List(ctx context.Context, params ListSettlementsParams) (*SettlementList, error) {
	var list SettlementList
	if err := s.client.do(ctx, http.MethodGet, "/settlements", params.values(), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Get retrieves a single settlement by ID.
func (s *SettlementService) Get(ctx context.Context, id string) (*Settlement, error) {
	var settlement Settlement
	if err := s.client.do(ctx, http.MethodGet, "/settlements/"+url.PathEscape(id), nil, nil, &settlement); err != nil {
		return nil, err
	}
	return &settlement, nil
}

// ListLines lists individual transaction lines that make up a settlement.
func (s *SettlementService) ListLines(ctx context.Context, id string, params ListSettlementLinesParams) (*SettlementLineList, error) {
	var list SettlementLineList
	if err := s.client.do(ctx, http.MethodGet, "/settlements/"+url.PathEscape(id)+"/lines", params.values(), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (p ListSettlementsParams) values() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.ID != "" {
		values.Set("id", p.ID)
	}
	if p.CreatedGT != "" {
		values.Set("created[gt]", p.CreatedGT)
	}
	if p.CreatedLT != "" {
		values.Set("created[lt]", p.CreatedLT)
	}
	return values
}

func (p ListSettlementLinesParams) values() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	return values
}
