package zooz

import (
	"context"
	"encoding/json"
	"fmt"
)

// RefundClient is a client for work with Refund entity.
// https://developers.paymentsos.com/docs/api#/reference/refunds
type RefundClient struct {
	Caller Caller
}

// Refund is a entity model.
type Refund struct {
	RefundParams

	ID                    string                `json:"id"`
	Result                Result                `json:"result"`
	Created               json.Number           `json:"created"`
	ProviderData          ProviderData          `json:"provider_data"`
	AdditionalDetails     AdditionalDetails     `json:"additional_details"`
	ProviderConfiguration ProviderConfiguration `json:"provider_configuration"`
}

// RefundParams is a set of params for creating entity.
type RefundParams struct {
	ReconciliationID string `json:"reconciliation_id,omitempty"`
	Amount           int64  `json:"amount,omitempty"`
	CaptureID        string `json:"capture_id,omitempty"`
	Reason           string `json:"reason,omitempty"`
}

// New creates new Refund entity.
func (c *RefundClient) New(ctx context.Context, idempotencyKey string, paymentID string, params *RefundParams) (*Refund, error) {
	refund := &Refund{}
	if err := c.Caller.Call(ctx, "POST", c.refundsPath(paymentID), map[string]string{headerIdempotencyKey: idempotencyKey}, params, refund); err != nil {
		return nil, err
	}
	return refund, nil
}

// Get returns Refund entity.
func (c *RefundClient) Get(ctx context.Context, paymentID string, refundID string) (*Refund, error) {
	refund := &Refund{}
	if err := c.Caller.Call(ctx, "GET", c.refundPath(paymentID, refundID), nil, nil, refund); err != nil {
		return nil, err
	}
	return refund, nil
}

// GetList returns a list of Refunds for given payment.
func (c *RefundClient) GetList(ctx context.Context, paymentID string) ([]Refund, error) {
	var refunds []Refund
	if err := c.Caller.Call(ctx, "GET", c.refundsPath(paymentID), nil, nil, &refunds); err != nil {
		return nil, err
	}
	return refunds, nil
}

func (c *RefundClient) refundsPath(paymentID string) string {
	return fmt.Sprintf("%s/%s/refunds", paymentsPath, paymentID)
}

func (c *RefundClient) refundPath(paymentID string, refundID string) string {
	return fmt.Sprintf("%s/%s", c.refundsPath(paymentID), refundID)
}
