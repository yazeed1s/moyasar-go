package moyasar

// Metadata contains merchant-defined key/value data attached to Moyasar objects.
//
// Moyasar supports metadata on payments, invoices, payouts, and tokens. Do not
// store sensitive data such as card or bank account details in metadata.
type Metadata map[string]string
