package zooz

import (
	"context"
	"encoding/json"
	"fmt"
)

// PaymentMethodClient is a client for work with PaymentMethod entity.
// https://developers.paymentsos.com/docs/api#/reference/payment-methods
type PaymentMethodClient struct {
	Caller Caller
}

// PaymentMethod is a entity model.
type PaymentMethod struct {
	Href               string            `json:"href"`
	Type               string            `json:"type"`
	TokenType          string            `json:"token_type"`
	PassLuhnValidation bool              `json:"pass_luhn_validation"`
	Token              string            `json:"token"`
	Created            json.Number       `json:"created"`
	Customer           string            `json:"customer"`
	AdditionalDetails  AdditionalDetails `json:"additional_details"`
	BinNumber          json.Number       `json:"bin_number"`
	Vendor             string            `json:"vendor"`
	Issuer             string            `json:"issuer"`
	CardType           string            `json:"card_type"`
	Level              string            `json:"level"`
	CountryCode        string            `json:"country_code"`
	HolderName         string            `json:"holder_name"`
	ExpirationDate     string            `json:"expiration_date"`
	Last4Digits        string            `json:"last_4_digits"`
	IdentityDocument   *IdentityDocument `json:"identity_document"`
	ShippingAddress    *Address          `json:"shipping_address"`
	BillingAddress     *Address          `json:"billing_address"`
	FingerPrint        string            `json:"fingerprint"`
}

// IdentityDocument represents some identity document.
type IdentityDocument struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

// New creates new PaymentMethod entity.
func (c *PaymentMethodClient) New(ctx context.Context, idempotencyKey string, customerID string, token string) (*PaymentMethod, error) {
	paymentMethod := &PaymentMethod{}
	if err := c.Caller.Call(ctx, "POST", c.tokenPath(customerID, token), map[string]string{headerIdempotencyKey: idempotencyKey}, nil, paymentMethod); err != nil {
		return nil, err
	}
	return paymentMethod, nil
}

// Get returns PaymentMethod entity by customer ID and token.
func (c *PaymentMethodClient) Get(ctx context.Context, customerID string, token string) (*PaymentMethod, error) {
	paymentMethod := &PaymentMethod{}
	if err := c.Caller.Call(ctx, "GET", c.tokenPath(customerID, token), nil, nil, paymentMethod); err != nil {
		return nil, err
	}
	return paymentMethod, nil
}

// GetList returns list of PaymentMethods for given customer.
func (c *PaymentMethodClient) GetList(ctx context.Context, customerID string) ([]PaymentMethod, error) {
	var paymentMethods []PaymentMethod
	if err := c.Caller.Call(ctx, "GET", c.paymentMethodsPath(customerID), nil, nil, &paymentMethods); err != nil {
		return nil, err
	}
	return paymentMethods, nil
}

// Delete customer PaymentMethod by token.
func (c *PaymentMethodClient) Delete(ctx context.Context, customerID string, token string) error {
	return c.Caller.Call(ctx, "DELETE", c.tokenPath(customerID, token), nil, nil, nil)
}

func (c *PaymentMethodClient) paymentMethodsPath(customerID string) string {
	return fmt.Sprintf("%s/%s/payment-methods", customersPath, customerID)
}

func (c *PaymentMethodClient) tokenPath(customerID, token string) string {
	return fmt.Sprintf("%s/%s", c.paymentMethodsPath(customerID), token)
}
