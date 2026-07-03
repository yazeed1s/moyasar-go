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

	invoice, err := client.Invoices.Create(context.Background(), moyasar.CreateInvoiceRequest{
		Amount:      100,
		Currency:    "SAR",
		Description: "Test invoice",
		CallbackURL: "https://example.com/invoice-callback",
		SuccessURL:  "https://example.com/success",
		BackURL:     "https://example.com/back",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("invoice id:", invoice.ID)
	fmt.Println("checkout url:", invoice.URL)
}
