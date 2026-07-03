package moyasar

// Test card brands used by Moyasar sandbox.
const (
	TestCardBrandMada       = "mada"
	TestCardBrandVisa       = "visa"
	TestCardBrandMastercard = "mastercard"
	TestCardBrandAmex       = "amex"
	TestCardBrandUnionPay   = "unionpay"
)

// Common Moyasar sandbox test cards.
const (
	TestCardMadaPaid                = "4201320111111010"
	TestCardVisaPaid                = "4111111111111111"
	TestCardVisaFrictionlessPaid    = "4111114005765430"
	TestCardMastercardPaid          = "5421080101000000"
	TestCardAmexPaid                = "340000000900000"
	TestCardUnionPayPaid            = "6200000000000005"
	TestCardVisaInsufficientFunds   = "4123120001090000"
	TestCardMastercardDeclined      = "5204730000002514"
	TestCardMadaStolenCard          = "4201321144311528"
	TestCardAmexExpiredCard         = "340000018441278"
	TestCardUnionPayWithdrawalLimit = "6200000000000062"
)

// TestCard is a Moyasar sandbox card scenario.
//
// These cards are only for Moyasar sandbox. Moyasar documents that using any
// card not listed in the sandbox card table will result in a failed payment.
type TestCard struct {
	// Number is the sandbox card number.
	Number string
	// Brand is the card scheme, such as visa, mada, master, amex, or unionpay.
	Brand string
	// ExpectedStatus is the payment status Moyasar documents for this card.
	ExpectedStatus PaymentStatus
	// Message is the documented sandbox response message.
	Message string
	// ResponseCode is the documented issuer response code when available.
	ResponseCode string
	// ThreeDSNote is the documented 3DS behavior for this card.
	ThreeDSNote string
}

