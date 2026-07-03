package moyasar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// PayoutService provides access to Moyasar payout and payout account APIs.
type PayoutService struct {
	client *Client
}

// PayoutAccountType indicates the payout account type.
type PayoutAccountType string

const (
	PayoutAccountTypeBank   PayoutAccountType = "bank"
	PayoutAccountTypeWallet PayoutAccountType = "wallet"
)

// PayoutStatus is the lifecycle status of a Moyasar payout.
type PayoutStatus string

const (
	PayoutStatusQueued    PayoutStatus = "queued"
	PayoutStatusInitiated PayoutStatus = "initiated"
	PayoutStatusPaid      PayoutStatus = "paid"
	PayoutStatusFailed    PayoutStatus = "failed"
	PayoutStatusCanceled  PayoutStatus = "canceled"
	PayoutStatusReturned  PayoutStatus = "returned"
)

// PayoutChannel is the channel through which a payout is sent.
type PayoutChannel string

const (
	PayoutChannelInternal PayoutChannel = "internal"
	PayoutChannelIPS      PayoutChannel = "ips"
	PayoutChannelSarie    PayoutChannel = "sarie"
)

// PayoutPurpose is the purpose value required when creating a payout.
type PayoutPurpose string

const (
	PayoutPurposeBillsOrRent           PayoutPurpose = "bills_or_rent"
	PayoutPurposeExpensesServices      PayoutPurpose = "expenses_services"
	PayoutPurposePurchaseAssets        PayoutPurpose = "purchase_assets"
	PayoutPurposeSavingInvestment      PayoutPurpose = "saving_investment"
	PayoutPurposeGovernmentDues        PayoutPurpose = "government_dues"
	PayoutPurposeMoneyExchange         PayoutPurpose = "money_exchange"
	PayoutPurposeCreditCardLoan        PayoutPurpose = "credit_card_loan"
	PayoutPurposeGiftOrReward          PayoutPurpose = "gift_or_reward"
	PayoutPurposePersonal              PayoutPurpose = "personal"
	PayoutPurposeInvestmentTransaction PayoutPurpose = "investment_transaction"
	PayoutPurposeFamilyAssistance      PayoutPurpose = "family_assistance"
	PayoutPurposeDonation              PayoutPurpose = "donation"
	PayoutPurposePayrollBenefits       PayoutPurpose = "payroll_benefits"
	PayoutPurposeOnlinePurchase        PayoutPurpose = "online_purchase"
	PayoutPurposeHajjAndUmra           PayoutPurpose = "hajj_and_umra"
	PayoutPurposeDividendPayment       PayoutPurpose = "dividend_payment"
	PayoutPurposeGovernmentPayment     PayoutPurpose = "government_payment"
	PayoutPurposeInvestmentHouse       PayoutPurpose = "investment_house"
	PayoutPurposePaymentToMerchant     PayoutPurpose = "payment_to_merchant"
	PayoutPurposeOwnAccountTransfer    PayoutPurpose = "own_account_transfer"
)

// PayoutObject contains arbitrary provider-specific payout account properties
// or credentials.
type PayoutObject map[string]any

