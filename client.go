package moyasar

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL     = "https://api.moyasar.com/v1"
	defaultTransferURL = "https://apimig.moyasar.com/v1"
)

// Client is a Moyasar API client.
//
// The client authenticates requests using HTTP Basic Auth with the API key as
// the username and an empty password, as required by Moyasar.
type Client struct {
	apiKey      string
	baseURL     *url.URL
	transferURL *url.URL
	httpClient  *http.Client

	// Payments provides access to Moyasar payment APIs.
	Payments *PaymentService
	// Invoices provides access to Moyasar invoice APIs.
	Invoices *InvoiceService
	// Tokens provides access to Moyasar token APIs.
	Tokens *TokenService
	// Sources provides access to Moyasar payment source APIs.
	Sources *SourceService
	// CardAuths provides access to Moyasar standalone 3D Secure authentication APIs.
	CardAuths *CardAuthService
	// Webhooks provides access to Moyasar webhook registration and delivery APIs.
	Webhooks *WebhookService
	// Payouts provides access to Moyasar payout and payout account APIs.
	Payouts *PayoutService
	// Settlements provides access to Moyasar settlement APIs.
	Settlements *SettlementService
	// InternalTransactions provides access to Moyasar internal transaction APIs.
	InternalTransactions *InternalTransactionService
	// Transfers provides access to aggregation merchant transfer APIs.
	Transfers *TransferService
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets the HTTP client used for requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithBaseURL sets the base URL used for the main Moyasar API.
//
// The default is https://api.moyasar.com/v1.
func WithBaseURL(rawURL string) Option {
	return func(c *Client) {
		if u := mustParseBaseURL(rawURL); u != nil {
			c.baseURL = u
		}
	}
}

// WithTransferURL sets the base URL used for the transfers API.
//
// The default is https://apimig.moyasar.com/v1.
func WithTransferURL(rawURL string) Option {
	return func(c *Client) {
		if u := mustParseBaseURL(rawURL); u != nil {
			c.transferURL = u
		}
	}
}

// NewClient creates a Moyasar API client using HTTP Basic Auth.
//
// Secret keys should be used for backend operations. Publishable keys are only
// intended for frontend-safe operations such as payment creation and
// tokenization flows.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:      apiKey,
		baseURL:     mustParseBaseURL(defaultBaseURL),
		transferURL: mustParseBaseURL(defaultTransferURL),
		httpClient:  http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.Payments = &PaymentService{client: c}
	c.Invoices = &InvoiceService{client: c}
	c.Tokens = &TokenService{client: c}
	c.Sources = &SourceService{client: c}
	c.CardAuths = &CardAuthService{client: c}
	c.Webhooks = &WebhookService{client: c}
	c.Payouts = &PayoutService{client: c}
	c.Settlements = &SettlementService{client: c}
	c.InternalTransactions = &InternalTransactionService{client: c}
	c.Transfers = &TransferService{client: c}
	return c
}

func mustParseBaseURL(rawURL string) *url.URL {
	if rawURL == "" {
		return nil
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	u.Path = strings.TrimRight(u.Path, "/")
	return u
}
