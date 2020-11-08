package zooz_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gtforge/go-zooz"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestDecodeWebhookRequest_Payment(t *testing.T) {
	expected := zooz.PaymentCallback{
		CallbackCommon: zooz.CallbackCommon{
			EventType:      "payment.payment.create",
			XPaymentsOSEnv: "live",
			XZoozRequestID: "test-x-zooz-request-id",

			ID:        "test-webhook-id",
			Created:   time.Date(2018, 10, 03, 04, 58, 35, 385000000, time.UTC), // "2018-10-03T04:58:35.385Z"
			AccountID: "test-account-id",
			AppID:     "t",
			PaymentID: "test-payment-id",
		},
		Data: zooz.Payment{
			ID: "test-transaction-id",
			PaymentParams: zooz.PaymentParams{
				Amount:     1000,
				Currency:   "USD",
				CustomerID: "test-customer-id",
			},
			Status: zooz.PaymentStatusAuthorized,
		},
	}

	keyProvider := zooz.FixedPrivateKeyProvider{expected.AppID: []byte("test-private-key")}

	body := `{
		"id": "` + expected.ID + `",
		"created": "2018-10-03T04:58:35.385Z",
		"account_id": "` + expected.AccountID + `",
		"app_id": "` + expected.AppID + `",
		"payment_id": "` + expected.PaymentID + `",
		"data": {
			"id": "` + expected.Data.ID + `",
			"status": "` + string(expected.Data.Status) + `",
			"amount": ` + strconv.FormatInt(expected.Data.Amount, 10) + `,
			"currency":"` + expected.Data.Currency + `",
			"customer_id": "` + expected.Data.CustomerID + `"
		}
	}`
	h := http.Header{}
	h.Set("event-type", expected.EventType)
	h.Set("x-payments-os-env", expected.XPaymentsOSEnv)
	h.Set("x-zooz-request-id", expected.XZoozRequestID)
	h.Set("signature", calcRequestSignature(t, keyProvider, []byte(body), h))

	cb, err := zooz.DecodeWebhookRequest(context.Background(), []byte(body), h, keyProvider)
	require.NoError(t, err)
	require.Equal(t, expected, cb)
}

func TestDecodeWebhookRequest_Authorization(t *testing.T) {
	expected := zooz.AuthorizationCallback{
		CallbackCommon: zooz.CallbackCommon{
			EventType:      "payment.authorization.create",
			XPaymentsOSEnv: "live",
			XZoozRequestID: "test-x-zooz-request-id",

			ID:        "test-webhook-id",
			Created:   time.Date(2018, 10, 03, 04, 58, 35, 385000000, time.UTC), // "2018-10-03T04:58:35.385Z"
			AccountID: "test-account-id",
			AppID:     "test-app-id",
			PaymentID: "test-payment-id",
		},
		Data: zooz.Authorization{
			ID: "test-transaction-id",
			Result: zooz.Result{
				Status:      "Pending",
				Category:    "payment_method_declined",
				SubCategory: "declined_by_issuing_bank",
				Description: "The transaction was declined by the Issuing bank.",
			},
			Amount:           1000,
			ReconciliationID: "test-reconciliation-id",
			ProviderData: zooz.ProviderData{
				ProviderName: "test-provider-name",
				ResponseCode: "test-response-code",
			},
		},
	}

	keyProvider := zooz.FixedPrivateKeyProvider{expected.AppID: []byte("test-private-key")}

	body := `{
		"id": "` + expected.ID + `",
		"created": "2018-10-03T04:58:35.385Z",
		"account_id": "` + expected.AccountID + `",
		"app_id": "` + expected.AppID + `",
		"payment_id": "` + expected.PaymentID + `",
		"data": {
			"id": "` + expected.Data.ID + `",
			"result": {
				"status": "` + expected.Data.Result.Status + `",
				"category": "` + expected.Data.Result.Category + `",
				"sub_category": "` + expected.Data.Result.SubCategory + `",
				"description": "` + expected.Data.Result.Description + `"
			},
			"amount": ` + strconv.FormatInt(expected.Data.Amount, 10) + `,
			"reconciliation_id": "` + expected.Data.ReconciliationID + `",
			"provider_data": {
				"provider_name": "` + expected.Data.ProviderData.ProviderName + `",
				"response_code": "` + expected.Data.ProviderData.ResponseCode + `"
			}
		}
	}`
	h := http.Header{}
	h.Set("event-type", expected.EventType)
	h.Set("x-payments-os-env", expected.XPaymentsOSEnv)
	h.Set("x-zooz-request-id", expected.XZoozRequestID)
	h.Set("signature", calcRequestSignature(t, keyProvider, []byte(body), h))

	cb, err := zooz.DecodeWebhookRequest(context.Background(), []byte(body), h, keyProvider)
	require.NoError(t, err)
	require.Equal(t, expected, cb)
}

