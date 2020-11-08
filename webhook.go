package zooz

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

func (PaymentCallback) isCallback()       {}
func (AuthorizationCallback) isCallback() {}
func (CaptureCallback) isCallback()       {}
func (VoidCallback) isCallback()          {}
func (RefundCallback) isCallback()        {}

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
	PrivateKey(ctx context.Context, appID string) ([]byte, error)
}

type PrivateKeyProviderFunc func(ctx context.Context, appID string) ([]byte, error)

func (f PrivateKeyProviderFunc) PrivateKey(ctx context.Context, appID string) ([]byte, error) {
	return f(ctx, appID)
}

type FixedPrivateKeyProvider map[string][]byte

func (m FixedPrivateKeyProvider) PrivateKey(_ context.Context, appID string) ([]byte, error) {
	return m[appID], nil
}

type ErrBadRequest struct {
	Err error
}

func (e ErrBadRequest) Error() string { return e.Err.Error() }

func (e ErrBadRequest) Unwrap() error { return e.Err }

// DecodeWebhookRequest decodes PaymentsOS webhook http request into callback entity.
// Supports webhook version >= 1.2.0
// reqBody should be a raw webhook http request body.
// reqHeader should be webhook http request headers.
//
// Use type switch to determine concrete callback entity.
// Returned callback entity is a value, never a pointer (i.e. AuthorizationCallback, not *AuthorizationCallback).
//
//  cb, _ := DecodeWebhookRequest(...)
//  switch cb := cb.(type) {
//    case zooz.AuthorizationCallback:
//	  case zooz.CaptureCallback:
//	  ...
//  }
//
// Will return ErrBadRequest if the error is permanent and request should not be retried.
// Bear in mind that your webhook handler should respond with 2xx status code anyway, otherwise PaymentsOS will continue
// resending this request.
// ErrBadRequest errors include:
//  * wrong body format (broken/invalid json, ...)
//  * validation error (missing required fields or headers, unexpected values)
//  * unknown business unit (app_id)
//  * incorrect request signature
func DecodeWebhookRequest(ctx context.Context, body []byte, header http.Header, keyProvider privateKeyProvider) (Callback, error) {
	// Validate request signature
	signature, err := CalculateWebhookSignature(ctx, body, header, keyProvider)
	if err != nil {
		return nil, errors.WithMessage(err, "calculate request signature")
	}
	if "sig1="+signature != header.Get("signature") {
		return nil, ErrBadRequest{errors.New("incorrect signature")}
	}

	// Decode request into appropriate entity type based on "event-type" header
	eventType := header.Get(headerEventType)
	common := CallbackCommon{
		EventType:      eventType,
		XPaymentsOSEnv: header.Get("x-payments-os-env"),
		XZoozRequestID: header.Get("x-zooz-request-id"),
	}
	var cbRef Callback
	switch {
	case strings.HasPrefix(eventType, "payment.payment."):
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

// CalculateWebhookSignature calculates signature for webhook request (without "sig1=" prefix).
// Generally, you would compare it to "signature" request header sent by PaymentsOS:
//
//  signature, _ := CalculateWebhookSignature(...)
//  if "sig1=" + signature == r.Header.Get("signature") {
//    // webhook request is valid
//  }
//
// reqBody should be a raw webhook http request body.
// reqHeader should be webhook http request headers (currently only 'event-type' header matters).
// Will return ErrBadRequest{} if error is permanent and should not be retried (malformed json body, unknown app_id)
// Also used by DecodeWebhookRequest.
func CalculateWebhookSignature(ctx context.Context, reqBody []byte, reqHeader http.Header, keyProvider privateKeyProvider) (string, error) {
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
			Amount           *int64 `json:"amount"`
			Currency         string `json:"currency"`
		} `json:"data"`
	}

	f := signatureFields{}
	if err := json.Unmarshal(reqBody, &f); err != nil {
		return "", ErrBadRequest{errors.Wrapf(err, "unmarshal request body into entity of type %T", f)}
	}

	key, err := keyProvider.PrivateKey(ctx, f.AppID)
	if err != nil {
		return "", errors.Wrap(err, "select private key by app_id")
	}
	if len(key) == 0 {
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
