package zooz

import (
	"context"
	"encoding/json"
	"fmt"
)

// AuthorizationClient is a client for work with Authorization entity.
// https://developers.paymentsos.com/docs/api#/reference/authorizations
type AuthorizationClient struct {
	Caller Caller
}

// Authorization is a model of entity.
type Authorization struct {
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
	ProviderConfiguration      ProviderConfiguration   `json:"provider_configuration"`
	OriginatingPurchaseCountry string                  `json:"originating_purchase_country"`
	IPAddress                  string                  `json:"ip_address"`
	Redirection                *Redirection            `json:"redirection"`
	AdditionalDetails          AdditionalDetails       `json:"additional_details"`
	DecisionEngineExecution    DecisionEngineExecution `json:"decision_engine_execution"`
}

// AuthorizationParams is a set of params for creating entity.
type AuthorizationParams struct {
	PaymentMethod            PaymentMethodDetails      `json:"payment_method"`
	MerchantSiteURL          string                    `json:"merchant_site_url,omitempty"`
	ReconciliationID         string                    `json:"reconciliation_id,omitempty"`
	ThreeDSecureAttributes   *ThreeDSecureAttributes   `json:"three_d_secure_attributes,omitempty"`
	Installments             *Installments             `json:"installments,omitempty"`
	ProviderSpecificData     map[string]interface{}    `json:"provider_specific_data,omitempty"`
	AdditionalDetails        map[string]string         `json:"additional_details,omitempty"`
	COFTransactionIndicators *COFTransactionIndicators `json:"cof_transaction_indicators,omitempty"`
}

type COFTransactionIndicators struct {
	CardEntryMode           string `json:"card_entry_mode"`
	COFConsentTransactionID string `json:"cof_consent_transaction_id"`
}

// New creates new Authorization entity.
func (c *AuthorizationClient) New(ctx context.Context, idempotencyKey string, paymentID string, params *AuthorizationParams, clientInfo *ClientInfo) (*Authorization, error) {
	authorization := &Authorization{}

	headers := map[string]string{headerIdempotencyKey: idempotencyKey}

	if clientInfo != nil {
		headers[headerClientIPAddress] = clientInfo.IPAddress
		headers[headerClientUserAgent] = clientInfo.UserAgent
	}

	if err := c.Caller.Call(ctx, "POST", c.authorizationsPath(paymentID), headers, params, authorization); err != nil {
		return nil, err
	}
	return authorization, nil
}

// Get returns Authorization entity.
func (c *AuthorizationClient) Get(ctx context.Context, paymentID string, authorizationID string) (*Authorization, error) {
	authorization := &Authorization{}
	if err := c.Caller.Call(ctx, "GET", c.authorizationPath(paymentID, authorizationID), nil, nil, authorization); err != nil {
		return nil, err
	}
	return authorization, nil
}

// GetList returns list of all Authorizations for given payment ID.
func (c *AuthorizationClient) GetList(ctx context.Context, paymentID string) ([]Authorization, error) {
	var authorizations []Authorization
	if err := c.Caller.Call(ctx, "GET", c.authorizationsPath(paymentID), nil, nil, &authorizations); err != nil {
		return nil, err
	}
	return authorizations, nil
}

func (c *AuthorizationClient) authorizationsPath(paymentID string) string {
	return fmt.Sprintf("%s/%s/authorizations", paymentsPath, paymentID)
}

func (c *AuthorizationClient) authorizationPath(paymentID string, authorizationID string) string {
	return fmt.Sprintf("%s/%s", c.authorizationsPath(paymentID), authorizationID)
}
