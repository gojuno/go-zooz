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
	ID                         string                 `json:"id"`
	Result                     Result                 `json:"result"`
	Amount                     float64                `json:"amount"`
	Created                    json.Number            `json:"created"`
	ReconciliationID           string                 `json:"reconciliation_id"`
	PaymentMethod              PaymentMethod          `json:"payment_method"`
	ThreeDSecureAttributes     ThreeDSecureAttributes `json:"three_d_secure_attributes"`
	Installments               *Installments          `json:"installments"`
	ProviderData               ProviderData           `json:"provider_data"`
	ProviderSpecificData       map[string]interface{} `json:"provider_specific_data"`
	OriginatingPurchaseCountry string                 `json:"originating_purchase_country"`
	IpAddress                  string                 `json:"ip_address"`
	Redirection                *Redirection           `json:"redirection"`
}

// AuthorizationParams is a set of params for creating entity.
type AuthorizationParams struct {
	PaymentMethodToken     string                 `json:"payment_method_token"`
	CreditCardCvv          string                 `json:"credit_card_cvv"`
	MerchantSiteUrl        string                 `json:"merchant_site_url"`
	ReconciliationID       string                 `json:"reconciliation_id"`
	ThreeDSecureAttributes ThreeDSecureAttributes `json:"three_d_secure_attributes"`
	Installments           *Installments          `json:"installments"`
	ProviderSpecificData   map[string]interface{} `json:"provider_specific_data"`
}

// New creates new Authorization entity.
func (c *AuthorizationClient) New(ctx context.Context, idempotencyKey string, paymentID string, params *AuthorizationParams, clientInfo *ClientInfo) (*Authorization, error) {
	authorization := &Authorization{}

	headers := map[string]string{headerIdempotencyKey: idempotencyKey}

	if clientInfo != nil {
		headers[headerClientIpAddress] = clientInfo.IpAddress
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