func TestDecodeWebhookRequest_Capture(t *testing.T) {
	expected := zooz.CaptureCallback{
		CallbackCommon: zooz.CallbackCommon{
			EventType:      "payment.capture.update",
			XPaymentsOSEnv: "live",
			XZoozRequestID: "test-x-zooz-request-id",

			ID:        "test-webhook-id",
			Created:   time.Date(2018, 10, 03, 05, 14, 17, 196000000, time.UTC), // "2018-10-03T05:14:17.196Z"
			AccountID: "test-account-id",
			AppID:     "test-app-id",
			PaymentID: "test-payment-id",
		},
		Data: zooz.Capture{
			ID: "test-transaction-id",
			Result: zooz.Result{
				Status:      "Pending",
				Category:    "payment_method_declined",
				SubCategory: "declined_by_issuing_bank",
				Description: "The transaction was declined by the Issuing bank.",
			},
			CaptureParams: zooz.CaptureParams{
				ReconciliationID: "test-reconciliation-id",
				Amount:           2000,
			},
			ProviderData: zooz.ProviderData{
				ProviderName: "test-provider-name",
				ResponseCode: "test-response-code",
			},
		},
	}

	keyProvider := zooz.FixedPrivateKeyProvider{expected.AppID: []byte("test-private-key")}

	body := `{
		"id": "` + expected.ID + `",
		"created": "2018-10-03T05:14:17.196Z",
		"account_id": "` + expected.AccountID + `",
		"app_id": "` + expected.AppID + `",
		"payment_id": "` + expected.PaymentID + `",
		"data": {
			"id": "` + expected.Data.ID + `",
			"result": {
				"status": "` + expected.Data.Result.Status + `",
				"category": "` + expected.Data.Result.Category + `",
				"sub_category": "` + expected.Data.Result.SubCategory + `",
				"description": "` + expected.Data.Result.Description + `"
			},
			"amount": ` + strconv.FormatInt(expected.Data.Amount, 10) + `,
			"reconciliation_id": "` + expected.Data.ReconciliationID + `",
			"provider_data": {
				"provider_name": "` + expected.Data.ProviderData.ProviderName + `",
				"response_code": "` + expected.Data.ProviderData.ResponseCode + `"
			}
		}
	}`
	h := http.Header{}
	h.Set("event-type", expected.EventType)
	h.Set("x-payments-os-env", expected.XPaymentsOSEnv)
	h.Set("x-zooz-request-id", expected.XZoozRequestID)
	h.Set("signature", calcRequestSignature(t, keyProvider, []byte(body), h))

	cb, err := zooz.DecodeWebhookRequest(context.Background(), []byte(body), h, keyProvider)
	require.NoError(t, err)
	require.Equal(t, expected, cb)
}

func TestDecodeWebhookRequest_Void(t *testing.T) {
	expected := zooz.VoidCallback{
		CallbackCommon: zooz.CallbackCommon{
			EventType:      "payment.void.create",
			XPaymentsOSEnv: "live",
			XZoozRequestID: "test-x-zooz-request-id",

			ID:        "test-webhook-id",
			Created:   time.Date(2018, 10, 03, 05, 14, 17, 196000000, time.UTC), // "2018-10-03T05:14:17.196Z"
			AccountID: "test-account-id",
			AppID:     "test-app-id",
			PaymentID: "test-payment-id",
		},
		Data: zooz.Void{
			ID: "test-transaction-id",
			Result: zooz.Result{
				Status:      "Pending",
				Category:    "payment_method_declined",
				SubCategory: "declined_by_issuing_bank",
				Description: "The transaction was declined by the Issuing bank.",
			},
			ProviderData: zooz.ProviderData{
				ProviderName: "test-provider-name",
				ResponseCode: "test-response-code",
			},
		},
	}

	keyProvider := zooz.FixedPrivateKeyProvider{expected.AppID: []byte("test-private-key")}

	body := `{
		"id": "` + expected.ID + `",
		"created": "2018-10-03T05:14:17.196Z",
		"account_id": "` + expected.AccountID + `",
		"app_id": "` + expected.AppID + `",
		"payment_id": "` + expected.PaymentID + `",
		"data": {
			"id": "` + expected.Data.ID + `",
			"result": {
				"status": "` + expected.Data.Result.Status + `",
				"category": "` + expected.Data.Result.Category + `",
				"sub_category": "` + expected.Data.Result.SubCategory + `",
				"description": "` + expected.Data.Result.Description + `"
			},
			"provider_data": {
				"provider_name": "` + expected.Data.ProviderData.ProviderName + `",
				"response_code": "` + expected.Data.ProviderData.ResponseCode + `"
			}
		}
	}`
	h := http.Header{}
	h.Set("event-type", expected.EventType)
	h.Set("x-payments-os-env", expected.XPaymentsOSEnv)
	h.Set("x-zooz-request-id", expected.XZoozRequestID)
	h.Set("signature", calcRequestSignature(t, keyProvider, []byte(body), h))

	cb, err := zooz.DecodeWebhookRequest(context.Background(), []byte(body), h, keyProvider)
	require.NoError(t, err)
	require.Equal(t, expected, cb)
}

