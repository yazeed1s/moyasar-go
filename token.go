package moyasar

import (
	"context"
	"net/http"
	"net/url"
)

// TokenService provides access to Moyasar token APIs.
type TokenService struct {
	client *Client
}

// TokenStatus is the lifecycle status of a Moyasar card token.
type TokenStatus string

const (
	TokenStatusInitiated TokenStatus = "initiated"
	TokenStatusActive    TokenStatus = "active"
	TokenStatusInactive  TokenStatus = "inactive"
)

// Token is a Moyasar card token object.
type Token struct {
	// ID is the token's unique ID.
	ID string `json:"id"`
	// Status is the token status. The default status is initiated.
	Status TokenStatus `json:"status"`
	// Brand is the card brand, such as visa, master, mada, amex, or unionpay.
	Brand string `json:"brand"`
	// Funding is the card funding type, such as credit or debit.
	Funding string `json:"funding"`
	// Country is the card issuing country.
	Country string `json:"country"`
	// Month is the card expiration month.
	Month string `json:"month"`
	// Year is the card expiration year.
	Year string `json:"year"`
	// Name is the cardholder name.
	Name string `json:"name"`
	// LastFour is the card's last four digits.
	LastFour string `json:"last_four"`
	// Metadata is merchant-defined key/value data.
	Metadata Metadata `json:"metadata"`
	// Message is a human-readable message returned by Moyasar.
	Message *string `json:"message"`
	// VerificationURL is the 3D Secure verification process URL.
	VerificationURL *string `json:"verification_url"`
	// CreatedAt is the token creation timestamp in ISO 8601 format.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the token modification timestamp in ISO 8601 format.
	UpdatedAt string `json:"updated_at"`
}

// CreateTokenRequest generates a token for a Mada or credit card.
//
// Moyasar documents this request as a frontend request made directly to Moyasar
// with a publishable key.
type CreateTokenRequest struct {
	// Name is the cardholder name.
	Name string
	// Number is the card number without separators.
	Number string
	// Month is the card expiration month.
	Month string
	// Year is the card expiration year.
	Year string
	// CVC is the card security code.
	CVC string
	// CallbackURL returns the payer after token verification.
	CallbackURL string
	// Metadata is merchant-defined key/value data.
	Metadata Metadata
}

// Create generates a token for a given Mada or credit card.
func (s *TokenService) Create(ctx context.Context, req CreateTokenRequest) (*Token, error) {
	var token Token
	if err := s.client.doForm(ctx, http.MethodPost, "/tokens", req.values(), &token); err != nil {
		return nil, err
	}
	return &token, nil
}

// Get fetches an individual token by its unique ID.
func (s *TokenService) Get(ctx context.Context, id string) (*Token, error) {
	var token Token
	if err := s.client.do(ctx, http.MethodGet, "/tokens/"+url.PathEscape(id), nil, nil, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

// Delete deletes an individual token by its unique ID.
func (s *TokenService) Delete(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodDelete, "/tokens/"+url.PathEscape(id), nil, nil, nil)
}

func (r CreateTokenRequest) values() url.Values {
	values := url.Values{}
	values.Set("name", r.Name)
	values.Set("number", r.Number)
	values.Set("month", r.Month)
	values.Set("year", r.Year)
	values.Set("cvc", r.CVC)
	if r.CallbackURL != "" {
		values.Set("callback_url", r.CallbackURL)
	}
	for key, value := range r.Metadata {
		values.Set("metadata["+key+"]", value)
	}
	return values
}
