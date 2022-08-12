package sandbox

import (
	"context"
	"net/http"
	"testing"

	"github.com/gtforge/go-zooz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				HolderName:         "holder name",
				ExpirationDate:     tokenCreated.ExpirationDate,
				Last4Digits:        "0005",
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
			ExpirationDate:    "", // looks like it is not required
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
				HolderName:         "holder name",
				ExpirationDate:     "",
				Last4Digits:        "1881",
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

	t.Run("idempotency works only for the same card", func(t *testing.T) {
		t.Parallel()

		idempotencyKey1 := randomString(32)
		idempotencyKey2 := randomString(32) // different key -> new token
		const name1, card1 = "name1", "4012888888881881"
		const name2, card2 = "name2", "6011601160116611" // different card -> new token

		token1, err := client.CreditCardToken().New(context.Background(), idempotencyKey1, &zooz.CreditCardTokenParams{
			HolderName: name1,
			CardNumber: card1,
		})
		require.NoError(t, err)

		token2, err := client.CreditCardToken().New(context.Background(), idempotencyKey1, &zooz.CreditCardTokenParams{
			HolderName: name1,
			CardNumber: card1,
		})
		require.NoError(t, err)
		require.Equal(t, token1, token2)

		token3, err := client.CreditCardToken().New(context.Background(), idempotencyKey1, &zooz.CreditCardTokenParams{
			HolderName: name2,
			CardNumber: card2, // different card
		})
		require.NoError(t, err)
		require.NotEqual(t, token1.Token, token3.Token)

		token4, err := client.CreditCardToken().New(context.Background(), idempotencyKey2, &zooz.CreditCardTokenParams{ // different key
			HolderName: name1,
			CardNumber: card1,
		})
		require.NoError(t, err)
		require.NotEqual(t, token1.Token, token4.Token)
	})

	t.Run("get - unknown token", func(t *testing.T) {
		t.Parallel()

		_, err := client.CreditCardToken().Get(context.Background(), UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "Not Found, Payment method token was not found",
		})
	})
}

// PrepareToken is a helper to create new token in zooz.
func PrepareToken(t *testing.T, client *zooz.Client) (*zooz.CreditCardToken, *zooz.CreditCardTokenParams) {
	tokenParams := &zooz.CreditCardTokenParams{
		HolderName:     "holder name",
		ExpirationDate: "12-51",
		IdentityDocument: &zooz.IdentityDocument{
			Type:   "identity type",
			Number: "identity number",
		},
		CardNumber: "4012888888881881",
		ShippingAddress: &zooz.Address{
			Country:   "RUS",
			State:     "token shipping state",
			City:      "token shipping city",
			Line1:     "token shipping line1",
			Line2:     "token shipping line2",
			ZipCode:   "token shipping zip code",
			Title:     "token shipping title",
			FirstName: "token shipping first name",
			LastName:  "token shipping last name",
			Phone:     "token shipping phone",
			Email:     "token-shipping-address@email.com",
		},
		BillingAddress: &zooz.Address{
			Country:   "RUS",
			State:     "token billing state",
			City:      "token billing city",
			Line1:     "token billing line1",
			Line2:     "token billing line2",
			ZipCode:   "token billing zip code",
			Title:     "token billing title",
			FirstName: "token billing first name",
			LastName:  "token billing last name",
			Phone:     "token billing phone",
			Email:     "token-billing-address@email.com",
		},
		AdditionalDetails: zooz.AdditionalDetails{
			"token detail 1": "token value 1",
			"token detail 2": "token value 2",
		},
		CreditCardCVV: "123",
	}
	token, err := client.CreditCardToken().New(context.Background(), randomString(32), tokenParams)
	require.NoError(t, err)
	return token, tokenParams
}
