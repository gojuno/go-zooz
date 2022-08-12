package zooz

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreditCardTokenClient_New(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "POST",
		expectedURL:    "/tokens",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		expectedBodyJSON: /* language=json */ `{
			"token_type": "credit_card",
			"holder_name": "holder name",
			"expiration_date": "12-2051",
			"identity_document": {
				"type": "identity type",
				"number": "identity number"
			},
			"card_number": "378282246310005",
			"shipping_address": {
				"country": "RUS",
				"state": "shipping state",
				"city": "shipping city",
				"line1": "shipping line1",
				"line2": "shipping line2",
				"zip_code": "shipping zip code",
				"title": "shipping title",
				"first_name": "shipping first name",
				"last_name": "shipping last name",
				"phone": "shipping phone",
				"email": "shipping-address@email.com"
			},
			"billing_address": {
				"country": "RUS",
				"state": "billing state",
				"city": "billing city",
				"line1": "billing line1",
				"line2": "billing line2",
				"zip_code": "billing zip code",
				"title": "billing title",
				"first_name": "billing first name",
				"last_name": "billing last name",
				"phone": "billing phone",
				"email": "billing-address@email.com"
			},
			"additional_details": {
				"token detail 1": "value 1",
				"token detail 2": "value 2"
			},
			"credit_card_cvv": "123"
		}`,
		responseBody: /* language=json */ ` {
			"token": "78565e54-9439-4cbf-91e0-fc7fc33703b6",
			"created": "1630440556137",
			"pass_luhn_validation": true,
			"encrypted_cvv": "xxx",
			"token_type": "credit_card",
			"type": "tokenized",
			"state": "created",
			"bin_number": "378282",
			"vendor": "AMERICAN EXPRESS",
			"card_type": "CREDIT",
			"issuer": "AMERICAN EXPRESS US (CARS)",
			"level": "CORPORATE",
			"country_code": "USA",
			"holder_name": "John Doe",
			"expiration_date": "12/2051",
			"last_4_digits": "0005"
		}`,
	}

	c := &CreditCardTokenClient{Caller: New(OptHTTPClient(cli))}

	creditCardTokenParams := CreditCardTokenParams{
		HolderName:     "holder name",
		ExpirationDate: "12-2051",
		IdentityDocument: &IdentityDocument{
			Type:   "identity type",
			Number: "identity number",
		},
		CardNumber: "378282246310005",
		ShippingAddress: &Address{
			Country:   "RUS",
			State:     "shipping state",
			City:      "shipping city",
			Line1:     "shipping line1",
			Line2:     "shipping line2",
			ZipCode:   "shipping zip code",
			Title:     "shipping title",
			FirstName: "shipping first name",
			LastName:  "shipping last name",
			Phone:     "shipping phone",
			Email:     "shipping-address@email.com",
		},
		BillingAddress: &Address{
			Country:   "RUS",
			State:     "billing state",
			City:      "billing city",
			Line1:     "billing line1",
			Line2:     "billing line2",
			ZipCode:   "billing zip code",
			Title:     "billing title",
			FirstName: "billing first name",
			LastName:  "billing last name",
			Phone:     "billing phone",
			Email:     "billing-address@email.com",
		},
		AdditionalDetails: AdditionalDetails{
			"token detail 1": "value 1",
			"token detail 2": "value 2",
		},
		CreditCardCVV: "123",
	}

	token, err := c.New(
		context.Background(),
		"idempotency_key",
		&creditCardTokenParams,
	)

	require.NoError(t, err)
	require.Equal(t, &CreditCardToken{
		TokenType:          "credit_card",
		State:              "created",
		PassLuhnValidation: true,
		BinNumber:          "378282",
		Vendor:             "AMERICAN EXPRESS",
		Issuer:             "AMERICAN EXPRESS US (CARS)",
		CardType:           "CREDIT",
		Level:              "CORPORATE",
		CountryCode:        "USA",
		HolderName:         "John Doe",
		ExpirationDate:     "12/2051",
		Last4Digits:        "0005",
		Token:              "78565e54-9439-4cbf-91e0-fc7fc33703b6",
		Created:            "1630440556137",
		Type:               "tokenized",
		EncryptedCVV:       "xxx",
	}, token)
}

func TestCreditCardTokenClient_Get(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/tokens/78565e54-9439-4cbf-91e0-fc7fc33703b6",
		responseBody: /* language=json */ `{
			"token": "78565e54-9439-4cbf-91e0-fc7fc33703b6",
			"created": "1630440556137",
			"pass_luhn_validation": true,
			"token_type": "credit_card",
			"type": "tokenized",
			"state": "created",
			"bin_number": "378282",
			"vendor": "AMERICAN EXPRESS",
			"card_type": "CREDIT",
			"issuer": "AMERICAN EXPRESS US (CARS)",
			"level": "CORPORATE",
			"country_code": "USA",
			"holder_name": "John Doe",
			"expiration_date": "12/2051",
			"last_4_digits": "0005"
		}`,
	}

	c := &CreditCardTokenClient{Caller: New(OptHTTPClient(cli))}

	token, err := c.Get(
		context.Background(),
		"78565e54-9439-4cbf-91e0-fc7fc33703b6",
	)

	require.NoError(t, err)
	require.Equal(t, &CreditCardToken{
		TokenType:          "credit_card",
		State:              "created",
		PassLuhnValidation: true,
		BinNumber:          "378282",
		Vendor:             "AMERICAN EXPRESS",
		Issuer:             "AMERICAN EXPRESS US (CARS)",
		CardType:           "CREDIT",
		Level:              "CORPORATE",
		CountryCode:        "USA",
		HolderName:         "John Doe",
		ExpirationDate:     "12/2051",
		Last4Digits:        "0005",
		Token:              "78565e54-9439-4cbf-91e0-fc7fc33703b6",
		Created:            "1630440556137",
		Type:               "tokenized",
		EncryptedCVV:       "",
	}, token)
}
