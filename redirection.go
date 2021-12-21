package zooz

import (
	"context"
	"encoding/json"
	"fmt"
)

// RedirectionClient is a client for work with Redirection entity.
// https://developers.paymentsos.com/docs/api#/reference/redirections
type RedirectionClient struct {
	Caller Caller
}

// Redirection is a entity model.
type Redirection struct {
	ID              string      `json:"id"`
	Created         json.Number `json:"created"`
	MerchantSiteURL string      `json:"merchant_site_url"`
	URL             string      `json:"url"`
	OperationType   string      `json:"operation_type"`
}

// Get creates new Redirection entity.
func (c *RedirectionClient) Get(ctx context.Context, paymentID string, redirectionID string) (*Redirection, error) {
	redirection := &Redirection{}
	if err := c.Caller.Call(ctx, "GET", c.redirectionPath(paymentID, redirectionID), nil, nil, redirection); err != nil {
		return nil, err
	}
	return redirection, nil
}

// GetList returns a list of Redirections for given payment.
func (c *RedirectionClient) GetList(ctx context.Context, paymentID string) ([]Redirection, error) {
	var redirections []Redirection
	if err := c.Caller.Call(ctx, "GET", c.redirectionsPath(paymentID), nil, nil, &redirections); err != nil {
		return nil, err
	}
	return redirections, nil
}

func (c *RedirectionClient) redirectionsPath(paymentID string) string {
	return fmt.Sprintf("%s/%s/redirections", paymentsPath, paymentID)
}

func (c *RedirectionClient) redirectionPath(paymentID string, redirectionID string) string {
	return fmt.Sprintf("%s/%s", c.redirectionsPath(paymentID), redirectionID)
}
