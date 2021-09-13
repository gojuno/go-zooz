package sandbox

import (
	"context"
	"encoding/json"
	"github.com/gtforge/go-zooz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestPaymentMethod(t *testing.T) {
	t.Parallel()
	client := GetClient(t)

	t.Run("new & get & get-list", func(t *testing.T) {
		t.Parallel()

		customer := PrepareCustomer(t, client)
		token, tokenParams := PrepareToken(t, client)

		paymentMethodCreated, err := client.PaymentMethod().New(context.Background(), randomString(32), customer.ID, token.Token)
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, paymentMethodCreated.FingerPrint)
			assert.Equal(t, &zooz.PaymentMethod{
				Href:               zooz.ApiURL + "/customers/" + customer.ID + "/payment-methods/" + token.Token,
				Type:               token.Type,
				TokenType:          string(token.TokenType),
				PassLuhnValidation: token.PassLuhnValidation,
				Token:              token.Token,
				Created:            json.Number(token.Created),
				Customer:           CustomerOnlyHref(customer).Href,
				AdditionalDetails:  tokenParams.AdditionalDetails,
				BinNumber:          json.Number(token.BinNumber),
				Vendor:             token.Vendor,
				Issuer:             token.Issuer,
				CardType:           token.CardType,
				Level:              token.Level,
				CountryCode:        token.CountryCode,
				HolderName:         tokenParams.HolderName,
				ExpirationDate:     normalizeExpirationDate(tokenParams.ExpirationDate),
				Last4Digits:        last4(tokenParams.CardNumber),
				ShippingAddress:    nil, // why missing?
				BillingAddress:     tokenParams.BillingAddress,
				FingerPrint:        paymentMethodCreated.FingerPrint, // ignore
			}, paymentMethodCreated)
		})

		paymentMethodRetrieved, err := client.PaymentMethod().Get(context.Background(), customer.ID, token.Token)
		require.NoError(t, err)
		require.Equal(t, paymentMethodCreated, paymentMethodRetrieved)

		paymentMethodsRetrieved, err := client.PaymentMethod().GetList(context.Background(), customer.ID)
		require.NoError(t, err)
		require.ElementsMatch(t, []zooz.PaymentMethod{*paymentMethodCreated}, paymentMethodsRetrieved)
	})

	t.Run("get-list (several)", func(t *testing.T) {
		t.Parallel()

		customer := PrepareCustomer(t, client)
		token1, _ := PrepareToken(t, client)
		token2, _ := PrepareToken(t, client)

		paymentMethod1, err := client.PaymentMethod().New(context.Background(), randomString(32), customer.ID, token1.Token)
		require.NoError(t, err)

		paymentMethod2, err := client.PaymentMethod().New(context.Background(), randomString(32), customer.ID, token2.Token)
		require.NoError(t, err)

		paymentMethodsRetrieved, err := client.PaymentMethod().GetList(context.Background(), customer.ID)
		require.NoError(t, err)
		require.ElementsMatch(t, []zooz.PaymentMethod{*paymentMethod1, *paymentMethod2}, paymentMethodsRetrieved)
	})

	t.Run("delete & get & get-list", func(t *testing.T) {
		t.Parallel()

		customer := PrepareCustomer(t, client)
		token, _ := PrepareToken(t, client)

		_, err := client.PaymentMethod().New(context.Background(), randomString(32), customer.ID, token.Token)
		require.NoError(t, err)

		err = client.PaymentMethod().Delete(context.Background(), customer.ID, token.Token)
		require.NoError(t, err)

		paymentMethodsRetrieved, err := client.PaymentMethod().GetList(context.Background(), customer.ID)
		require.NoError(t, err)
		require.Empty(t, paymentMethodsRetrieved)
	})

	t.Run("idempotency key means nothing, token is idempotency key", func(t *testing.T) {
		t.Parallel()

		idempotencyKey1 := randomString(32)
		idempotencyKey2 := randomString(32) // different key -> same payment method
		token1, _ := PrepareToken(t, client)
		token2, _ := PrepareToken(t, client) // different token -> new payment method
		customer := PrepareCustomer(t, client)

		paymentMethod1, err := client.PaymentMethod().New(context.Background(), idempotencyKey1, customer.ID, token1.Token)
		require.NoError(t, err)

		paymentMethod2, err := client.PaymentMethod().New(context.Background(), idempotencyKey1, customer.ID, token1.Token)
		require.NoError(t, err)
		require.Equal(t, paymentMethod1, paymentMethod2)

		paymentMethod3, err := client.PaymentMethod().New(context.Background(), idempotencyKey2, customer.ID, token1.Token) // different key
		require.NoError(t, err)
		require.Equal(t, paymentMethod1, paymentMethod3)

		paymentMethod4, err := client.PaymentMethod().New(context.Background(), idempotencyKey1, customer.ID, token2.Token) // different token
		require.NoError(t, err)
		require.NotEqual(t, paymentMethod1.Token, paymentMethod4.Token)
	})

	t.Run("get - unknown token", func(t *testing.T) {
		t.Parallel()

		customer := PrepareCustomer(t, client)

		_, err := client.PaymentMethod().Get(context.Background(), customer.ID, UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "",
		})
	})

	t.Run("get - unknown customer", func(t *testing.T) {
		t.Parallel()

		token, _ := PrepareToken(t, client)

		_, err := client.PaymentMethod().Get(context.Background(), UnknownUUID, token.Token)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "",
		})
	})

	t.Run("get-list - unknown customer", func(t *testing.T) {
		t.Parallel()

		_, err := client.PaymentMethod().GetList(context.Background(), UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "",
		})
	})
}
