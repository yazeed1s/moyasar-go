package moyasar

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// PaymentService provides access to Moyasar payment APIs.
type PaymentService struct {
	client *Client
}

// PaymentStatus is the lifecycle status of a Moyasar payment.
type PaymentStatus string

const (
	PaymentStatusInitiated  PaymentStatus = "initiated"
	PaymentStatusPaid       PaymentStatus = "paid"
	PaymentStatusAuthorized PaymentStatus = "authorized"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusCaptured   PaymentStatus = "captured"
	PaymentStatusVoided     PaymentStatus = "voided"
	PaymentStatusVerified   PaymentStatus = "verified"
)

// Payment is a Moyasar payment object.
type Payment struct {
	// ID is the unique payment identifier. If given_id is supplied on create,
	// Moyasar uses that value as the payment ID.
	ID string `json:"id"`
	// Status indicates the payment lifecycle state.
	Status PaymentStatus `json:"status"`
	// Amount is the payment amount in the smallest currency unit.
	Amount int `json:"amount"`
	// Fee is the estimated payment fee, including VAT.
	Fee int `json:"fee"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// Refunded is the refunded amount in the smallest currency unit.
	Refunded int `json:"refunded"`
	// RefundedAt is the time the payment was refunded, when present.
	RefundedAt *string `json:"refunded_at"`
	// Captured is the captured amount in the smallest currency unit.
	Captured int `json:"captured"`
	// CapturedAt is the time the payment was captured, when present.
	CapturedAt *string `json:"captured_at"`
	// VoidedAt is the time the payment was voided, when present.
	VoidedAt *string `json:"voided_at"`
	// Description is the merchant-facing payment description.
	Description string `json:"description"`
	// AmountFormat is the formatted payment amount with currency.
	AmountFormat string `json:"amount_format"`
	// FeeFormat is the formatted payment fee with currency.
	FeeFormat string `json:"fee_format"`
	// RefundedFormat is the formatted refunded amount with currency.
	RefundedFormat string `json:"refunded_format"`
	// CapturedFormat is the formatted captured amount with currency.
	CapturedFormat string `json:"captured_format"`
	// InvoiceID is the invoice ID this payment is used to pay, when present.
	InvoiceID string `json:"invoice_id"`
	// IP is the payer IPv4 address collected by Moyasar, when present.
	IP string `json:"ip"`
	// CallbackURL is the URL used to return the payer after card payment.
	CallbackURL string `json:"callback_url"`
	// CreatedAt is the time the payment was created.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the time the payment was last updated.
	UpdatedAt string `json:"updated_at"`
	// Metadata is merchant-defined key/value data returned with the payment.
	Metadata Metadata `json:"metadata"`
	// Source is the payment source object returned by Moyasar.
	Source *RawPaymentSource `json:"source"`
	// Splits are settlement split rules associated with the payment.
	Splits []PaymentSplit `json:"splits,omitempty"`
}

// PaymentSplit describes a split of payment funds to a recipient.
type PaymentSplit struct {
	// Amount is the split amount in the smallest currency unit.
	Amount int `json:"amount"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency,omitempty"`
	// RecipientType is Entity, Platform, or Beneficiary.
	RecipientType string `json:"recipient_type"`
	// RecipientID is the recipient Entity, Platform, or Beneficiary ID.
	RecipientID string `json:"recipient_id"`
	// FeeSource indicates whether the recipient is the payment fee source.
	FeeSource bool `json:"fee_source"`
	// Reference is a merchant-provided reference.
	Reference string `json:"reference,omitempty"`
	// Description is a merchant-provided split description.
	Description string `json:"description,omitempty"`
	// Refundable indicates whether the split should be reversed on refund.
	Refundable bool `json:"refundable"`
	// Metadata is merchant-defined key/value data for the split.
	Metadata Metadata `json:"metadata,omitempty"`
}

