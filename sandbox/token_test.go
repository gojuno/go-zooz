package sandbox

import (
	"context"
	"github.com/gtforge/go-zooz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestToken(t *testing.T) {
	t.Parallel()
	client := GetClient(t)

	t.Run("new (all fields) & get", func(t *testing.T) {
		t.Parallel()

		tokenCreated, err := client.CreditCardToken().New(context.Background(), randomString(32), &zooz.CreditCardTokenParams{
			HolderName:     "holder name",
			ExpirationDate: "12-2051",
			IdentityDocument: &zooz.IdentityDocument{
				Type:   "identity type",
				Number: "identity number",
			},
			CardNumber: "378282246310005",
			ShippingAddress: &zooz.Address{
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
			BillingAddress: &zooz.Address{
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
			AdditionalDetails: zooz.AdditionalDetails{
				"token detail 1": "value 1",
				"token detail 2": "value 2",
			},
			CreditCardCVV: "123",
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, tokenCreated.Token)
			assert.NotEmpty(t, tokenCreated.Created)
			assert.NotEmpty(t, tokenCreated.EncryptedCVV)
			assert.Equal(t, &zooz.CreditCardToken{
				TokenType:          zooz.TokenTypeCreditCard,
				State:              "created",
				PassLuhnValidation: true,
				BinNumber:          "378282",
				Vendor:             "AMERICAN EXPRESS",
				Issuer:             "AMERICAN EXPRESS US (CARS)",
				CardType:           "CREDIT",
				Level:              "CORPORATE",
				CountryCode:        "USA",
				Token:              tokenCreated.Token,
				Created:            tokenCreated.Created,
				Type:               "tokenized",
				EncryptedCVV:       tokenCreated.EncryptedCVV,
			}, tokenCreated)
		})

		tokenRetrieved, err := client.CreditCardToken().Get(context.Background(), tokenCreated.Token)
		require.NoError(t, err)
		must(t, func() {
			assert.Empty(t, tokenRetrieved.EncryptedCVV) // encrypted CVV returned only once
			tokenRetrieved.EncryptedCVV = tokenCreated.EncryptedCVV
			assert.Equal(t, tokenCreated, tokenRetrieved)
		})
	})

	t.Run("new (required fields) & get", func(t *testing.T) {
		t.Parallel()

		tokenCreated, err := client.CreditCardToken().New(context.Background(), randomString(32), &zooz.CreditCardTokenParams{
			HolderName:        "holder name",
			ExpirationDate:    "", // cool, it is not required
			IdentityDocument:  nil,
			CardNumber:        "4012888888881881",
			ShippingAddress:   nil,
			BillingAddress:    nil,
			AdditionalDetails: nil,
			CreditCardCVV:     "",
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, tokenCreated.Token)
			assert.NotEmpty(t, tokenCreated.Created)
			assert.Equal(t, &zooz.CreditCardToken{
				TokenType:          zooz.TokenTypeCreditCard,
				State:              "created",
				PassLuhnValidation: true,
				BinNumber:          "401288",
				Vendor:             "VISA",
				Issuer:             "PJSC CREDIT AGRICOLE BANK",
				CardType:           "CREDIT",
				Level:              "INFINITE",
				CountryCode:        "UKR",
				Token:              tokenCreated.Token,
				Created:            tokenCreated.Created,
				Type:               "tokenized",
				EncryptedCVV:       "",
			}, tokenCreated)
		})

		tokenRetrieved, err := client.CreditCardToken().Get(context.Background(), tokenCreated.Token)
		require.NoError(t, err)
		require.Equal(t, tokenCreated, tokenRetrieved)
	})

	t.Run("get - unknown token", func(t *testing.T) {
		t.Parallel()

		_, err := client.CreditCardToken().Get(context.Background(), "00000000-0000-1000-8000-000000000000")
		zoozErr := &zooz.Error{}
		require.ErrorAs(t, err, &zoozErr)
		require.Equal(t, &zooz.Error{
			StatusCode: http.StatusNotFound,
			RequestID:  zoozErr.RequestID, // ignore
			APIError: zooz.APIError{
				Category:    "api_request_error",
				Description: "The resource was not found.",
				MoreInfo:    "Not Found, Payment method token was not found",
			},
		}, zoozErr)
	})
}