// TestCards contains Moyasar sandbox card scenarios from the test card docs.
var TestCards = []TestCard{
	{Number: "4201320111111010", Brand: TestCardBrandMada, ExpectedStatus: PaymentStatusPaid, Message: "APPROVED", ResponseCode: "00"},
	{Number: "4201320000013020", Brand: TestCardBrandMada, ExpectedStatus: PaymentStatusFailed, Message: "UNSPECIFIED FAILURE", ResponseCode: "99"},
	{Number: "4201320000311101", Brand: TestCardBrandMada, ExpectedStatus: PaymentStatusFailed, Message: "INSUFFICIENT FUNDS", ResponseCode: "51"},
	{Number: "4201320131000508", Brand: TestCardBrandMada, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: LOST CARD", ResponseCode: "41"},
	{Number: "4201321234411220", Brand: TestCardBrandMada, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED", ResponseCode: "05"},
	{Number: "4201322267774310", Brand: TestCardBrandMada, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: EXPIRED CARD", ResponseCode: "54"},
	{Number: "4201326324640570", Brand: TestCardBrandMada, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: EXCEEDS WITHDRAWAL LIMIT", ResponseCode: "61"},
	{Number: "4201321144311528", Brand: TestCardBrandMada, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: STOLEN CARD", ResponseCode: "43"},

	{Number: "4111118250252531", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "3DS: attempted but not available, please ensure that you have enabled Online Purchase from your bank portal.", ThreeDSNote: "ECI 06"},
	{Number: "4111114005765430", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusPaid, Message: "APPROVED", ResponseCode: "00", ThreeDSNote: "Frictionless Authentication"},
	{Number: "4111111111111111", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusPaid, Message: "APPROVED", ResponseCode: "00"},
	{Number: "4111113343111067", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "3DS service error occurred.", ThreeDSNote: "3DS fails during enrollement check"},
	{Number: "4111116611600661", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "The card is not enrolled in 3DS service."},
	{Number: "4111112205628150", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "3DS service error occurred.", ThreeDSNote: "3DS fails during authentication attempt"},
	{Number: "4111115784228433", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "The authentication attempt was rejected by the issuer bank."},
	{Number: "4111115620358287", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "The authentication is unavailable, please try again later or contact issuer bank if problem persisted."},
	{Number: "4123120000000000", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "UNSPECIFIED FAILURE", ResponseCode: "99"},
	{Number: "4123120001090000", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "INSUFFICIENT FUNDS", ResponseCode: "51"},
	{Number: "4123450131000508", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: LOST CARD", ResponseCode: "41"},
	{Number: "4123120001090109", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED", ResponseCode: "05"},
	{Number: "4123128518640738", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: EXPIRED CARD", ResponseCode: "54"},
	{Number: "4123123033308648", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: EXCEEDS WITHDRAWAL LIMIT", ResponseCode: "61"},
	{Number: "4123125276780003", Brand: TestCardBrandVisa, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: STOLEN CARD", ResponseCode: "43"},

	{Number: "5421080101000000", Brand: TestCardBrandMastercard, ExpectedStatus: PaymentStatusPaid, Message: "APPROVED", ResponseCode: "00"},
	{Number: "5105105105105100", Brand: TestCardBrandMastercard, ExpectedStatus: PaymentStatusFailed, Message: "UNSPECIFIED FAILURE", ResponseCode: "99"},
	{Number: "5457210001000092", Brand: TestCardBrandMastercard, ExpectedStatus: PaymentStatusFailed, Message: "INSUFFICIENT FUNDS", ResponseCode: "51"},
	{Number: "5204010101000000", Brand: TestCardBrandMastercard, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: LOST CARD", ResponseCode: "41"},
	{Number: "5204730000002514", Brand: TestCardBrandMastercard, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED", ResponseCode: "05"},
	{Number: "5105107550274126", Brand: TestCardBrandMastercard, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: EXPIRED CARD", ResponseCode: "54"},
	{Number: "5105106475101067", Brand: TestCardBrandMastercard, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: EXCEEDS WITHDRAWAL LIMIT", ResponseCode: "61"},
	{Number: "5105107304607225", Brand: TestCardBrandMastercard, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: STOLEN CARD", ResponseCode: "43"},

	{Number: "340000000900000", Brand: TestCardBrandAmex, ExpectedStatus: PaymentStatusPaid, Message: "APPROVED", ResponseCode: "00"},
	{Number: "371111111111114", Brand: TestCardBrandAmex, ExpectedStatus: PaymentStatusFailed, Message: "UNSPECIFIED FAILURE", ResponseCode: "99"},
	{Number: "340033000000000", Brand: TestCardBrandAmex, ExpectedStatus: PaymentStatusFailed, Message: "INSUFFICIENT FUNDS", ResponseCode: "51"},
	{Number: "340012340501000", Brand: TestCardBrandAmex, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: LOST CARD", ResponseCode: "41"},
	{Number: "340033000000133", Brand: TestCardBrandAmex, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED", ResponseCode: "05"},
	{Number: "340000018441278", Brand: TestCardBrandAmex, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: EXPIRED CARD", ResponseCode: "54"},
	{Number: "340000753060788", Brand: TestCardBrandAmex, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: EXCEEDS WITHDRAWAL LIMIT", ResponseCode: "61"},
	{Number: "340000418501838", Brand: TestCardBrandAmex, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: STOLEN CARD", ResponseCode: "43"},

	{Number: "6200000000000005", Brand: TestCardBrandUnionPay, ExpectedStatus: PaymentStatusPaid, Message: "APPROVED", ResponseCode: "00"},
	{Number: "6200000000000013", Brand: TestCardBrandUnionPay, ExpectedStatus: PaymentStatusFailed, Message: "UNSPECIFIED FAILURE", ResponseCode: "99"},
	{Number: "6200000000000021", Brand: TestCardBrandUnionPay, ExpectedStatus: PaymentStatusFailed, Message: "INSUFFICIENT FUNDS", ResponseCode: "51"},
	{Number: "6200000000000039", Brand: TestCardBrandUnionPay, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: LOST CARD", ResponseCode: "41"},
	{Number: "6200000000000047", Brand: TestCardBrandUnionPay, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED", ResponseCode: "05"},
	{Number: "6200000000000054", Brand: TestCardBrandUnionPay, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: EXPIRED CARD", ResponseCode: "54"},
	{Number: "6200000000000062", Brand: TestCardBrandUnionPay, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: EXCEEDS WITHDRAWAL LIMIT", ResponseCode: "61"},
	{Number: "6200000000000070", Brand: TestCardBrandUnionPay, ExpectedStatus: PaymentStatusFailed, Message: "DECLINED: STOLEN CARD", ResponseCode: "43"},
}
