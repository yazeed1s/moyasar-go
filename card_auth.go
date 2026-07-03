package moyasar

import (
	"context"
	"net/http"
	"net/url"
)

// CardAuthService provides access to Moyasar standalone 3D Secure
// authentication APIs.
type CardAuthService struct {
	client *Client
}

// CardAuthStatus is the lifecycle status of a standalone card authentication.
type CardAuthStatus string

const (
	CardAuthStatusInitiated     CardAuthStatus = "initiated"
	CardAuthStatusAvailable     CardAuthStatus = "available"
	CardAuthStatusInProgress    CardAuthStatus = "in_progress"
	CardAuthStatusAuthenticated CardAuthStatus = "authenticated"
	CardAuthStatusFailed        CardAuthStatus = "failed"
)

// CardAuthSource is the credit card source to authenticate without charging it.
type CardAuthSource struct {
	// Name is the cardholder name.
	Name string `json:"name"`
	// Number is the card number without separators.
	Number string `json:"number"`
	// Month is the card expiry month.
	Month string `json:"month"`
	// Year is the card expiry year.
	Year string `json:"year"`
	// CVC is the card security code.
	CVC string `json:"cvc"`
}

func (s CardAuthSource) MarshalJSON() ([]byte, error) {
	type alias CardAuthSource
	return marshalPaymentSource("creditcard", alias(s))
}

// CreateCardAuthRequest starts a standalone 3D Secure authentication for a
// card, without charging it.
type CreateCardAuthRequest struct {
	// Amount is the intended payment amount in the smallest currency unit.
	Amount int `json:"amount"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// CallbackURL returns the cardholder after completing authentication.
	CallbackURL string `json:"callback_url"`
	// Source is the credit card source to authenticate.
	Source CardAuthSource `json:"source"`
}

// CardAuth is a standalone 3D Secure authentication object.
type CardAuth struct {
	// ID is the unique identifier of the card authentication.
	ID string `json:"id"`
	// Status indicates the authentication lifecycle state.
	Status CardAuthStatus `json:"status"`
	// Amount is the intended payment amount in the smallest currency unit.
	Amount int `json:"amount"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// CallbackURL returns the cardholder after completing authentication.
	CallbackURL string `json:"callback_url"`
	// TransactionURL is where the cardholder is redirected to complete the challenge.
	TransactionURL *string `json:"transaction_url"`
	// Card contains masked card details.
	Card CardAuthCard `json:"card"`
	// Result is populated once authentication reaches authenticated or failed.
	// It is returned only with secret key authentication.
	Result *CardAuthResult `json:"result"`
	// Message is a human-readable message, set on failures.
	Message *string `json:"message"`
	// CreatedAt is the time the card authentication was created.
	CreatedAt string `json:"created_at"`
}

// CardAuthCard contains masked card details for a card authentication.
type CardAuthCard struct {
	// Company is the scheme through which the payment is processed.
	Company string `json:"company"`
	// LastDigits are the masked card number digits returned by Moyasar.
	LastDigits string `json:"last_digits"`
}

// CardAuthResult contains 3D Secure authentication values.
type CardAuthResult struct {
	// ECI is the Electronic Commerce Indicator.
	ECI string `json:"eci"`
	// AuthenticationValue is the CAVV/AAV value, Base64 encoded.
	AuthenticationValue string `json:"authentication_value"`
	// DSTransactionID is the Directory Server transaction ID.
	DSTransactionID string `json:"ds_transaction_id"`
	// Version is the 3DS protocol version.
	Version string `json:"version"`
	// TransactionStatus is the EMVCo transaction status.
	TransactionStatus string `json:"transaction_status"`
	// TransactionStatusReason is the EMVCo transaction status reason code.
	TransactionStatusReason *string `json:"transaction_status_reason"`
	// AuthScheme is the directory server used to perform authentication.
	AuthScheme string `json:"auth_scheme"`
	// ACSTransactionID is the ACS transaction ID.
	ACSTransactionID string `json:"acs_transaction_id"`
	// DSReferenceNumber is the Directory Server reference number.
	DSReferenceNumber string `json:"ds_reference_number"`
	// ACSReferenceNumber is the ACS reference number.
	ACSReferenceNumber string `json:"acs_reference_number"`
	// ThreeDSServerTransactionID is the 3DS Server transaction ID.
	ThreeDSServerTransactionID string `json:"three_ds_server_transaction_id"`
	// IsFrictionless is true if authenticated without a challenge.
	IsFrictionless bool `json:"is_frictionless"`
}

// Create starts a standalone 3D Secure authentication for a card.
//
// When the returned status is available, redirect the cardholder to
// transaction_url to complete the challenge.
func (s *CardAuthService) Create(ctx context.Context, req CreateCardAuthRequest) (*CardAuth, error) {
	var auth CardAuth
	if err := s.client.do(ctx, http.MethodPost, "/card_auths", nil, req, &auth); err != nil {
		return nil, err
	}
	return &auth, nil
}

// Get retrieves a card authentication and its result once the cardholder has
// finished.
func (s *CardAuthService) Get(ctx context.Context, id string) (*CardAuth, error) {
	var auth CardAuth
	if err := s.client.do(ctx, http.MethodGet, "/card_auths/"+url.PathEscape(id), nil, nil, &auth); err != nil {
		return nil, err
	}
	return &auth, nil
}