func TestDecodeWebhookRequest_Refund(t *testing.T) {
	expected := zooz.RefundCallback{
		CallbackCommon: zooz.CallbackCommon{
			EventType:      "payment.refund.update",
			XPaymentsOSEnv: "live",
			XZoozRequestID: "test-x-zooz-request-id",

			ID:        "test-webhook-id",
			Created:   time.Date(2018, 10, 03, 05, 22, 45, 610000000, time.UTC), // "2018-10-03T05:22:45.610Z"
			AccountID: "test-account-id",
			AppID:     "test-app-id",
			PaymentID: "test-payment-id",
		},
		Data: zooz.Refund{
			ID: "test-transaction-id",
			Result: zooz.Result{
				Status:      "Pending",
				Category:    "payment_method_declined",
				SubCategory: "declined_by_issuing_bank",
				Description: "The transaction was declined by the Issuing bank.",
			},
			RefundParams: zooz.RefundParams{
				ReconciliationID: "test-reconciliation-id",
				Amount:           2000,
				CaptureID:        "test-capture-id",
				Reason:           "reason for the refund",
			},
			ProviderData: zooz.ProviderData{
				ProviderName: "test-provider-name",
				ResponseCode: "test-response-code",
			},
		},
	}

	keyProvider := zooz.FixedPrivateKeyProvider{expected.AppID: []byte("test-private-key")}

	body := `{
		"id": "` + expected.ID + `",
		"created": "2018-10-03T05:22:45.610Z",
		"account_id": "` + expected.AccountID + `",
		"app_id": "` + expected.AppID + `",
		"payment_id": "` + expected.PaymentID + `",
		"data": {
			"id": "` + expected.Data.ID + `",
			"result": {
				"status": "` + expected.Data.Result.Status + `",
				"category": "` + expected.Data.Result.Category + `",
				"sub_category": "` + expected.Data.Result.SubCategory + `",
				"description": "` + expected.Data.Result.Description + `"
  			},
  			"amount": ` + strconv.FormatInt(expected.Data.Amount, 10) + `,
			"reconciliation_id": "` + expected.Data.ReconciliationID + `",
			"capture_id": "` + expected.Data.CaptureID + `",
			"reason": "` + expected.Data.Reason + `",
			"provider_data": {
				"provider_name": "` + expected.Data.ProviderData.ProviderName + `",
				"response_code": "` + expected.Data.ProviderData.ResponseCode + `"
			}
		}
	}`
	h := http.Header{}
	h.Set("event-type", expected.EventType)
	h.Set("x-payments-os-env", expected.XPaymentsOSEnv)
	h.Set("x-zooz-request-id", expected.XZoozRequestID)
	h.Set("signature", calcRequestSignature(t, keyProvider, []byte(body), h))

	cb, err := zooz.DecodeWebhookRequest(context.Background(), []byte(body), h, keyProvider)
	require.NoError(t, err)
	require.Equal(t, expected, cb)
}

func TestDecodeWebhookRequest_BadRequestError_BrokenJson(t *testing.T) {
	keyProvider := zooz.FixedPrivateKeyProvider{"test-app-id": []byte("test-private-key")}

	body := `{
		"id": "test-webhook-id",
		"created": "2018-10-03T05:22:45.610Z",
		"account_id": "test-account-id",
		"app_id": "test-app-id",
`
	h := http.Header{}
	h.Set("event-type", "payment.void.create")
	h.Set("signature", "sig1=doesntmatter")

	_, err := zooz.DecodeWebhookRequest(context.Background(), []byte(body), h, keyProvider)
	require.Error(t, err)
	require.True(t, errors.As(err, &zooz.ErrBadRequest{}))
	require.Contains(t, err.Error(), "unmarshal request body")
}

