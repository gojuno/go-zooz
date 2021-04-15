package zooz

import (
	"context"
	"encoding/json"
	"fmt"
)

// CaptureClient is a client for work with Capture entity.
// https://developers.paymentsos.com/docs/api#/reference/captures
type CaptureClient struct {
	Caller Caller
}

// Capture is a model of entity.
type Capture struct {
	CaptureParams

	ID                    string                `json:"id"`
	Result                Result                `json:"result"`
	Created               json.Number           `json:"created"`
	ProviderData          ProviderData          `json:"provider_data"`
	ProviderSpecificData  DecodedJSON           `json:"provider_specific_data"`
	Level23               Level23               `json:"level_2_3"`
	ProviderConfiguration ProviderConfiguration `json:"provider_configuration"`
	AdditionalDetails     AdditionalDetails     `json:"additional_details"`
}

// CaptureParams is a set of params for creating entity.
type CaptureParams struct {
	ReconciliationID string `json:"reconciliation_id,omitempty"`
	Amount           int64  `json:"amount,omitempty"`
}

// New creates new Capture entity.
func (c *CaptureClient) New(ctx context.Context, idempotencyKey string, paymentID string, params *CaptureParams) (*Capture, error) {
	capture := &Capture{}
	if err := c.Caller.Call(ctx, "POST", c.capturesPath(paymentID), map[string]string{headerIdempotencyKey: idempotencyKey}, params, capture); err != nil {
		return nil, err
	}
	return capture, nil
}

// Get returns Capture entity.
func (c *CaptureClient) Get(ctx context.Context, paymentID string, captureID string) (*Capture, error) {
	capture := &Capture{}
	if err := c.Caller.Call(ctx, "GET", c.capturePath(paymentID, captureID), nil, nil, capture); err != nil {
		return nil, err
	}
	return capture, nil
}

// GetList returns list of Captures for given payment ID.
func (c *CaptureClient) GetList(ctx context.Context, paymentID string) ([]Capture, error) {
	var captures []Capture
	if err := c.Caller.Call(ctx, "GET", c.capturesPath(paymentID), nil, nil, &captures); err != nil {
		return nil, err
	}
	return captures, nil
}

func (c *CaptureClient) capturesPath(paymentID string) string {
	return fmt.Sprintf("%s/%s/captures", paymentsPath, paymentID)
}

func (c *CaptureClient) capturePath(paymentID string, captureID string) string {
	return fmt.Sprintf("%s/%s", c.capturesPath(paymentID), captureID)
}
