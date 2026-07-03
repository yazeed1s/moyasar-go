# moyasar-go

Small Go SDK for Moyasar.

It wraps the Moyasar REST API with typed services and request structs.

## Install

```sh
go get github.com/yazeed1s/moyasar-go
```

## Create a client

```go
client := moyasar.NewClient("sk_test_xxx")
```

The SDK uses Basic Auth like Moyasar docs say. The API key is the username, and
the password is empty.

You can also use your own HTTP client:

```go
client := moyasar.NewClient(
	"sk_test_xxx",
	moyasar.WithHTTPClient(http.DefaultClient),
)
```

## Create a payment

```go
payment, err := client.Payments.Create(ctx, moyasar.CreatePaymentRequest{
	GivenID:     "a1168bd1-47a4-4b97-8a50-dd5caaccacf2",
	Amount:      100,
	Currency:    "SAR",
	Description: "Order #123",
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
	return err
}
```

`GivenID` is used for Moyasar payment idempotency. Use a UUID from your app.

## List payments

```go
payments, err := client.Payments.List(ctx, moyasar.ListPaymentsParams{
	Page:   1,
	Status: moyasar.PaymentStatusPaid,
})
```

List responses include `Meta` for pagination.

## Services

The client has these services:

- `Payments`
- `Invoices`
- `Tokens`
- `Sources`
- `CardAuths`
- `Webhooks`
- `Payouts`
- `Settlements`
- `InternalTransactions`
- `Transfers`

Transfers use the Moyasar aggregation base URL. You can change it with:

```go
client := moyasar.NewClient(
	"sk_test_xxx",
	moyasar.WithTransferURL("https://apimig.moyasar.com/v1"),
)
```

## Errors

When Moyasar returns a non-2xx response, the SDK returns `*moyasar.APIError`.

```go
var apiErr *moyasar.APIError
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.StatusCode)
	fmt.Println(apiErr.Type)
	fmt.Println(apiErr.Message)
}
```

## Examples

See the `examples/` folder.

Run one example:

```sh
MOYASAR_API_KEY=sk_test_xxx go run ./examples/list-payments
```

Use test keys while trying the SDK.