func TestDecodeWebhookRequest_BadRequestError_UnsupportedEventType(t *testing.T) {
	keyProvider := zooz.FixedPrivateKeyProvider{"test-app-id": []byte("test-private-key")}

	body := `{
		"id": "test-webhook-id",
		"created": "2018-10-03T05:14:17.196Z",
		"account_id": "test-account-id",
		"app_id": "test-app-id",
		"payment_id": "test-payment-id",
		"data": {
			"id": "test-transaction-id",
			"result": {
				"status": "Success"
			},
			"provider_data": {
				"provider_name": "test-provider-name",
				"response_code": "0"
			}
		}
	}`
	h := http.Header{}
	h.Set("event-type", "payment.XXXX.create")
	h.Set("signature", calcRequestSignature(t, keyProvider, []byte(body), h))

	_, err := zooz.DecodeWebhookRequest(context.Background(), []byte(body), h, keyProvider)
	require.Error(t, err)
	require.True(t, errors.As(err, &zooz.ErrBadRequest{}))
	require.Contains(t, err.Error(), `unsupported event type: "payment.XXXX.create"`)
}

func TestDecodeWebhookRequest_BadRequestError_UnknownBusinessUnit(t *testing.T) {
	keyProvider := zooz.FixedPrivateKeyProvider{"test-app-id": []byte("test-private-key")} // we don't know 'UNKNOWN-app-id'

	body := `{
		"id": "test-webhook-id",
		"created": "2018-10-03T05:14:17.196Z",
		"account_id": "test-account-id",
		"app_id": "UNKNOWN-app-id",
		"payment_id": "test-payment-id",
		"data": {
			"id": "test-transaction-id",
			"result": {
				"status": "Success"
			},
			"provider_data": {
				"provider_name": "test-provider-name",
				"response_code": "0"
			}
		}
	}`
	h := http.Header{}
	h.Set("event-type", "payment.void.create")
	h.Set("signature", "sig1=doesntmatter")

	_, err := zooz.DecodeWebhookRequest(context.Background(), []byte(body), h, keyProvider)
	require.Error(t, err)
	require.True(t, errors.As(err, &zooz.ErrBadRequest{}))
	require.Contains(t, err.Error(), `unknown app_id "UNKNOWN-app-id"`)
}

func TestDecodeWebhookRequest_BadRequestError_IncorrectSignature(t *testing.T) {
	keyProvider := zooz.FixedPrivateKeyProvider{"test-app-id": []byte("test-private-key")}

	body := `{
		"id": "test-webhook-id",
		"created": "2018-10-03T05:14:17.196Z",
		"account_id": "test-account-id",
		"app_id": "test-app-id",
		"payment_id": "test-payment-id",
		"data": {
			"id": "test-transaction-id",
			"result": {
				"status": "Success"
			},
			"provider_data": {
				"provider_name": "test-provider-name",
				"response_code": "0"
			}
		}
	}`
	h := http.Header{}
	h.Set("event-type", "payment.void.create")
	h.Set("signature", "sig1=iaminvalid")

	_, err := zooz.DecodeWebhookRequest(context.Background(), []byte(body), h, keyProvider)
	require.Error(t, err)
	require.True(t, errors.As(err, &zooz.ErrBadRequest{}))
	require.Contains(t, err.Error(), `incorrect signature`)
}

