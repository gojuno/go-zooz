package sandbox

import (
	"context"
	"github.com/gtforge/go-zooz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestVoid(t *testing.T) {
	t.Parallel()
	client := GetClient(t)

	t.Run("new (successful) & get", func(t *testing.T) {
		t.Parallel()

		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 5000, nil)
		_ = PrepareAuthorization(t, client, payment, token)

		voidCreated, err := client.Void().New(context.Background(), randomString(32), payment.ID)
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, voidCreated.ID)
			assert.NotEmpty(t, voidCreated.Created)
			assert.NotEmpty(t, voidCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, voidCreated.ProviderData.AuthorizationCode)
			assert.NotEmpty(t, voidCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, voidCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, voidCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Void{
				ID: voidCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Succeed",
					Category:    "",
					SubCategory: "",
					Description: "",
				},
				Created: voidCreated.Created, // ignore
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "0",
					Description:           "Canceled.",
					RawResponse:           voidCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     voidCreated.ProviderData.AuthorizationCode, // ignore
					TransactionID:         voidCreated.ProviderData.TransactionID,     // ignore
					ExternalID:            voidCreated.ProviderData.ExternalID,        // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				AdditionalDetails:     nil,
				ProviderConfiguration: voidCreated.ProviderConfiguration, // ignore
			}, voidCreated)
		})

		voidRetrieved, err := client.Void().Get(context.Background(), payment.ID, voidCreated.ID)
		require.NoError(t, err)
		require.Equal(t, voidCreated, voidRetrieved)
	})

	t.Run("new (failed) & get", func(t *testing.T) {
		t.Parallel()

		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 1000, nil)
		_ = PrepareAuthorization(t, client, payment, token)

		voidCreated, err := client.Void().New(context.Background(), randomString(32), payment.ID)
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, voidCreated.ID)
			assert.NotEmpty(t, voidCreated.Created)
			assert.NotEmpty(t, voidCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, voidCreated.ProviderData.AuthorizationCode)
			assert.NotEmpty(t, voidCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, voidCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, voidCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Void{
				ID: voidCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Failed",
					Category:    "provider_error",
					SubCategory: "",
					Description: "Something went wrong on the provider's side.",
				},
				Created: voidCreated.Created, // ignore
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "102",
					Description:           "void failed.",
					RawResponse:           voidCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     voidCreated.ProviderData.AuthorizationCode, // ignore
					TransactionID:         voidCreated.ProviderData.TransactionID,     // ignore
					ExternalID:            voidCreated.ProviderData.ExternalID,        // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				AdditionalDetails:     nil,
				ProviderConfiguration: voidCreated.ProviderConfiguration, // ignore
			}, voidCreated)
		})

		voidRetrieved, err := client.Void().Get(context.Background(), payment.ID, voidCreated.ID)
		require.NoError(t, err)
		require.Equal(t, voidCreated, voidRetrieved)
	})

	t.Run("idempotency", func(t *testing.T) {
		t.Parallel()

		idempotencyKey1 := randomString(32)
		idempotencyKey2 := randomString(32) // different key -> new void
		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 5000, nil)
		_ = PrepareAuthorization(t, client, payment, token)

		void1, err := client.Void().New(context.Background(), idempotencyKey1, payment.ID)
		require.NoError(t, err)

		void2, err := client.Void().New(context.Background(), idempotencyKey1, payment.ID)
		require.NoError(t, err)
		require.Equal(t, void1, void2)

		void3, err := client.Void().New(context.Background(), idempotencyKey2, payment.ID) // different key
		require.NoError(t, err)
		require.NotEqual(t, void1.ID, void3.ID)
	})

	t.Run("new - can't create void without authorization", func(t *testing.T) {
		t.Parallel()

		payment := PreparePayment(t, client, 1000, nil)

		_, err := client.Void().New(context.Background(), randomString(32), payment.ID)
		requireZoozError(t, err, http.StatusConflict, zooz.APIError{
			Category:    "api_request_error",
			Description: "There was conflict with payment resource current state.",
			MoreInfo:    "Please check the current state of the payment.",
		})
	})

	t.Run("get - unknown void", func(t *testing.T) {
		t.Parallel()

		payment := PreparePayment(t, client, 123, nil)

		_, err := client.Void().Get(context.Background(), payment.ID, UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "voids resource does not exits",
		})
	})

	t.Run("get - unknown payment", func(t *testing.T) {
		t.Parallel()

		_, err := client.Void().Get(context.Background(), UnknownUUID, UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "Payment resource does not exists",
		})
	})
}
