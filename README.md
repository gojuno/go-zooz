# Zooz API client [![GoDoc](https://godoc.org/github.com/gojuno/go-zooz?status.svg)](http://godoc.org/github.com/gojuno/go-zooz) [![Build Status](https://travis-ci.org/gojuno/go-zooz.svg?branch=master)](https://travis-ci.org/gojuno/go-zooz) [![Go Report Card](https://goreportcard.com/badge/github.com/gojuno/go-zooz)](https://goreportcard.com/report/github.com/gojuno/go-zooz)

This repo contains Zooz API client written in Go.

Zooz API documentation: https://developers.paymentsos.com/docs/api

Before using this client you need to register and configure Zooz account: https://developers.paymentsos.com/docs/quick-start.html

## How to install

Download package:
```
go get github.com/gojuno/go-zooz
```

Client uses `github.com/pkg/errors`, so you may need to download this package as well:
```
go get github.com/pkg/errors
```

## How to use

To init client you will need `private_key` and `app_id` which you can get from your Zooz account profile.
```
import "github.com/gojuno/go-zooz"
...
// Init client
client := zooz.New(
	zooz.OptAppID("com.yourhost.go_client"),
	zooz.OptPrivateKey("a630518c-22da-4eaa-bb39-502ad7832030"),
)

// Create new customer
customer, customerErr := client.Customer().New(
	context.Background(),
	"customer_idempotency_key",
	&zooz.CustomerParams{
		CustomerReference: "1234",
		FirstName:         "John",
		LastName:          "Doe",
	},
)

// Create new payment method
paymentMethod, paymentMethodErr := client.PaymentMethod().New(
	context.Background(),
	"payment_method_idempotency_key",
	customer.ID,
	"918a917e-4cf9-4303-949c-d0cd7ff7f619",
)

// Delete customer
deleteCustomerErr := client.Customer().Delete(context.Background(), customer.ID)
```

## Custom HTTP client

By default Zooz client uses `http.DefaultClient`. You can set custom HTTP client using `zooz.OptHttpClient` option:
```
httpClient := &http.Client{
	Timeout: time.Minute,
}

client := zooz.New(
	zooz.OptAppID("com.yourhost.go_client"),
	zooz.OptPrivateKey("a630518c-22da-4eaa-bb39-502ad7832030"),
	zooz.OptHttpClient(httpClient),
)
```
You can use any HTTP client, implementing `zooz.httpClient` interface with method `Do(r *http.Request) (*http.Response, error)`. Built-in `net/http` client implements it, of course.

## Test/live environment

Zooz supports test and live environment. Environment is defined by `x-payments-os-env` request header.

By default, client sends `test` value. You can redefine this value to `live` using `zooz.OptEnv(zooz.EnvLive)` option.
```
client := zooz.New(
	zooz.OptAppID("com.yourhost.go_client"),
	zooz.OptPrivateKey("a630518c-22da-4eaa-bb39-502ad7832030"),
	zooz.OptEnv(zooz.EnvLive),
)
```

## Tokens

API methods for Tokens are not implemented in this client, because they are supposed to be used on client-side, not server-side. See example here: https://developers.paymentsos.com/docs/collecting-payment-details.html