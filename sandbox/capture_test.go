package sandbox

import (
	"context"
	"github.com/gtforge/go-zooz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCapture(t *testing.T) {
	t.Parallel()
	client := GetClient(t)

	t.Run("new (successful, different amount) & get", func(t *testing.T) {
		t.Parallel()

		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 5000, nil)
		_ = PrepareAuthorization(t, client, payment, token)
		reconciliationID := randomString(32)

		captureCreated, err := client.Capture().New(context.Background(), randomString(32), payment.ID, &zooz.CaptureParams{
			ReconciliationID: reconciliationID,
			Amount:           4500,
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, captureCreated.ID)
			assert.NotEmpty(t, captureCreated.Created)
			assert.NotEmpty(t, captureCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, captureCreated.ProviderData.AuthorizationCode)
			assert.NotEmpty(t, captureCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, captureCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, captureCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Capture{
				CaptureParams: zooz.CaptureParams{
					ReconciliationID: reconciliationID,
					Amount:           4500,
				},
				ID: captureCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Succeed",
					Category:    "",
					SubCategory: "",
					Description: "",
				},
				Created: captureCreated.Created, // ignore
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "0",
					Description:           "Captured.",
					RawResponse:           captureCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     captureCreated.ProviderData.AuthorizationCode, // ignore
					TransactionID:         captureCreated.ProviderData.TransactionID,     // ignore
					ExternalID:            captureCreated.ProviderData.ExternalID,        // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				ProviderSpecificData: nil,
				Level23: zooz.Level23{
					OrderID:             nil,
					TaxMode:             "",
					TaxAmount:           0,
					ShippingAmount:      0,
					FromShippingZipCode: "",
					DutyAmount:          0,
					DiscountAmount:      0,
					LineItems:           nil,
					ShippingAddress:     nil,
				},
				ProviderConfiguration: captureCreated.ProviderConfiguration, // ignore
				AdditionalDetails:     nil,
			}, captureCreated)
		})

		captureRetrieved, err := client.Capture().Get(context.Background(), payment.ID, captureCreated.ID)
		require.NoError(t, err)
		require.Equal(t, captureCreated, captureRetrieved)
	})

	t.Run("new (successful, w/o explicit amount) & get", func(t *testing.T) {
		t.Parallel()

		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 5000, nil)
		_ = PrepareAuthorization(t, client, payment, token)
		reconciliationID := randomString(32)

		captureCreated, err := client.Capture().New(context.Background(), randomString(32), payment.ID, &zooz.CaptureParams{
			ReconciliationID: reconciliationID,
			Amount:           0, // should use amount from payment
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, captureCreated.ID)
			assert.NotEmpty(t, captureCreated.Created)
			assert.NotEmpty(t, captureCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, captureCreated.ProviderData.AuthorizationCode)
			assert.NotEmpty(t, captureCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, captureCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, captureCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Capture{
				CaptureParams: zooz.CaptureParams{
					ReconciliationID: reconciliationID,
					Amount:           5000,
				},
				ID: captureCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Succeed",
					Category:    "",
					SubCategory: "",
					Description: "",
				},
				Created: captureCreated.Created, // ignore
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "0",
					Description:           "Captured.",
					RawResponse:           captureCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     captureCreated.ProviderData.AuthorizationCode, // ignore
					TransactionID:         captureCreated.ProviderData.TransactionID,     // ignore
					ExternalID:            captureCreated.ProviderData.ExternalID,        // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				ProviderSpecificData: nil,
				Level23: zooz.Level23{
					OrderID:             nil,
					TaxMode:             "",
					TaxAmount:           0,
					ShippingAmount:      0,
					FromShippingZipCode: "",
					DutyAmount:          0,
					DiscountAmount:      0,
					LineItems:           nil,
					ShippingAddress:     nil,
				},
				ProviderConfiguration: captureCreated.ProviderConfiguration, // ignore
				AdditionalDetails:     nil,
			}, captureCreated)
		})

		captureRetrieved, err := client.Capture().Get(context.Background(), payment.ID, captureCreated.ID)
		require.NoError(t, err)
		require.Equal(t, captureCreated, captureRetrieved)
	})

	t.Run("new (failed) & get", func(t *testing.T) {
		t.Parallel()

		const amount = 1000 // to fail capture, see https://developers.paymentsos.com/docs/testing/mockprovider-reference.html#capture-or-refund-requests

		token, _ := PrepareToken(t, client)
		payment := PreparePayment(t, client, 5000, nil)
		_ = PrepareAuthorization(t, client, payment, token)
		reconciliationID := randomString(32)

		captureCreated, err := client.Capture().New(context.Background(), randomString(32), payment.ID, &zooz.CaptureParams{
			ReconciliationID: reconciliationID,
			Amount:           amount,
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, captureCreated.ID)
			assert.NotEmpty(t, captureCreated.Created)
			assert.NotEmpty(t, captureCreated.ProviderData.RawResponse)
			assert.NotEmpty(t, captureCreated.ProviderData.AuthorizationCode)
			assert.NotEmpty(t, captureCreated.ProviderData.TransactionID)
			assert.NotEmpty(t, captureCreated.ProviderData.ExternalID)
			assert.NotEmpty(t, captureCreated.ProviderConfiguration)
			assert.Equal(t, &zooz.Capture{
				CaptureParams: zooz.CaptureParams{
					ReconciliationID: reconciliationID,
					Amount:           amount,
				},
				ID: captureCreated.ID, // ignore
				Result: zooz.Result{
					Status:      "Failed",
					Category:    "provider_error",
					SubCategory: "",
					Description: "Something went wrong on the provider's side.",
				},
				Created: captureCreated.Created, // ignore
				ProviderData: zooz.ProviderData{
					ProviderName:          "MockProcessor",
					ResponseCode:          "102",
					Description:           "capture failed.",
					RawResponse:           captureCreated.ProviderData.RawResponse, // ignore
					AvsCode:               "",
					AuthorizationCode:     captureCreated.ProviderData.AuthorizationCode, // ignore
					TransactionID:         captureCreated.ProviderData.TransactionID,     // ignore
					ExternalID:            captureCreated.ProviderData.ExternalID,        // ignore
					Documents:             nil,
					AdditionalInformation: nil,
					NetworkTransactionID:  "",
				},
				ProviderSpecificData: nil,
				Level23: zooz.Level23{
					OrderID:             nil,
					TaxMode:             "",
					TaxAmount:           0,
					ShippingAmount:      0,
					FromShippingZipCode: "",
					DutyAmount:          0,
					DiscountAmount:      0,
					LineItems:           nil,
					ShippingAddress:     nil,
				},
				ProviderConfiguration: captureCreated.ProviderConfiguration, // ignore
				AdditionalDetails:     nil,
			}, captureCreated)
		})

		captureRetrieved, err := client.Capture().Get(context.Background(), payment.ID, captureCreated.ID)
		require.NoError(t, err)
		require.Equal(t, captureCreated, captureRetrieved)
	})

	t.Run("new - can't create capture without authorization", func(t *testing.T) {
		t.Parallel()

		payment := PreparePayment(t, client, 10000, nil)

		_, err := client.Capture().New(context.Background(), randomString(32), payment.ID, &zooz.CaptureParams{
			ReconciliationID: randomString(32),
			Amount:           5000,
		})
		requireZoozError(t, err, http.StatusConflict, zooz.APIError{
			Category:    "api_request_error",
			Description: "There was conflict with payment resource current state.",
			MoreInfo:    "Please check the current state of the payment.",
		})
	})

	t.Run("get - unknown capture", func(t *testing.T) {
		t.Parallel()

		payment := PreparePayment(t, client, 123, nil)

		_, err := client.Capture().Get(context.Background(), payment.ID, UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "captures resource does not exits",
		})
	})

	t.Run("get - unknown payment", func(t *testing.T) {
		t.Parallel()

		_, err := client.Capture().Get(context.Background(), UnknownUUID, UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "Payment resource does not exists",
		})
	})
}
