package moyasar

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// InvoiceService provides access to Moyasar invoice APIs.
type InvoiceService struct {
	client *Client
}

// InvoiceStatus is the lifecycle status of a Moyasar invoice.
type InvoiceStatus string

const (
	InvoiceStatusInitiated InvoiceStatus = "initiated"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusFailed    InvoiceStatus = "failed"
	InvoiceStatusRefunded  InvoiceStatus = "refunded"
	InvoiceStatusCanceled  InvoiceStatus = "canceled"
	InvoiceStatusOnHold    InvoiceStatus = "on_hold"
	InvoiceStatusExpired   InvoiceStatus = "expired"
	InvoiceStatusVoided    InvoiceStatus = "voided"
)

// Invoice is a Moyasar invoice used to bill a customer through a hosted
// payment page.
type Invoice struct {
	// ID is the unique invoice identifier.
	ID string `json:"id"`
	// Status indicates the invoice lifecycle state.
	Status InvoiceStatus `json:"status"`
	// Amount is the invoice amount in the smallest currency unit.
	Amount int `json:"amount"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// Description is displayed on the invoice alongside the amount.
	Description string `json:"description"`
	// LogoURL is the entity logo configured through Moyasar Dashboard.
	LogoURL string `json:"logo_url"`
	// AmountFormat is the formatted invoice amount with currency.
	AmountFormat string `json:"amount_format"`
	// URL is the checkout page URL that the merchant presents to the payer.
	URL string `json:"url"`
	// CallbackURL receives a POST request with the invoice object when paid.
	CallbackURL string `json:"callback_url"`
	// ExpiredAt is the time after which the invoice can no longer be paid.
	ExpiredAt *string `json:"expired_at"`
	// CreatedAt is the time the invoice was created.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the time the invoice was last updated.
	UpdatedAt string `json:"updated_at"`
	// BackURL is used when the payer clicks the back button.
	BackURL string `json:"back_url"`
	// SuccessURL is where Moyasar redirects the payer when the invoice is paid.
	SuccessURL string `json:"success_url"`
	// Payments contains payment attempts made against this invoice.
	Payments []Payment `json:"payments"`
	// Metadata is merchant-defined key/value data returned in responses and webhooks.
	Metadata Metadata `json:"metadata"`
}

// CreateInvoiceRequest creates a Moyasar invoice to bill a customer and collect
// payment through a hosted payment page.
type CreateInvoiceRequest struct {
	// Amount is the invoice amount in the smallest currency unit.
	Amount int `json:"amount"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// Description is displayed on the invoice alongside the amount.
	Description string `json:"description"`
	// CallbackURL receives a POST request with the invoice object when paid.
	// Unlike Payment, this is not used to redirect the payer.
	CallbackURL string `json:"callback_url,omitempty"`
	// SuccessURL is where Moyasar redirects the payer when the invoice is paid.
	SuccessURL string `json:"success_url,omitempty"`
	// BackURL is used when the payer clicks the back button.
	BackURL string `json:"back_url,omitempty"`
	// ExpiredAt prevents payment after the supplied ISO 8601 date or datetime.
	ExpiredAt string `json:"expired_at,omitempty"`
	// Metadata is merchant-defined key/value data returned in responses and webhooks.
	Metadata Metadata `json:"metadata,omitempty"`
}

// BulkCreateInvoicesRequest creates multiple Moyasar invoices in one request.
type BulkCreateInvoicesRequest struct {
	// Invoices contains up to 50 invoice creation requests.
	Invoices []CreateInvoiceRequest `json:"invoices"`
}

// BulkCreateInvoicesResponse is returned by the bulk invoice creation API.
type BulkCreateInvoicesResponse struct {
	// Invoices contains the created invoice objects.
	Invoices []Invoice `json:"invoices"`
}

// ListInvoicesParams filters and paginates invoice list results.
type ListInvoicesParams struct {
	// Page is the requested page number.
	Page int
	// ID filters by invoice ID.
	ID string
	// Status filters by invoice status.
	Status InvoiceStatus
	// CreatedGT filters invoices created equal to or greater than this timestamp.
	CreatedGT string
	// CreatedLT filters invoices created equal to or less than this timestamp.
	CreatedLT string
	// Metadata filters by metadata key/value pairs.
	Metadata Metadata
}

// InvoiceList is a page of invoices.
type InvoiceList struct {
	// Invoices contains the returned invoice objects.
	Invoices []Invoice `json:"invoices"`
	// Meta contains pagination metadata.
	Meta PageMeta `json:"meta"`
}

// UpdateInvoiceRequest updates mutable invoice metadata.
type UpdateInvoiceRequest struct {
	// Metadata is merchant-defined key/value data returned in responses and webhooks.
	Metadata Metadata `json:"metadata,omitempty"`
}

// Create creates a Moyasar invoice to bill a customer.
func (s *InvoiceService) Create(ctx context.Context, req CreateInvoiceRequest) (*Invoice, error) {
	var invoice Invoice
	if err := s.client.do(ctx, http.MethodPost, "/invoices", nil, req, &invoice); err != nil {
		return nil, err
	}
	return &invoice, nil
}

// BulkCreate creates multiple Moyasar invoices in one request.
func (s *InvoiceService) BulkCreate(ctx context.Context, req BulkCreateInvoicesRequest) (*BulkCreateInvoicesResponse, error) {
	var response BulkCreateInvoicesResponse
	if err := s.client.do(ctx, http.MethodPost, "/invoices/bulk", nil, req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// List returns a paginated list of invoices on the account.
func (s *InvoiceService) List(ctx context.Context, params ListInvoicesParams) (*InvoiceList, error) {
	var list InvoiceList
	if err := s.client.do(ctx, http.MethodGet, "/invoices", params.values(), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Get retrieves a single Moyasar invoice by ID.
func (s *InvoiceService) Get(ctx context.Context, id string) (*Invoice, error) {
	var invoice Invoice
	if err := s.client.do(ctx, http.MethodGet, "/invoices/"+url.PathEscape(id), nil, nil, &invoice); err != nil {
		return nil, err
	}
	return &invoice, nil
}

// Update updates metadata on an existing Moyasar invoice.
func (s *InvoiceService) Update(ctx context.Context, id string, req UpdateInvoiceRequest) (*Invoice, error) {
	var invoice Invoice
	if err := s.client.do(ctx, http.MethodPut, "/invoices/"+url.PathEscape(id), nil, req, &invoice); err != nil {
		return nil, err
	}
	return &invoice, nil
}

// Cancel cancels an invoice so it can no longer be paid by the customer.
func (s *InvoiceService) Cancel(ctx context.Context, id string) (*Invoice, error) {
	var invoice Invoice
	if err := s.client.do(ctx, http.MethodPut, "/invoices/"+url.PathEscape(id)+"/cancel", nil, nil, &invoice); err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (p ListInvoicesParams) values() url.Values {
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
	if p.CreatedGT != "" {
		values.Set("created[gt]", p.CreatedGT)
	}
	if p.CreatedLT != "" {
		values.Set("created[lt]", p.CreatedLT)
	}
	for key, value := range p.Metadata {
		values.Set("metadata["+key+"]", value)
	}
	return values
}