func TestCalculateWebhookSignature_AllFields(t *testing.T) {
	const (
		privateKey = "test-private-key"

		eventType                      = "payment.authorization.create"
		id                             = "test-id"
		created                        = "2018-10-03T05:14:17.196Z"
		accountID                      = "test-account-id"
		appID                          = "test-app-id"
		paymentID                      = "test-payment-id"
		data_ID                        = "test-transaction-id"
		data_Result_Status             = "Pending"
		data_Result_Category           = "payment_method_declined"
		data_Result_SubCategory        = "declined_by_issuing_bank"
		data_ProviderData_ResponseCode = "124"
		data_ReconciliationID          = "test-reconciliation-id"
		data_Amount                    = "42"
		data_Currency                  = "RUB"
	)

	keyProvider := zooz.FixedPrivateKeyProvider{appID: []byte(privateKey)}

	body := `{
		"id": "` + id + `",
		"created": "` + created + `",
		"account_id": "` + accountID + `",
		"app_id": "` + appID + `",
		"payment_id": "` + paymentID + `",
		"data": {
			"id": "` + data_ID + `",
			"result": {
				"status": "` + data_Result_Status + `",
				"category": "` + data_Result_Category + `",
				"sub_category": "` + data_Result_SubCategory + `"
			},
			"provider_data": {
				"response_code": "` + data_ProviderData_ResponseCode + `"
			},
			"reconciliation_id": "` + data_ReconciliationID + `",
			"amount": ` + data_Amount + `,
			"currency": "` + data_Currency + `"
		}
	}`
	h := http.Header{}
	h.Set("event-type", eventType)

	sign, err := zooz.CalculateWebhookSignature(context.Background(), []byte(body), h, keyProvider)
	require.NoError(t, err)
	// double check that all signature fields are covered
	require.Equal(t,
		signature([]string{
			eventType,
			id,
			accountID,
			paymentID,
			created,
			appID,
			data_ID,
			data_Result_Status,
			data_Result_Category,
			data_Result_SubCategory,
			data_ProviderData_ResponseCode,
			data_ReconciliationID,
			data_Amount,
			data_Currency,
		}, privateKey),
		sign)
	require.Equal(t, "b5457f3a7c7e8bdeea150485240b0a2c041f63d7ec5460c186ccde1e18453c3d", sign)
}

func TestCalculateWebhookSignature_MissingFields(t *testing.T) {
	const (
		privateKey = "test-private-key"

		eventType                      = "payment.authorization.create"
		id                             = "test-id"
		created                        = "2018-10-03T05:14:17.196Z"
		accountID                      = "test-account-id"
		appID                          = "test-app-id"
		paymentID                      = "test-app-id"
		data_ID                        = "test-transaction-id"
		data_Result_Status             = "Pending"
		data_ProviderData_ResponseCode = "124"
		data_ReconciliationID          = "test-reconciliation-id"
	)

	keyProvider := zooz.FixedPrivateKeyProvider{appID: []byte(privateKey)}

	body := `{
		"id": "` + id + `",
		"created": "` + created + `",
		"account_id": "` + accountID + `",
		"app_id": "` + appID + `",
		"payment_id": "` + paymentID + `",
		"data": {
			"id": "` + data_ID + `",
			"result": {
				"status": "` + data_Result_Status + `",
				"category": ""
			},
			"provider_data": {
				"response_code": "` + data_ProviderData_ResponseCode + `"
			},
			"reconciliation_id": "` + data_ReconciliationID + `"
		}
	}`
	h := http.Header{}
	h.Set("event-type", eventType)

	sign, err := zooz.CalculateWebhookSignature(context.Background(), []byte(body), h, keyProvider)
	require.NoError(t, err)
	require.Equal(t,
		signature([]string{
			eventType,
			id,
			accountID,
			paymentID,
			created,
			appID,
			data_ID,
			data_Result_Status,
			"", // empty category
			"", // missing subcategory
			data_ProviderData_ResponseCode,
			data_ReconciliationID,
			"", // no amount
			"", // no currency
		}, privateKey),
		sign)
	require.Equal(t, "f4a3827af9a1d3a33e33db9dc226e54c676e46295e2e27c3a8115280def4f21c", sign)
}

func TestCalculateWebhookSignature_NoAmount(t *testing.T) {
	keyProvider := zooz.FixedPrivateKeyProvider{"test-app-id": []byte("test-private-key")}

	body := `{
		"id": "test-id",
		"created": "2018-10-03T05:14:17.196Z",
		"account_id": "test-account-id",
		"app_id": "test-app-id",
		"payment_id": "test-payment-id",
		"data": {
			"id": "test-transaction-id",
			"result": {
				"status": "Pending",
				"category": "payment_method_declined",
				"sub_category": "declined_by_issuing_bank"
			},
			"provider_data": {
				"response_code": "124"
			},
			"reconciliation_id": "test-reconciliation-id",
			"currency": "RUB"
		}
	}`
	h := http.Header{}
	h.Set("event-type", "payment.void.create")

	sign, err := zooz.CalculateWebhookSignature(context.Background(), []byte(body), h, keyProvider)
	require.NoError(t, err)
	require.Equal(t, "a7e5661bb528f36271142af2c935fdd12865f83919d0a8c5051023ee55339f82", sign)
}

