package zooz

import (
	"context"
	"fmt"
)

type TokenType string

const (
	TokenTypeCreditCard       TokenType = "credit_card"
	TokenTypeCardCVVCode      TokenType = "card_cvv_code"
	TokenTypeBillingAgreement TokenType = "billing_agreement"
)

// CreditCardTokenClient is a client for work with Token entity (type = credit_card).
// https://developers.paymentsos.com/docs/api#/reference/tokens
type CreditCardTokenClient struct {
	Caller Caller
}

// CreditCardToken is a model of entity.
type CreditCardToken struct {
	TokenType TokenType `json:"token_type"`

	// State - Enum: "valid" "created" "used" "assigned".
	// Reflects the token's usage. The state can be one of the following:
	// 	* valid: token has been successfully created and can be used. Only applies to tokens created in API version 1.1.0 or below.
	// 	* created: token has been successfully created but has not yet been used or assigned to a customer.
	// 	* used: token has been used in a successful authorize or charge request.
	// 	* assigned: the token is assigned to a customer. The token will remain in this state also if it is used in a new authorize or charge request.
	State string `json:"state"`

	// PassLuhnValidation - When token_type is credit_card, then this indicates if the credit card number passed the Luhn validation.
	PassLuhnValidation bool `json:"pass_luhn_validation"`

	// BinNumber - The initial four to six numbers that appear on the card.
	BinNumber string `json:"bin_number"`

	// Vendor - The name of the credit card corporation.
	Vendor string `json:"vendor"`

	// Issuer - The name of the bank that issued the card.
	Issuer string `json:"issuer"` // TODO: WTF? in the API docs it is issuer_name, but actually it is issuer

	// CardType - The type of card.
	CardType string `json:"card_type"`

	// Level - The level of benefits or services available with the card.
	Level string `json:"level"`

	// CountryCode - The 3-letter country code defined in ISO 3166-1 alpha-3 format, identifying the country in which the card was issued.
	CountryCode string `json:"country_code"`

	// HolderName - Name of the credit card holder.
	HolderName string `json:"holder_name"`

	// Credit card expiration date.
	ExpirationDate ExpirationDate `json:"expiration_date"`

	// Token - Depending on the `token_type`, the token represents one of the following:
	// a customer's credit card number, the card cvv code or a billing agreement.
	Token string `json:"token"`

	// Created - <timestamp> The date and time that the token was created.
	Created string `json:"created"`

	// Type - Value: "tokenized". Depending on the `token_type`, this field represents either
	// the card's or billing agreement's representation as it will be used in an authorization or charge request.
	Type string `json:"type"`

	// EncryptedCVV - <missing in API docs>
	// Returned only for new token creation. And only if `credit_card_cvv` was sent.
	// Expires after three hours.
	EncryptedCVV string `json:"encrypted_cvv"`
}

// CreditCardTokenParams is a set of params for creating entity.
type CreditCardTokenParams struct {
	// HolderName - Name of the credit card holder.
	HolderName string `json:"holder_name"`

	// ExpirationDate - ^(0[1-9]|1[0-2])(\/|\-|\.| )\d{2,4} Credit card expiration date.
	// Possible formats: mm-yyyy, mm-yy, mm.yyyy, mm.yy, mm/yy, mm/yyyy, mm yyyy, or mm yy.
	ExpirationDate string `json:"expiration_date,omitempty"`

	// IdentityDocument - National identity document of the card holder.
	IdentityDocument *IdentityDocument `json:"identity_document,omitempty"`

	// CardNumber - \d{8}|\d{12,19} Credit card number.
	CardNumber string `json:"card_number"`

	ShippingAddress   *Address          `json:"shipping_address,omitempty"`
	BillingAddress    *Address          `json:"billing_address,omitempty"`
	AdditionalDetails AdditionalDetails `json:"additional_details,omitempty"`

	// CreditCardCVV - The CVV number on the card (3 or 4 digits) to be encrypted.
	// Sending this field returns an encrypted_cvv field, which expires after three hours.
	CreditCardCVV string `json:"credit_card_cvv,omitempty"`
}

type privateCreditCardTokenParams struct {
	TokenType TokenType `json:"token_type"`
	*CreditCardTokenParams
}

const (
	tokensPath = "tokens"
)

// New creates new Token entity (type = credit_card).
func (c *CreditCardTokenClient) New(ctx context.Context, idempotencyKey string, params *CreditCardTokenParams) (*CreditCardToken, error) {
	token := &CreditCardToken{}
	err := c.Caller.Call(
		ctx,
		"POST",
		tokensPath,
		map[string]string{headerIdempotencyKey: idempotencyKey},
		privateCreditCardTokenParams{
			TokenType:             TokenTypeCreditCard,
			CreditCardTokenParams: params,
		},
		token,
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Get returns Token entity.
func (c *CreditCardTokenClient) Get(ctx context.Context, token string) (*CreditCardToken, error) {
	creditCardToken := &CreditCardToken{}
	if err := c.Caller.Call(ctx, "GET", c.tokenPath(token), nil, nil, creditCardToken); err != nil {
		return nil, err
	}
	return creditCardToken, nil
}

func (c *CreditCardTokenClient) tokenPath(token string) string {
	return fmt.Sprintf("%s/%s", tokensPath, token)
}
