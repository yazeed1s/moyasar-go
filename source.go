package moyasar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// SourceService provides access to Moyasar payment source APIs.
type SourceService struct {
	client *Client
}

// IssuerSource is implemented by source types accepted by Retrieve Issuer.
type IssuerSource interface {
	issuerSourceType() string
}

// IssuerCreditCardSource is a simplified credit card source used for issuer lookup.
type IssuerCreditCardSource struct {
	// Number is the card number without separators.
	Number string `json:"number"`
}

func (IssuerCreditCardSource) issuerSourceType() string { return "creditcard" }

func (s IssuerCreditCardSource) MarshalJSON() ([]byte, error) {
	type alias IssuerCreditCardSource
	return marshalIssuerSource("creditcard", alias(s))
}

// IssuerApplePaySource is a device-payment source used for issuer lookup.
type IssuerApplePaySource struct {
	// Token is the encrypted Apple Pay token to be processed.
	Token string `json:"token"`
}

func (IssuerApplePaySource) issuerSourceType() string { return "applepay" }

func (s IssuerApplePaySource) MarshalJSON() ([]byte, error) {
	type alias IssuerApplePaySource
	return marshalIssuerSource("applepay", alias(s))
}

// RetrieveIssuerRequest looks up card issuing bank and card metadata without
// creating a payment.
type RetrieveIssuerRequest struct {
	// Source is the simplified payment source to inspect.
	Source IssuerSource `json:"source"`
}

// Issuer contains issuing bank and card metadata inferred by Moyasar.
type Issuer struct {
	// IssuerName is the name of the card issuing bank.
	IssuerName string `json:"issuer_name"`
	// IssuerCountry is the origin country of the card issuer as ISO 3166 alpha-2.
	IssuerCountry string `json:"issuer_country"`
	// IssuerCardType is debit, credit, charge_card, or unspecified.
	IssuerCardType string `json:"issuer_card_type"`
	// IssuerCardCategory is the card category or product type.
	IssuerCardCategory string `json:"issuer_card_category"`
	// Company is the scheme through which the payment is processed.
	Company string `json:"company"`
	// FirstDigits are the first 6 to 8 digits of the card BIN/IIN.
	FirstDigits string `json:"first_digits"`
	// LastDigits are the last 4 digits when a full card number is supplied.
	LastDigits string `json:"last_digits"`
}

// ExtendedPaymentSource preserves the extended payment source payload returned
// by Moyasar.
//
// The response shape depends on the payment source type and can include
// operational details such as acquirer references, gateway IDs, reconciliation
// dates, and network transaction IDs.
type ExtendedPaymentSource struct {
	// Type is the source type, such as creditcard, applepay, samsungpay, or stcpay.
	Type string
	// Raw is the full extended source JSON object.
	Raw json.RawMessage
}

func (s *ExtendedPaymentSource) UnmarshalJSON(data []byte) error {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}
	s.Type = probe.Type
	s.Raw = append(s.Raw[:0], data...)
	return nil
}

// RetrieveIssuer looks up the issuing bank and card metadata for a source.
func (s *SourceService) RetrieveIssuer(ctx context.Context, req RetrieveIssuerRequest) (*Issuer, error) {
	var issuer Issuer
	if err := s.client.do(ctx, http.MethodPost, "/source/issuer", nil, req, &issuer); err != nil {
		return nil, err
	}
	return &issuer, nil
}

// GetPaymentSource fetches the extended source object for a payment.
func (s *SourceService) GetPaymentSource(ctx context.Context, paymentID string) (*ExtendedPaymentSource, error) {
	var source ExtendedPaymentSource
	if err := s.client.do(ctx, http.MethodGet, "/payments/"+url.PathEscape(paymentID)+"/source", nil, nil, &source); err != nil {
		return nil, err
	}
	return &source, nil
}

func marshalIssuerSource(sourceType string, source any) ([]byte, error) {
	data, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}
	var fields map[string]any
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	fields["type"] = sourceType
	return json.Marshal(fields)
}