func TestCalculateWebhookSignature_ZeroAmount(t *testing.T) {
	keyProvider := zooz.PrivateKeyProviderFunc(func(_ context.Context, appID string) ([]byte, error) {
		require.Equal(t, "test-app-id", appID)
		return []byte("test-private-key"), nil
	})

	body := `{
		"id": "test-id",
		"created": "2018-10-03T05:14:17.196Z",
		"account_id": "test-account-id",
		"app_id": "test-app-id",
		"payment_id": "test-app-id",
		"data": {
			"id": "test-transaction-id",
			"result": {
				"status": "Pending",
				"category": "payment_method_declined",
				"sub_category": "declined_by_issuing_bank"
			},
			"provider_data": {
				"response_code": "124"
			},
			"reconciliation_id": "test-reconciliation-id",
			"amount": 0,
			"currency": "RUB"
		}
	}`
	h := http.Header{}
	h.Set("event-type", "payment.authorization.create")

	sign, err := zooz.CalculateWebhookSignature(context.Background(), []byte(body), h, keyProvider)
	require.NoError(t, err)
	require.Equal(t, "bb5ce0186b035d21a7dbb23958873082f453acf19b29f7a939cfddf61926acf8", sign)
}

// was validated via zooz sandbox
func TestCalculateWebhookSignature_Payment(t *testing.T) {
	keyProvider := zooz.FixedPrivateKeyProvider{"com.gojuno.gett_development": []byte("test-private-key")}

	body := `{
		"id": "150f1279-6d2c-4728-9e23-6b303f3f1f2f-2020-11-06T14:33:46.150Z-fe14a744-ade7-4a15-8df4-a0b12ac11590",
		"created": "2020-11-06T14:33:46.150Z",
		"account_id": "2058ac1f-ca4b-497f-8181-341e0eea5392",
		"app_id": "com.gojuno.gett_development",
		"payment_id": "150f1279-6d2c-4728-9e23-6b303f3f1f2f",
		"data": {
			"id": "150f1279-6d2c-4728-9e23-6b303f3f1f2f",
			"amount": 100,
			"currency": "RUB"
		}
	}`
	h := http.Header{}
	h.Set("event-type", "payment.payment.create")

	sign, err := zooz.CalculateWebhookSignature(context.Background(), []byte(body), h, keyProvider)
	require.NoError(t, err)
	require.Equal(t, "a74ee728238d42eef7bed9e832eba5b9ab9d95d487027610b8ae50557f09cd05", sign)
}

// was validated via zooz sandbox
func TestCalculateWebhookSignature_Void(t *testing.T) {
	keyProvider := zooz.FixedPrivateKeyProvider{"com.gojuno.gett_development": []byte("test-private-key")}

	body := `{
		"id": "83bc25f0-1710-4ea7-b610-8028f4de6d63-2020-11-06T21:09:22.356Z-fe14a744-ade7-4a15-8df4-a0b12ac11590",
		"created": "2020-11-06T21:09:22.356Z",
		"account_id": "2058ac1f-ca4b-497f-8181-341e0eea5392",
		"app_id": "com.gojuno.gett_development",
		"payment_id": "83bc25f0-1710-4ea7-b610-8028f4de6d63",
		"data": {
			"id": "c5f5ade4-674f-481e-a119-5827a095d352",
			"result": {
				"status": "Succeed"
			},
			"provider_data": {
				"response_code": "0"
			}
		}
	}`
	h := http.Header{}
	h.Set("event-type", "payment.void.create")

	sign, err := zooz.CalculateWebhookSignature(context.Background(), []byte(body), h, keyProvider)
	require.NoError(t, err)
	require.Equal(t, "e2113f4e1728128b18ae945933f8f3e2b01d1e07d2bd8395688109cb18942377", sign)
}

func signature(values []string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	if _, err := mac.Write([]byte(strings.Join(values, ","))); err != nil {
		panic(err)
	}
	return hex.EncodeToString(mac.Sum(nil))
}

func calcRequestSignature(t *testing.T,
	keyProvider interface {
		PrivateKey(ctx context.Context, appID string) ([]byte, error)
	},
	reqBody []byte,
	reqHeader http.Header,
) string {
	signature, err := zooz.CalculateWebhookSignature(context.Background(), reqBody, reqHeader, keyProvider)
	require.NoError(t, err)
	return "sig1=" + signature
}
