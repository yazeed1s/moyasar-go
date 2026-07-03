package moyasar

import (
	"encoding/json"
)

// PaymentSource is implemented by request source types accepted by Create Payment.
//
// The source object defines how Moyasar starts the payment, such as card,
// token, Apple Pay, Samsung Pay, or STC Pay.
type PaymentSource interface {
	paymentSourceType() string
}

// RawPaymentSource preserves the source payload returned by Moyasar.
//
// Response source fields vary by payment method and status. Type contains the
// source discriminator and Raw contains the full JSON object returned by Moyasar.
type RawPaymentSource struct {
	// Type is the returned source type, such as creditcard, token, applepay,
	// samsungpay, or stcpay.
	Type string
	// Raw is the full source JSON object.
	Raw json.RawMessage
}

func (s *RawPaymentSource) UnmarshalJSON(data []byte) error {
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

// CardAuthData contains bring-your-own 3DS values for a credit card payment.
//
// Providing these values disables Moyasar's own 3DS flow for the payment. This
// feature is enabled only for selected merchants.
type CardAuthData struct {
	// Provider identifies where the 3DS data originated.
	Provider string `json:"provider"`
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
	// AuthScheme is the directory server used, such as visa, mada, or master.
	AuthScheme string `json:"auth_scheme,omitempty"`
	// ACSTransactionID is the ACS transaction ID.
	ACSTransactionID string `json:"acs_transaction_id,omitempty"`
	// DSReferenceNumber is the Directory Server reference number.
	DSReferenceNumber string `json:"ds_reference_number,omitempty"`
	// ACSReferenceNumber is the ACS reference number.
	ACSReferenceNumber string `json:"acs_reference_number,omitempty"`
}

// CreditCardSource charges a raw Mada or credit card.
//
// Cardholder data should be sent directly from the customer device to Moyasar
// when using publishable-key frontend flows. Do not pass cardholder data through
// your backend unless your integration is allowed to handle it.
type CreditCardSource struct {
	// Name is the cardholder name. Moyasar requires at least two English names.
	Name string `json:"name"`
	// Number is the card number without separators.
	Number string `json:"number"`
	// Month is the card expiry month from 1 to 12.
	Month int `json:"month"`
	// Year is the card expiry year.
	Year int `json:"year"`
	// CVC is the card security code.
	CVC string `json:"cvc"`
	// StatementDescriptor adds extra issuer statement descriptor information.
	StatementDescriptor string `json:"statement_descriptor,omitempty"`
	// ThreeDS controls whether 3DS is used. Moyasar defaults it to true.
	ThreeDS *bool `json:"3ds,omitempty"`
	// Manual authorizes the payment without capturing it immediately.
	Manual *bool `json:"manual,omitempty"`
	// SaveCard asks Moyasar to generate a token in source.token when supported.
	SaveCard *bool `json:"save_card,omitempty"`
	// CardAuthID reuses an authenticated standalone card_auth for this payment.
	CardAuthID string `json:"card_auth_id,omitempty"`
	// CardAuthData supplies external 3DS authentication values.
	CardAuthData *CardAuthData `json:"card_auth_data,omitempty"`
}

func (CreditCardSource) paymentSourceType() string { return "creditcard" }

func (s CreditCardSource) MarshalJSON() ([]byte, error) {
	type alias CreditCardSource
	return marshalPaymentSource("creditcard", alias(s))
}

// TokenSource charges a previously created Moyasar card token.
type TokenSource struct {
	// Token is the Moyasar token ID. It starts with token_.
	Token string `json:"token"`
	// CVC is the optional card security code for the tokenized card.
	CVC string `json:"cvc,omitempty"`
	// StatementDescriptor adds extra issuer statement descriptor information.
	StatementDescriptor string `json:"statement_descriptor,omitempty"`
	// ThreeDS controls whether 3DS is used for the token payment.
	ThreeDS *bool `json:"3ds,omitempty"`
	// Manual authorizes the payment without capturing it immediately.
	Manual *bool `json:"manual,omitempty"`
}

func (TokenSource) paymentSourceType() string { return "token" }

func (s TokenSource) MarshalJSON() ([]byte, error) {
	type alias TokenSource
	return marshalPaymentSource("token", alias(s))
}

// ApplePaySource starts an Apple Pay payment.
//
// Token contains the encrypted Apple Pay token payload. Device PAN fields are
// also modeled for integrations that provide decrypted network token data.
type ApplePaySource struct {
	// Token is the encrypted Apple Pay token payload.
	Token string `json:"token,omitempty"`
	// Number is the Device Primary Account Number.
	Number string `json:"number,omitempty"`
	// Month is the card expiry month from the device token.
	Month int `json:"month,omitempty"`
	// Year is the card expiry year from the device token.
	Year int `json:"year,omitempty"`
	// Cryptogram is the network token cryptogram generated by the device wallet.
	Cryptogram string `json:"cryptogram,omitempty"`
	// DeviceID is the unique identifier assigned to the device or wallet.
	DeviceID string `json:"device_id,omitempty"`
	// LastFour is the last four digits of the card.
	LastFour string `json:"last_four,omitempty"`
	// ECI is the Electronic Commerce Indicator.
	ECI string `json:"eci,omitempty"`
	// Manual authorizes the payment without capturing it immediately.
	Manual *bool `json:"manual,omitempty"`
	// SaveCard asks Moyasar to generate a token in source.token when supported.
	SaveCard *bool `json:"save_card,omitempty"`
	// StatementDescriptor adds extra issuer statement descriptor information.
	StatementDescriptor string `json:"statement_descriptor,omitempty"`
}

func (ApplePaySource) paymentSourceType() string { return "applepay" }

func (s ApplePaySource) MarshalJSON() ([]byte, error) {
	type alias ApplePaySource
	return marshalPaymentSource("applepay", alias(s))
}

// SamsungPaySource starts a Samsung Pay payment.
type SamsungPaySource struct {
	// Token is the encrypted Samsung Pay token payload.
	Token string `json:"token,omitempty"`
	// Manual authorizes the payment without capturing it immediately.
	Manual *bool `json:"manual,omitempty"`
	// SaveCard asks Moyasar to generate a token in source.token when supported.
	SaveCard *bool `json:"save_card,omitempty"`
	// StatementDescriptor adds extra issuer statement descriptor information.
	StatementDescriptor string `json:"statement_descriptor,omitempty"`
}

func (SamsungPaySource) paymentSourceType() string { return "samsungpay" }

func (s SamsungPaySource) MarshalJSON() ([]byte, error) {
	type alias SamsungPaySource
	return marshalPaymentSource("samsungpay", alias(s))
}

// STCPaySource starts an STC Pay payment.
type STCPaySource struct {
	// Mobile is a Saudi mobile number accepted by STC Pay.
	Mobile string `json:"mobile"`
	// Cashier is an optional cashier identifier shown in the dashboard.
	Cashier string `json:"cashier,omitempty"`
	// Branch is an optional branch identifier shown in the dashboard.
	Branch string `json:"branch,omitempty"`
}

func (STCPaySource) paymentSourceType() string { return "stcpay" }

func (s STCPaySource) MarshalJSON() ([]byte, error) {
	type alias STCPaySource
	return marshalPaymentSource("stcpay", alias(s))
}

func marshalPaymentSource(sourceType string, source any) ([]byte, error) {
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