// PayoutAccount is a Moyasar payout account, bank or wallet, used as a source
// or destination for sending payouts.
type PayoutAccount struct {
	// ID is the payout account ID.
	ID string `json:"id"`
	// AccountType indicates whether the payout account is bank or wallet.
	AccountType PayoutAccountType `json:"account_type"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// Properties contains public information relevant to the payout account.
	Properties PayoutObject `json:"properties"`
	// CreatedAt is the time the payout account was created.
	CreatedAt string `json:"created_at"`
}

// CreatePayoutAccountRequest creates a payout account to use as the source or
// destination for sending payouts.
type CreatePayoutAccountRequest struct {
	// AccountType indicates whether the payout account is bank or wallet.
	AccountType PayoutAccountType `json:"account_type"`
	// Properties contains public information relevant to the payout account.
	Properties PayoutObject `json:"properties"`
	// Credentials contains secret information relevant to the payout account.
	Credentials PayoutObject `json:"credentials"`
}

// PayoutAccountList is a page of payout accounts.
type PayoutAccountList struct {
	// PayoutAccounts contains the returned payout account objects.
	PayoutAccounts []PayoutAccount `json:"payout_accounts"`
	// Meta contains pagination metadata.
	Meta PageMeta `json:"meta"`
}

// ListPayoutAccountsParams paginates payout account list results.
type ListPayoutAccountsParams struct {
	// Page is the requested page number.
	Page int
}

// PayoutDestination is implemented by payout destination types.
type PayoutDestination interface {
	payoutDestinationType() string
}

// BankPayoutDestination contains bank beneficiary information.
type BankPayoutDestination struct {
	// IBAN is the IBAN of the beneficiary.
	IBAN string `json:"iban"`
	// Name is the beneficiary's name.
	Name string `json:"name"`
	// Mobile is the beneficiary's mobile.
	Mobile string `json:"mobile"`
	// Country is the beneficiary's country.
	Country string `json:"country"`
	// City is the beneficiary's city.
	City string `json:"city"`
}

func (BankPayoutDestination) payoutDestinationType() string { return "bank" }

func (d BankPayoutDestination) MarshalJSON() ([]byte, error) {
	type alias BankPayoutDestination
	return marshalPayoutDestination("bank", alias(d))
}

// WalletPayoutDestination contains wallet beneficiary information.
type WalletPayoutDestination struct {
	// Mobile is the beneficiary's mobile.
	Mobile string `json:"mobile"`
}

func (WalletPayoutDestination) payoutDestinationType() string { return "wallet" }

func (d WalletPayoutDestination) MarshalJSON() ([]byte, error) {
	type alias WalletPayoutDestination
	return marshalPayoutDestination("wallet", alias(d))
}

// RawPayoutDestination preserves the destination payload returned by Moyasar.
type RawPayoutDestination struct {
	// Type is the destination type, bank or wallet.
	Type string
	// Raw is the full destination JSON object.
	Raw json.RawMessage
}

func (d *RawPayoutDestination) UnmarshalJSON(data []byte) error {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}
	d.Type = probe.Type
	d.Raw = append(d.Raw[:0], data...)
	return nil
}

// Payout is a Moyasar payout sent from a payout account to a bank or wallet
// beneficiary.
type Payout struct {
	// ID is the payout ID.
	ID string `json:"id"`
	// SourceID is the payout account ID used as the source account.
	SourceID string `json:"source_id"`
	// SequenceNumber is the reference number created by you or generated by Moyasar.
	SequenceNumber string `json:"sequence_number"`
	// Channel is the channel through which the payout is sent.
	Channel PayoutChannel `json:"channel"`
	// Status is the payout lifecycle status.
	Status PayoutStatus `json:"status"`
	// Amount is the payout amount in the smallest currency unit.
	Amount int `json:"amount"`
	// Currency is the ISO-4217 three-letter currency code.
	Currency string `json:"currency"`
	// Purpose is the purpose supplied for the payout.
	Purpose PayoutPurpose `json:"purpose"`
	// Comment is the comment provided when creating the payout.
	Comment string `json:"comment"`
	// Destination contains the beneficiary information returned by Moyasar.
	Destination *RawPayoutDestination `json:"destination"`
	// Message is a human-readable message explaining the status.
	Message string `json:"message"`
	// FailureReason is a classification of failure, if any.
	FailureReason string `json:"failure_reason"`
	// CreatedAt is the time the payout was created.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the time the payout was last updated.
	UpdatedAt string `json:"updated_at"`
	// Metadata is merchant-defined key/value data returned in responses and webhooks.
	Metadata Metadata `json:"metadata"`
}

// CreatePayoutRequest sends a payout from a payout account to a bank or wallet
// beneficiary.
type CreatePayoutRequest struct {
	// SourceID is the payout account ID to use as the source account.
	SourceID string `json:"source_id"`
	// SequenceNumber is an optional 16 digit reference. Moyasar generates one if omitted.
	SequenceNumber string `json:"sequence_number,omitempty"`
	// Amount is the payout amount in the smallest currency unit.
	Amount int `json:"amount"`
	// Purpose is the required payout purpose.
	Purpose PayoutPurpose `json:"purpose"`
	// Destination contains the beneficiary information.
	Destination PayoutDestination `json:"destination"`
	// Comment is an optional comment for the payout.
	Comment string `json:"comment,omitempty"`
	// Metadata is merchant-defined key/value data returned in responses and webhooks.
	Metadata Metadata `json:"metadata,omitempty"`
}

// BulkPayoutItem is one payout item in a bulk payout request.
type BulkPayoutItem struct {
	// SequenceNumber is an optional 16 digit reference. Moyasar generates one if omitted.
	SequenceNumber string `json:"sequence_number,omitempty"`
	// Amount is the payout amount in the smallest currency unit.
	Amount int `json:"amount"`
	// Purpose is the required payout purpose.
	Purpose PayoutPurpose `json:"purpose"`
	// Destination contains the beneficiary information.
	Destination PayoutDestination `json:"destination"`
	// Comment is an optional comment for the payout.
	Comment string `json:"comment,omitempty"`
	// Metadata is merchant-defined key/value data returned in responses and webhooks.
	Metadata Metadata `json:"metadata,omitempty"`
}

// CreateBulkPayoutRequest sends multiple payouts in a single request.
type CreateBulkPayoutRequest struct {
	// SourceID is the payout account ID to use as the source account.
	SourceID string `json:"source_id"`
	// Payouts contains the payouts to create.
	Payouts []BulkPayoutItem `json:"payouts"`
}

// BulkPayoutResponse contains payouts created by the bulk payout API.
type BulkPayoutResponse struct {
	// Payouts contains the created payout objects.
	Payouts []Payout `json:"payouts"`
}

// PayoutList is a page of payouts.
type PayoutList struct {
	// Payouts contains the returned payout objects.
	Payouts []Payout `json:"payouts"`
	// Meta contains pagination metadata.
	Meta PageMeta `json:"meta"`
}

// ListPayoutsParams paginates payout list results.
type ListPayoutsParams struct {
	// Page is the requested page number.
	Page int
}

// CreateAccount creates a Moyasar payout account.
func (s *PayoutService) CreateAccount(ctx context.Context, req CreatePayoutAccountRequest) (*PayoutAccount, error) {
	var account PayoutAccount
	if err := s.client.do(ctx, http.MethodPost, "/payout_accounts", nil, req, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// ListAccounts lists payout accounts on the account with pagination.
func (s *PayoutService) ListAccounts(ctx context.Context, params ListPayoutAccountsParams) (*PayoutAccountList, error) {
	var list PayoutAccountList
	if err := s.client.do(ctx, http.MethodGet, "/payout_accounts", params.values(), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// GetAccount retrieves a single payout account by ID.
func (s *PayoutService) GetAccount(ctx context.Context, id string) (*PayoutAccount, error) {
	var account PayoutAccount
	if err := s.client.do(ctx, http.MethodGet, "/payout_accounts/"+url.PathEscape(id), nil, nil, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// Create sends a Moyasar payout from a payout account to a beneficiary.
func (s *PayoutService) Create(ctx context.Context, req CreatePayoutRequest) (*Payout, error) {
	var payout Payout
	if err := s.client.do(ctx, http.MethodPost, "/payouts", nil, req, &payout); err != nil {
		return nil, err
	}
	return &payout, nil
}

// List lists payouts on the account with pagination.
func (s *PayoutService) List(ctx context.Context, params ListPayoutsParams) (*PayoutList, error) {
	var list PayoutList
	if err := s.client.do(ctx, http.MethodGet, "/payouts", params.values(), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Get retrieves a single payout by ID.
func (s *PayoutService) Get(ctx context.Context, id string) (*Payout, error) {
	var payout Payout
	if err := s.client.do(ctx, http.MethodGet, "/payout/"+url.PathEscape(id), nil, nil, &payout); err != nil {
		return nil, err
	}
	return &payout, nil
}

// BulkCreate sends multiple Moyasar payouts in a single request.
func (s *PayoutService) BulkCreate(ctx context.Context, req CreateBulkPayoutRequest) (*BulkPayoutResponse, error) {
	var response BulkPayoutResponse
	if err := s.client.do(ctx, http.MethodPost, "/payouts/bulk", nil, req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (p ListPayoutAccountsParams) values() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	return values
}

func (p ListPayoutsParams) values() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	return values
}

func marshalPayoutDestination(destinationType string, destination any) ([]byte, error) {
	data, err := json.Marshal(destination)
	if err != nil {
		return nil, err
	}
	var fields map[string]any
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	fields["type"] = destinationType
	return json.Marshal(fields)
}
