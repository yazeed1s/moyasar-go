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

	list, err := client.Payments.List(context.Background(), moyasar.ListPaymentsParams{
		Page: 1,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, payment := range list.Payments {
		fmt.Println(payment.ID, payment.Status, payment.Amount, payment.Currency)
	}

	fmt.Println("page:", list.Meta.CurrentPage)
}
