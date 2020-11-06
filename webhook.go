package zooz

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const headerEventType = "event-type"

type Callback interface {
	isCallback()
}

type CallbackCommon struct {
	EventType      string `json:"-"` // The type of resource that triggered the event. For example, "payment.authorization.create". Returned in "event-type" header.
	XPaymentsOSEnv string `json:"-"` // PaymentsOS environment, 'live' or 'test'. Returned in "x-payments-os-env" header.
	XZoozRequestID string `json:"-"` // The ID of the original request that triggered the webhook event. Returned in "x-zooz-request-id" header.

	ID        string    `json:"id"`      // The Webhook id. This id is unique per Webhook and can be used to validate that the request is unique.
	Created   time.Time `json:"created"` // The date and time the event was created. "2018-10-03T04:58:35.385Z".
	AccountID string    `json:"account_id"`
	AppID     string    `json:"app_id"`
	PaymentID string    `json:"payment_id"`
}

func (CallbackCommon) isCallback() {}

type PaymentCallback struct {
	CallbackCommon
	Data Payment `json:"data"`
}

type AuthorizationCallback struct {
	CallbackCommon
	Data Authorization `json:"data"`
}

type CaptureCallback struct {
	CallbackCommon
	Data Capture `json:"data"`
}

type VoidCallback struct {
	CallbackCommon
	Data Void `json:"data"`
}

type RefundCallback struct {
	CallbackCommon
	Data Refund `json:"data"`
}

type privateKeyProvider interface {
	// PrivateKey should return private key for given business unit. It is used to validate request signature.
	// For unknown business units (including empty appID) it should return (nil, nil).
	PrivateKey(appID string) ([]byte, error)
}

type PrivateKeyProviderFunc func(appID string) ([]byte, error)

func (f PrivateKeyProviderFunc) PrivateKey(appID string) ([]byte, error) { return f(appID) }

type ErrBadRequest struct {
	Err error
}

func (e ErrBadRequest) Error() string { return e.Err.Error() }

func (e ErrBadRequest) Unwrap() error { return e.Err }

// DecodeWebhookRequest can be used to decode incoming webhook request from PaymentsOS.
// Supports webhook version >= 1.2.0
//
// Will return ErrBadRequest if the error is permanent and request should not be retried.
// Bear in mind that your webhook handler should respond with 2xx status code anyway, otherwise PaymentsOS will continue
// resending this request.
// ErrBadRequest errors include:
// 	* wrong body format (broken/invalid json, ...)
//  * validation error (missing required fields or headers, unexpected values)
// 	* unknown business unit (app_id)
//  * request signature validation error
func DecodeWebhookRequest(_ context.Context, r *http.Request, keyProvider privateKeyProvider) (Callback, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read request body")
	}

	// Validate request signature
	signature, err := calculateWebhookSignature(keyProvider, body, r.Header)
	if err != nil {
		return nil, errors.WithMessage(err, "calculate request signature")
	}
	if "sig1="+signature != r.Header.Get("signature") {
		return nil, ErrBadRequest{errors.New("incorrect signature")}
	}

	// Decode request into appropriate entity type based on "event-type" header
	eventType := r.Header.Get(headerEventType)
	common := CallbackCommon{
		EventType:      eventType,
		XPaymentsOSEnv: r.Header.Get("x-payments-os-env"),
		XZoozRequestID: r.Header.Get("x-zooz-request-id"),
	}
	var cbRef Callback
	switch true {
	case strings.HasPrefix(eventType, "payment.payment."): // @TODO: event-type format is not documented, this is my guess. Validate in zooz sandbox!
		cbRef = &PaymentCallback{CallbackCommon: common}
	case strings.HasPrefix(eventType, "payment.authorization."):
		cbRef = &AuthorizationCallback{CallbackCommon: common}
	case strings.HasPrefix(eventType, "payment.capture."):
		cbRef = &CaptureCallback{CallbackCommon: common}
	case strings.HasPrefix(eventType, "payment.void."):
		cbRef = &VoidCallback{CallbackCommon: common}
	case strings.HasPrefix(eventType, "payment.refund."):
		cbRef = &RefundCallback{CallbackCommon: common}
	default:
		return nil, ErrBadRequest{errors.Errorf("unsupported event type: %q", eventType)}
	}

	if err := json.Unmarshal(body, cbRef); err != nil {
		return nil, errors.Wrapf(err, "unmarshal request body into entity of type %T", cbRef)
	}
	return deref(cbRef), nil
}

// Calculates webhook request signature (without "sig1=" prefix).
// reqBody should be a raw webhook request body, eventType should be the value of "event-type" request header.
// Used by DecodeWebhookRequest.
func calculateWebhookSignature(keyProvider privateKeyProvider, reqBody []byte, reqHeader http.Header) (string, error) {
	type signatureFields struct {
		ID        string `json:"id"`
		Created   string `json:"created"`
		AccountID string `json:"account_id"`
		AppID     string `json:"app_id"`
		PaymentID string `json:"payment_id"`
		Data      struct {
			ID     string `json:"id"`
			Result struct {
				Status      string `json:"status"`
				Category    string `json:"category"`
				SubCategory string `json:"sub_category"`
			} `json:"result"`
			ProviderData struct {
				ResponseCode string `json:"response_code"`
			} `json:"provider_data"`
			ReconciliationID string `json:"reconciliation_id"`
			Amount           *int64 `json:"amount"` // @TODO: should missing amount be zero or an empty string? Test for voids!
			Currency         string `json:"currency"`
		} `json:"data"`
	}

	f := signatureFields{}
	if err := json.Unmarshal(reqBody, &f); err != nil {
		return "", ErrBadRequest{errors.Wrapf(err, "unmarshal request body into entity of type %T", f)}
	}

	key, err := keyProvider.PrivateKey(f.AppID)
	if err != nil {
		return "", errors.Wrap(err, "select private key by app_id")
	}
	if key == nil {
		return "", ErrBadRequest{errors.Errorf("unknown app_id %q", f.AppID)}
	}

	var amount string
	if f.Data.Amount != nil {
		amount = strconv.FormatInt(*f.Data.Amount, 10)
	} else {
		amount = ""
	}
	values := []string{
		reqHeader.Get(headerEventType),
		f.ID,
		f.AccountID,
		f.PaymentID,
		f.Created,
		f.AppID,
		f.Data.ID,
		f.Data.Result.Status,
		f.Data.Result.Category,
		f.Data.Result.SubCategory,
		f.Data.ProviderData.ResponseCode,
		f.Data.ReconciliationID,
		amount,
		f.Data.Currency,
	}
	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write([]byte(strings.Join(values, ","))); err != nil {
		return "", errors.Wrap(err, "sha256 calculation") // should never happen since sha256.digest::Write() never returns error
	}
	sha := hex.EncodeToString(mac.Sum(nil))
	return sha, nil
}

func deref(cb Callback) Callback {
	switch cb := cb.(type) {
	case *PaymentCallback:
		return *cb
	case *AuthorizationCallback:
		return *cb
	case *CaptureCallback:
		return *cb
	case *VoidCallback:
		return *cb
	case *RefundCallback:
		return *cb
	default:
		panic("should be a pointer to a known callback type")
	}
}
