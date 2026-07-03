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
		GivenID:     "a1168bd1-47a4-4b97-8a50-dd5caaccacf2",
		Amount:      100,
		Currency:    "SAR",
		Description: "Test order",
		CallbackURL: "https://example.com/callback",
		Source: moyasar.CreditCardSource{
			Name:   "John Doe",
			Number: "4111111111111111",
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
}
