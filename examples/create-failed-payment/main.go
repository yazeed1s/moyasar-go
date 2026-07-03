package main

import (
	"context"
	"fmt"
	"log"
	"os"

	moyasar "github.com/yazeed1s/moyasar-go"
)

func main() {
	apiKey := os.Getenv("MOYASAR_API_KEY")
	if apiKey == "" {
		log.Fatal("set MOYASAR_API_KEY")
	}

	client := moyasar.NewClient(apiKey)

	payment, err := client.Payments.Create(context.Background(), moyasar.CreatePaymentRequest{
		Amount:      100,
		Currency:    "SAR",
		Description: "Failed payment test",
		CallbackURL: "https://example.com/callback",
		Source: moyasar.CreditCardSource{
			Name:   "John Doe",
			Number: moyasar.TestCardVisaInsufficientFunds,
			Month:  9,
			Year:   2030,
			CVC:    "123",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("payment id:", payment.ID)
	fmt.Println("status:", payment.Status)
	if payment.Source != nil {
		fmt.Println("source:", string(payment.Source.Raw))
	}
}
