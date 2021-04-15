package zooz

import (
	"context"
	"encoding/json"
	"fmt"
)

// VoidClient is a client for work with Void entity.
// https://developers.paymentsos.com/docs/api#/reference/voids
type VoidClient struct {
	Caller Caller
}

// Void is an entity model.
type Void struct {
	ID                    string                `json:"id"`
	Result                Result                `json:"result"`
	Created               json.Number           `json:"created"`
	ProviderData          ProviderData          `json:"provider_data"`
	AdditionalDetails     AdditionalDetails     `json:"additional_details"`
	ProviderConfiguration ProviderConfiguration `json:"provider_configuration"`
}

// New create new Void entity.
func (c *VoidClient) New(ctx context.Context, idempotencyKey string, paymentID string) (*Void, error) {
	void := &Void{}
	if err := c.Caller.Call(ctx, "POST", c.voidsPath(paymentID), map[string]string{headerIdempotencyKey: idempotencyKey}, nil, void); err != nil {
		return nil, err
	}
	return void, nil
}

// Get returns Void entity.
func (c *VoidClient) Get(ctx context.Context, paymentID string, voidID string) (*Void, error) {
	void := &Void{}
	if err := c.Caller.Call(ctx, "GET", c.voidPath(paymentID, voidID), nil, nil, void); err != nil {
		return nil, err
	}
	return void, nil
}

// GetList returns a list of Void for given payment.
func (c *VoidClient) GetList(ctx context.Context, paymentID string) ([]Void, error) {
	var voids []Void
	if err := c.Caller.Call(ctx, "GET", c.voidsPath(paymentID), nil, nil, &voids); err != nil {
		return nil, err
	}
	return voids, nil
}

func (c *VoidClient) voidsPath(paymentID string) string {
	return fmt.Sprintf("%s/%s/voids", paymentsPath, paymentID)
}

func (c *VoidClient) voidPath(paymentID string, voidID string) string {
	return fmt.Sprintf("%s/%s", c.voidsPath(paymentID), voidID)
}
