package zooz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// CustomerClient is a client for work with Customer entity.
// https://developers.paymentsos.com/docs/api#/reference/customers
type CustomerClient struct {
	Caller Caller
}

// Customer is a model of entity.
type Customer struct {
	CustomerParams

	ID             string          `json:"id"`
	Created        json.Number     `json:"created"`
	Modified       json.Number     `json:"modified"`
	PaymentMethods []PaymentMethod `json:"payment_methods"`
	Href           string          `json:"href"`
}

// CustomerParams is a set of params for creating and updating entity.
type CustomerParams struct {
	CustomerReference string            `json:"customer_reference"`
	FirstName         string            `json:"first_name,omitempty"`
	LastName          string            `json:"last_name,omitempty"`
	Email             string            `json:"email,omitempty"`
	AdditionalDetails AdditionalDetails `json:"additional_details,omitempty"`
	ShippingAddress   *Address          `json:"shipping_address,omitempty"`
}

const (
	customersPath = "customers"
)

// New creates new Customer entity.
func (c *CustomerClient) New(ctx context.Context, idempotencyKey string, params *CustomerParams) (*Customer, error) {
	customer := &Customer{}
	if err := c.Caller.Call(ctx, "POST", customersPath, map[string]string{headerIdempotencyKey: idempotencyKey}, params, customer); err != nil {
		return nil, err
	}
	return customer, nil
}

// Get returns Customer entity.
func (c *CustomerClient) Get(ctx context.Context, id string) (*Customer, error) {
	customer := &Customer{}
	if err := c.Caller.Call(ctx, "GET", c.customerPath(id), nil, nil, customer); err != nil {
		return nil, err
	}
	return customer, nil
}

// GetByReference returns Customer entity by reference.
func (c *CustomerClient) GetByReference(ctx context.Context, reference string) (*Customer, error) {
	var customers []*Customer
	path := customersPath + "?customer_reference=" + url.QueryEscape(reference)
	if err := c.Caller.Call(ctx, "GET", path, nil, nil, &customers); err != nil {
		return nil, err
	}
	if len(customers) == 0 { // Should not happen. If customer is not found PaymentsOS returns 404.
		return nil, errors.New("PaymentsOS returned empty array")
	}
	if len(customers) > 1 {
		return nil, errors.New("PaymentsOS returned array with more than one item")
	}
	return customers[0], nil
}

// Update updates Customer entity with given params and return updated Customer entity.
func (c *CustomerClient) Update(ctx context.Context, id string, params *CustomerParams) (*Customer, error) {
	customer := &Customer{}
	if err := c.Caller.Call(ctx, "PUT", c.customerPath(id), nil, params, customer); err != nil {
		return nil, err
	}
	return customer, nil
}

// Delete deletes Customer entity.
func (c *CustomerClient) Delete(ctx context.Context, id string) error {
	return c.Caller.Call(ctx, "DELETE", c.customerPath(id), nil, nil, nil)
}

func (c *CustomerClient) customerPath(id string) string {
	return fmt.Sprintf("%s/%s", customersPath, id)
}
