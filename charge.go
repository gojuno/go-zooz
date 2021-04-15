package zooz

import (
	"context"
	"encoding/json"
	"fmt"
)

// ChargeClient is a client for work with Charge entity.
// https://developers.paymentsos.com/docs/api#/reference/charges
type ChargeClient struct {
	Caller Caller
}

// Charge is a model of entity.
type Charge struct {
	ID                         string                  `json:"id"`
	Result                     Result                  `json:"result"`
	Amount                     int64                   `json:"amount"`
	Created                    json.Number             `json:"created"`
	ReconciliationID           string                  `json:"reconciliation_id"`
	PaymentMethod              PaymentMethod           `json:"payment_method"`
	ThreeDSecureAttributes     *ThreeDSecureAttributes `json:"three_d_secure_attributes"`
	Installments               *Installments           `json:"installments"`
	ProviderData               ProviderData            `json:"provider_data"`
	ProviderSpecificData       DecodedJSON             `json:"provider_specific_data"`
	OriginatingPurchaseCountry string                  `json:"originating_purchase_country"`
	IPAddress                  string                  `json:"ip_address"`
	Redirection                *Redirection            `json:"redirection"`
	ProviderConfiguration      ProviderConfiguration   `json:"provider_configuration"`
	AdditionalDetails          AdditionalDetails       `json:"additional_details"`
	DecisionEngineExecution    DecisionEngineExecution `json:"decision_engine_execution"`
}

// ChargeParams is a set of params for creating entity.
type ChargeParams struct {
	PaymentMethod          PaymentMethodDetails    `json:"payment_method"`
	MerchantSiteURL        string                  `json:"merchant_site_url,omitempty"`
	ReconciliationID       string                  `json:"reconciliation_id,omitempty"`
	ThreeDSecureAttributes *ThreeDSecureAttributes `json:"three_d_secure_attributes,omitempty"`
	Installments           *Installments           `json:"installments,omitempty"`
	ProviderSpecificData   map[string]interface{}  `json:"provider_specific_data,omitempty"`
}

// New creates new Charge entity.
func (c *ChargeClient) New(ctx context.Context, idempotencyKey string, paymentID string, params *ChargeParams, clientInfo *ClientInfo) (*Charge, error) {
	charge := &Charge{}

	headers := map[string]string{headerIdempotencyKey: idempotencyKey}

	if clientInfo != nil {
		headers[headerClientIPAddress] = clientInfo.IPAddress
		headers[headerClientUserAgent] = clientInfo.UserAgent
	}

	if err := c.Caller.Call(ctx, "POST", c.chargesPath(paymentID), headers, params, charge); err != nil {
		return nil, err
	}
	return charge, nil
}

// Get returns Charge entity.
func (c *ChargeClient) Get(ctx context.Context, paymentID string, chargeID string) (*Charge, error) {
	charge := &Charge{}
	if err := c.Caller.Call(ctx, "GET", c.chargePath(paymentID, chargeID), nil, nil, charge); err != nil {
		return nil, err
	}
	return charge, nil
}

// GetList returns a list of Charges for given payment ID.
func (c *ChargeClient) GetList(ctx context.Context, paymentID string) ([]Charge, error) {
	var charges []Charge
	if err := c.Caller.Call(ctx, "GET", c.chargesPath(paymentID), nil, nil, &charges); err != nil {
		return nil, err
	}
	return charges, nil
}

func (c *ChargeClient) chargesPath(paymentID string) string {
	return fmt.Sprintf("%s/%s/charges", paymentsPath, paymentID)
}

func (c *ChargeClient) chargePath(paymentID string, chargeID string) string {
	return fmt.Sprintf("%s/%s", c.chargesPath(paymentID), chargeID)
}
