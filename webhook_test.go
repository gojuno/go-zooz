package zooz_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gtforge/go-zooz"
	"github.com/stretchr/testify/require"
)

func TestDecodeWebhookRequest_Payment(t *testing.T) {
	const eventType = "payment.payment.create" // @TODO: validate against real Zooz!
	const privateKey = "test-private-key"

	expected := zooz.PaymentCallback{
		CallbackCommon: zooz.CallbackCommon{
			EventType:      eventType,
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

	keyProvider := zooz.PrivateKeyProviderFunc(func(appID string) ([]byte, error) {
		require.Equal(t, expected.AppID, appID)
		return []byte(privateKey), nil
	})

	body :=
		`{
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
	r := httptest.NewRequest("POST", "http://nevermind", strings.NewReader(body))
	r.Header.Set("event-type", eventType)
	r.Header.Set("x-payments-os-env", expected.XPaymentsOSEnv)
	r.Header.Set("x-zooz-request-id", expected.XZoozRequestID)
	r.Header.Set("signature", zooz.WebhookRequestSignature(t, keyProvider, []byte(body), r.Header))

	cb, err := zooz.DecodeWebhookRequest(context.Background(), r, keyProvider)
	require.NoError(t, err)
	require.Equal(t, expected, cb)
}

func TestDecodeWebhookRequest_Authorization(t *testing.T) {
	const eventType = "payment.authorization.create"
	const privateKey = "test-private-key"

	expected := zooz.AuthorizationCallback{
		CallbackCommon: zooz.CallbackCommon{
			EventType:      eventType,
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

	keyProvider := zooz.PrivateKeyProviderFunc(func(appID string) ([]byte, error) {
		require.Equal(t, expected.AppID, appID)
		return []byte(privateKey), nil
	})

	body :=
		`{
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
	r := httptest.NewRequest("POST", "http://nevermind", strings.NewReader(body))
	r.Header.Set("event-type", eventType)
	r.Header.Set("x-payments-os-env", expected.XPaymentsOSEnv)
	r.Header.Set("x-zooz-request-id", expected.XZoozRequestID)
	r.Header.Set("signature", zooz.WebhookRequestSignature(t, keyProvider, []byte(body), r.Header))

	cb, err := zooz.DecodeWebhookRequest(context.Background(), r, keyProvider)
	require.NoError(t, err)
	require.Equal(t, expected, cb)
}

func TestDecodeWebhookRequest_Capture(t *testing.T) {
	const eventType = "payment.capture.update"
	const privateKey = "test-private-key"

	expected := zooz.CaptureCallback{
		CallbackCommon: zooz.CallbackCommon{
			EventType:      eventType,
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

	keyProvider := zooz.PrivateKeyProviderFunc(func(appID string) ([]byte, error) {
		require.Equal(t, expected.AppID, appID)
		return []byte(privateKey), nil
	})

	body :=
		`{
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
	r := httptest.NewRequest("POST", "http://nevermind", strings.NewReader(body))
	r.Header.Set("event-type", eventType)
	r.Header.Set("x-payments-os-env", expected.XPaymentsOSEnv)
	r.Header.Set("x-zooz-request-id", expected.XZoozRequestID)
	r.Header.Set("signature", zooz.WebhookRequestSignature(t, keyProvider, []byte(body), r.Header))

	cb, err := zooz.DecodeWebhookRequest(context.Background(), r, keyProvider)
	require.NoError(t, err)
	require.Equal(t, expected, cb)
}

func TestDecodeWebhookRequest_Void(t *testing.T) {
	const eventType = "payment.void.create"
	const privateKey = "test-private-key"

	expected := zooz.VoidCallback{
		CallbackCommon: zooz.CallbackCommon{
			EventType:      eventType,
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

	keyProvider := zooz.PrivateKeyProviderFunc(func(appID string) ([]byte, error) {
		require.Equal(t, expected.AppID, appID)
		return []byte(privateKey), nil
	})

	body :=
		`{
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
	r := httptest.NewRequest("POST", "http://nevermind", strings.NewReader(body))
	r.Header.Set("event-type", eventType)
	r.Header.Set("x-payments-os-env", expected.XPaymentsOSEnv)
	r.Header.Set("x-zooz-request-id", expected.XZoozRequestID)
	r.Header.Set("signature", zooz.WebhookRequestSignature(t, keyProvider, []byte(body), r.Header))

	cb, err := zooz.DecodeWebhookRequest(context.Background(), r, keyProvider)
	require.NoError(t, err)
	require.Equal(t, expected, cb)
}

func TestDecodeWebhookRequest_Refund(t *testing.T) {
	const eventType = "payment.refund.update"
	const appID = "test-app-id"
	const privateKey = "test-private-key"

	expected := zooz.RefundCallback{
		CallbackCommon: zooz.CallbackCommon{
			EventType:      eventType,
			XPaymentsOSEnv: "live",
			XZoozRequestID: "test-x-zooz-request-id",

			ID:        "test-webhook-id",
			Created:   time.Date(2018, 10, 03, 05, 22, 45, 610000000, time.UTC), // "2018-10-03T05:22:45.610Z"
			AccountID: "test-account-id",
			AppID:     appID,
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

	keyProvider := zooz.PrivateKeyProviderFunc(func(appID string) ([]byte, error) {
		require.Equal(t, "test-app-id", appID)
		return []byte(privateKey), nil
	})

	body :=
		`{
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
	r := httptest.NewRequest("POST", "http://nevermind", strings.NewReader(body))
	r.Header.Set("event-type", eventType)
	r.Header.Set("x-payments-os-env", expected.XPaymentsOSEnv)
	r.Header.Set("x-zooz-request-id", expected.XZoozRequestID)
	r.Header.Set("signature", zooz.WebhookRequestSignature(t, keyProvider, []byte(body), r.Header))

	cb, err := zooz.DecodeWebhookRequest(context.Background(), r, keyProvider)
	require.NoError(t, err)
	require.Equal(t, expected, cb)
}

func TestSignature_AllFields(t *testing.T) {
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
		data_Result_Category           = "payment_method_declined"
		data_Result_SubCategory        = "declined_by_issuing_bank"
		data_ProviderData_ResponseCode = "124"
		data_ReconciliationID          = "test-reconciliation-id"
		data_Amount                    = "42"
		data_Currency                  = "RUB"
	)

	keyProvider := zooz.PrivateKeyProviderFunc(func(id string) ([]byte, error) {
		require.Equal(t, appID, id)
		return []byte(privateKey), nil
	})

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

	sign, err := zooz.CalculateWebhookSignature(keyProvider, []byte(body), h)
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
			data_Result_Category,
			data_Result_SubCategory,
			data_ProviderData_ResponseCode,
			data_ReconciliationID,
			data_Amount,
			data_Currency,
		}, privateKey),
		sign)
	require.Equal(t, "3d498ba3149ebd503164bbc7a48feeaea74acc1e5627df5e7dad58a339f7d4af", sign)
}

func TestSignature_NoAmount(t *testing.T) {
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
		data_Result_Category           = "payment_method_declined"
		data_Result_SubCategory        = "declined_by_issuing_bank"
		data_ProviderData_ResponseCode = "124"
		data_ReconciliationID          = "test-reconciliation-id"
		data_Currency                  = "RUB"
	)

	keyProvider := zooz.PrivateKeyProviderFunc(func(id string) ([]byte, error) {
		require.Equal(t, appID, id)
		return []byte(privateKey), nil
	})

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
		"currency": "` + data_Currency + `"
	}
}`
	h := http.Header{}
	h.Set("event-type", eventType)

	sign, err := zooz.CalculateWebhookSignature(keyProvider, []byte(body), h)
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
			data_Result_Category,
			data_Result_SubCategory,
			data_ProviderData_ResponseCode,
			data_ReconciliationID,
			"", // no amount
			data_Currency,
		}, privateKey),
		sign)
	require.Equal(t, "330ea79c4154892ec9a57edff161c60a2421bed4872060303c925f41265e0d16", sign)
}

func TestSignature_ZeroAmount(t *testing.T) {
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
		data_Result_Category           = "payment_method_declined"
		data_Result_SubCategory        = "declined_by_issuing_bank"
		data_ProviderData_ResponseCode = "124"
		data_ReconciliationID          = "test-reconciliation-id"
		data_Currency                  = "RUB"
	)

	keyProvider := zooz.PrivateKeyProviderFunc(func(id string) ([]byte, error) {
		require.Equal(t, appID, id)
		return []byte(privateKey), nil
	})

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
		"amount": 0,
		"currency": "` + data_Currency + `"
	}
}`
	h := http.Header{}
	h.Set("event-type", eventType)

	sign, err := zooz.CalculateWebhookSignature(keyProvider, []byte(body), h)
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
			data_Result_Category,
			data_Result_SubCategory,
			data_ProviderData_ResponseCode,
			data_ReconciliationID,
			"0",
			data_Currency,
		}, privateKey),
		sign)
	require.Equal(t, "bb5ce0186b035d21a7dbb23958873082f453acf19b29f7a939cfddf61926acf8", sign)
}

func TestSignature_MissingFields(t *testing.T) {
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

	keyProvider := zooz.PrivateKeyProviderFunc(func(id string) ([]byte, error) {
		require.Equal(t, appID, id)
		return []byte(privateKey), nil
	})

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

	sign, err := zooz.CalculateWebhookSignature(keyProvider, []byte(body), h)
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

//
//func TestSignature_POS(t *testing.T) {
//	const (
//		privateKey = "123456"
//		expectedSignature = "5f024b177e670949c1c5efc214987467c784aef50807758d22b9d287c8f06ed7"
//
//		eventType                      = "payment.authorization.create"
//		id                             = "ID"
//		created                        = "2018-10-03T04:58:35.385Z"
//		accountID                      = "accountID"
//		appID                          = "appID"
//		paymentID                      = "paymentID"
//		data_ID                        = "operationID"
//		data_Result_Status             = "Succeed"
//		data_Result_Category           = "category"
//		data_Result_SubCategory        = "subCategory"
//		data_ProviderData_ResponseCode = "responseCode"
//		data_ReconciliationID = "reconciliationID"
//		data_Amount = "1000"
//		data_Currency = "RUB"
//	)
//
//	keyProvider := zooz.PrivateKeyProviderFunc(func(id string) ([]byte, error) {
//		require.Equal(t, appID, id)
//		return []byte(privateKey), nil
//	})
//
//	body := `{
//	"id": "` + id + `",
//    "created": "` + created + `",
//    "account_id": "` + accountID + `",
//    "app_id": "` + appID + `",
//    "payment_id": "` + paymentID + `",
//	"data": {
//		"id": "` + data_ID + `",
//  		"result": {
//			"status": "` + data_Result_Status + `",
//    		"category": "` + data_Result_Category + `",
//    		"sub_category": "` + data_Result_SubCategory + `"
//  		},
//  		"provider_data": {
//    		"response_code": "` + data_ProviderData_ResponseCode + `"
//  		},
//		"reconciliation_id": "` + data_ReconciliationID +`",
//		"amount": ` + data_Amount +`,
//		"currency": "` + data_Currency +`"
//	}
//}`
//	h := http.Header{}
//	h.Set("event-type", eventType)
//
//	sign, err := zooz.CalculateWebhookSignature(keyProvider, []byte(body), h)
//	require.NoError(t, err)
//	require.Equal(t,
//		signature([]string{
//			eventType,
//			id,
//			accountID,
//			paymentID,
//			created,
//			appID,
//			data_ID,
//			data_Result_Status,
//			data_Result_Category,
//			data_Result_SubCategory,
//			data_ProviderData_ResponseCode,
//			data_ReconciliationID,
//			data_Amount,
//			data_Currency,
//		}, privateKey),
//		sign)
//	require.Equal(t, expectedSignature, sign)
//}

func signature(values []string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	if _, err := mac.Write([]byte(strings.Join(values, ","))); err != nil {
		panic(err)
	}
	return hex.EncodeToString(mac.Sum(nil))
}
