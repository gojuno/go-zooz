package sandbox

import (
	"context"
	"github.com/gtforge/go-zooz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRefund(t *testing.T) {
	t.Parallel()
	client := GetClient(t)

	t.Run("new (successful, different amount) & get", func(t *testing.T) {
		t.Parallel()

		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 5000, nil)
		_ = PrepareAuthorization(t, client, payment, token)
		_ = PrepareCapture(t, client, payment)
		reconciliationID := randomString(32)

		refundCreated, err := client.Refund().New(context.Background(), randomString(32), payment.ID, &zooz.RefundParams{
			ReconciliationID: reconciliationID,
			Amount:           4500,
			CaptureID:        "",
			Reason:           "xxx",
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, refundCreated.ID)
			assert.NotEmpty(t, refundCreated.Created)
			assert.NotEmpty(t, refundCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, refundCreated.ProviderData.AuthorizationCode)
			assert.NotEmpty(t, refundCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, refundCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, refundCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Refund{
				RefundParams: zooz.RefundParams{
					ReconciliationID: reconciliationID,
					Amount:           4500,
					CaptureID:        "",
					Reason:           "", // why empty?
				},
				ID: refundCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Succeed",
					Category:    "",
					SubCategory: "",
					Description: "",
				},
				Created: refundCreated.Created, // ignore
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "0",
					Description:           "Refunded.",
					RawResponse:           refundCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     refundCreated.ProviderData.AuthorizationCode, // ignore
					TransactionID:         refundCreated.ProviderData.TransactionID,     // ignore
					ExternalID:            refundCreated.ProviderData.ExternalID,        // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				AdditionalDetails:     nil,
				ProviderConfiguration: refundCreated.ProviderConfiguration, // ignore
			}, refundCreated)
		})

		refundRetrieved, err := client.Refund().Get(context.Background(), payment.ID, refundCreated.ID)
		require.NoError(t, err)
		require.Equal(t, refundCreated, refundRetrieved)
	})

	t.Run("new (successful, w/o explicit amount) & get", func(t *testing.T) {
		t.Parallel()

		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 5000, nil)
		_ = PrepareAuthorization(t, client, payment, token)
		_ = PrepareCapture(t, client, payment)
		reconciliationID := randomString(32)

		refundCreated, err := client.Refund().New(context.Background(), randomString(32), payment.ID, &zooz.RefundParams{
			ReconciliationID: reconciliationID,
			Amount:           0, // should use amount from payment
			CaptureID:        "",
			Reason:           "xxx",
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, refundCreated.ID)
			assert.NotEmpty(t, refundCreated.Created)
			assert.NotEmpty(t, refundCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, refundCreated.ProviderData.AuthorizationCode)
			assert.NotEmpty(t, refundCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, refundCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, refundCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Refund{
				RefundParams: zooz.RefundParams{
					ReconciliationID: reconciliationID,
					Amount:           5000,
					CaptureID:        "",
					Reason:           "", // why empty?
				},
				ID: refundCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Succeed",
					Category:    "",
					SubCategory: "",
					Description: "",
				},
				Created: refundCreated.Created, // ignore
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "0",
					Description:           "Refunded.",
					RawResponse:           refundCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     refundCreated.ProviderData.AuthorizationCode, // ignore
					TransactionID:         refundCreated.ProviderData.TransactionID,     // ignore
					ExternalID:            refundCreated.ProviderData.ExternalID,        // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				AdditionalDetails:     nil,
				ProviderConfiguration: refundCreated.ProviderConfiguration, // ignore
			}, refundCreated)
		})

		refundRetrieved, err := client.Refund().Get(context.Background(), payment.ID, refundCreated.ID)
		require.NoError(t, err)
		require.Equal(t, refundCreated, refundRetrieved)
	})

	t.Run("new (failed) & get", func(t *testing.T) {
		t.Parallel()

		const amount = 1000 // to fail capture, see https://developers.paymentsos.com/docs/testing/mockprovider-reference.html#capture-or-refund-requests

		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 5000, nil)
		_ = PrepareAuthorization(t, client, payment, token)
		_ = PrepareCapture(t, client, payment)
		reconciliationID := randomString(32)

		refundCreated, err := client.Refund().New(context.Background(), randomString(32), payment.ID, &zooz.RefundParams{
			ReconciliationID: reconciliationID,
			Amount:           amount,
			CaptureID:        "",
			Reason:           "xxx",
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, refundCreated.ID)
			assert.NotEmpty(t, refundCreated.Created)
			assert.NotEmpty(t, refundCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, refundCreated.ProviderData.AuthorizationCode)
			assert.NotEmpty(t, refundCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, refundCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, refundCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Refund{
				RefundParams: zooz.RefundParams{
					ReconciliationID: reconciliationID,
					Amount:           amount,
					CaptureID:        "",
					Reason:           "", // why empty?
				},
				ID: refundCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Failed",
					Category:    "provider_error",
					SubCategory: "",
					Description: "Something went wrong on the provider's side.",
				},
				Created: refundCreated.Created, // ignore
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "102",
					Description:           "refund failed",
					RawResponse:           refundCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     refundCreated.ProviderData.AuthorizationCode, // ignore
					TransactionID:         refundCreated.ProviderData.TransactionID,     // ignore
					ExternalID:            refundCreated.ProviderData.ExternalID,        // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				AdditionalDetails:     nil,
				ProviderConfiguration: refundCreated.ProviderConfiguration, // ignore
			}, refundCreated)
		})

		refundRetrieved, err := client.Refund().Get(context.Background(), payment.ID, refundCreated.ID)
		require.NoError(t, err)
		require.Equal(t, refundCreated, refundRetrieved)
	})

	t.Run("idempotency", func(t *testing.T) {
		t.Parallel()

		idempotencyKey1 := randomString(32)
		idempotencyKey2 := randomString(32) // different key -> new refund
		const amount1, amount2 = 4500, 4000 // can't change amount
		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 5000, nil)
		_ = PrepareAuthorization(t, client, payment, token)
		_ = PrepareCapture(t, client, payment)

		refund1, err := client.Refund().New(context.Background(), idempotencyKey1, payment.ID, &zooz.RefundParams{Amount: amount1})
		require.NoError(t, err)

		refund2, err := client.Refund().New(context.Background(), idempotencyKey1, payment.ID, &zooz.RefundParams{Amount: amount1})
		require.NoError(t, err)
		require.Equal(t, refund1, refund2)

		refund3, err := client.Refund().New(context.Background(), idempotencyKey1, payment.ID, &zooz.RefundParams{Amount: amount2}) // can't change amount
		require.NoError(t, err)
		require.Equal(t, refund1, refund3)

		refund4, err := client.Refund().New(context.Background(), idempotencyKey2, payment.ID, &zooz.RefundParams{Amount: amount1}) // different key
		require.NoError(t, err)
		require.NotEqual(t, refund1.ID, refund4.ID)
	})

	t.Run("new - can't create refund without capture", func(t *testing.T) {
		t.Parallel()

		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 10000, nil)
		_ = PrepareAuthorization(t, client, payment, token)

		_, err := client.Refund().New(context.Background(), randomString(32), payment.ID, &zooz.RefundParams{
			ReconciliationID: randomString(32),
			Amount:           5000,
			CaptureID:        "",
			Reason:           "abc",
		})
		requireZoozError(t, err, http.StatusConflict, zooz.APIError{
			Category:    "api_request_error",
			Description: "There was conflict with payment resource current state.",
			MoreInfo:    "Please check the current state of the payment.",
		})
	})

	t.Run("get - unknown refund", func(t *testing.T) {
		t.Parallel()

		payment := PreparePayment(t, client, 123, nil)

		_, err := client.Refund().Get(context.Background(), payment.ID, UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "refunds resource does not exits",
		})
	})

	t.Run("get - unknown payment", func(t *testing.T) {
		t.Parallel()

		_, err := client.Refund().Get(context.Background(), UnknownUUID, UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "Payment resource does not exists",
		})
	})
}