// CreatePaymentRequest starts a new card, token, Apple Pay, Samsung Pay, or STC
// Pay payment.
type CreatePaymentRequest struct {
	// GivenID is a merchant-generated UUID used for idempotent payment creation.
	// When provided, it becomes the created payment ID.
	GivenID string `json:"given_id,omitempty"`
	// Amount is the payment amount in the smallest currency unit.
	Amount int `json:"amount"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// Description is shown to the merchant only and not to the payer.
	Description string `json:"description,omitempty"`
	// CallbackURL returns the payer to the merchant website after card payment.
	// Moyasar requires it for creditcard and token sources.
	CallbackURL string `json:"callback_url,omitempty"`
	// Source defines the payment method to charge.
	Source PaymentSource `json:"source"`
	// Metadata is merchant-defined key/value data returned in responses and webhooks.
	Metadata Metadata `json:"metadata,omitempty"`
	// ApplyCoupon controls coupon application when needed.
	ApplyCoupon *bool `json:"apply_coupon,omitempty"`
	// Splits are settlement split rules for the payment.
	Splits []PaymentSplit `json:"splits,omitempty"`
}

// ListPaymentsParams filters and paginates payment list results.
type ListPaymentsParams struct {
	// Page is the requested page number.
	Page int
	// ID filters by payment ID.
	ID string
	// Status filters by payment status.
	Status PaymentStatus
	// Metadata filters by metadata key/value pairs.
	Metadata Metadata
}

// PaymentList is a page of payments.
type PaymentList struct {
	// Payments contains the returned payment objects.
	Payments []Payment `json:"payments"`
	// Meta contains pagination metadata.
	Meta PageMeta `json:"meta"`
}

// UpdatePaymentRequest updates mutable payment fields.
type UpdatePaymentRequest struct {
	// Description is the merchant-facing payment description.
	Description string `json:"description,omitempty"`
	// Metadata replaces or updates merchant-defined key/value data.
	Metadata Metadata `json:"metadata,omitempty"`
}

// RefundPaymentRequest refunds a captured Moyasar payment.
type RefundPaymentRequest struct {
	// Amount is the optional amount to refund. If omitted, Moyasar refunds the
	// full available amount.
	Amount int `json:"amount,omitempty"`
}

// CapturePaymentRequest captures an authorized Moyasar payment.
type CapturePaymentRequest struct {
	// Amount is the optional amount to capture.
	Amount int `json:"amount,omitempty"`
}

// Create starts a new Moyasar payment.
//
// When the returned payment has status initiated, the payer must usually
// complete the challenge at source.transaction_url.
func (s *PaymentService) Create(ctx context.Context, req CreatePaymentRequest) (*Payment, error) {
	var payment Payment
	if err := s.client.do(ctx, http.MethodPost, "/payments", nil, req, &payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

// Get fetches a payment by ID.
func (s *PaymentService) Get(ctx context.Context, id string) (*Payment, error) {
	var payment Payment
	if err := s.client.do(ctx, http.MethodGet, "/payments/"+url.PathEscape(id), nil, nil, &payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

// List returns a paginated reverse-chronological stream of payments.
func (s *PaymentService) List(ctx context.Context, params ListPaymentsParams) (*PaymentList, error) {
	var list PaymentList
	if err := s.client.do(ctx, http.MethodGet, "/payments", params.values(), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Update updates the description or metadata of an existing payment.
func (s *PaymentService) Update(ctx context.Context, id string, req UpdatePaymentRequest) (*Payment, error) {
	var payment Payment
	if err := s.client.do(ctx, http.MethodPut, "/payments/"+url.PathEscape(id), nil, req, &payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

// Refund refunds a paid or captured payment, fully or partially.
func (s *PaymentService) Refund(ctx context.Context, id string, req RefundPaymentRequest) (*Payment, error) {
	var payment Payment
	if err := s.client.do(ctx, http.MethodPost, "/payments/"+url.PathEscape(id)+"/refund", nil, req, &payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

// Capture captures a manually authorized payment.
func (s *PaymentService) Capture(ctx context.Context, id string, req CapturePaymentRequest) (*Payment, error) {
	var payment Payment
	if err := s.client.do(ctx, http.MethodPost, "/payments/"+url.PathEscape(id)+"/capture", nil, req, &payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

// Void cancels a paid, authorized, or captured payment when it has not settled.
func (s *PaymentService) Void(ctx context.Context, id string) (*Payment, error) {
	var payment Payment
	if err := s.client.do(ctx, http.MethodPost, "/payments/"+url.PathEscape(id)+"/void", nil, nil, &payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

func (p ListPaymentsParams) values() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.ID != "" {
		values.Set("id", p.ID)
	}
	if p.Status != "" {
		values.Set("status", string(p.Status))
	}
	for key, value := range p.Metadata {
		values.Set("metadata["+key+"]", value)
	}
	return values
}
