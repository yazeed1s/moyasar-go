package moyasar

import "testing"

func TestTestCardsIncludeCommonCards(t *testing.T) {
	t.Parallel()

	if len(TestCards) == 0 {
		t.Fatal("TestCards is empty")
	}

	foundVisaPaid := false
	foundVisaInsufficientFunds := false
	for _, card := range TestCards {
		if card.Number == TestCardVisaPaid && card.ExpectedStatus == PaymentStatusPaid {
			foundVisaPaid = true
		}
		if card.Number == TestCardVisaInsufficientFunds && card.ExpectedStatus == PaymentStatusFailed && card.ResponseCode == "51" {
			foundVisaInsufficientFunds = true
		}
	}

	if !foundVisaPaid {
		t.Fatal("TestCards does not include TestCardVisaPaid")
	}
	if !foundVisaInsufficientFunds {
		t.Fatal("TestCards does not include TestCardVisaInsufficientFunds")
	}
}
