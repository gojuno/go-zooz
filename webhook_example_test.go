package zooz_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gtforge/go-zooz"
	"github.com/pkg/errors"
)

func ExampleDecodeWebhookRequest() {
	keyProvider := zooz.FixedPrivateKeyProvider{"app-id": []byte("my-private-key")}

	_ = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		cb, err := zooz.DecodeWebhookRequest(r.Context(), body, r.Header, keyProvider)
		if err != nil {
			if errors.As(err, &zooz.ErrBadRequest{}) {
				// Invalid request. We don't want PaymentsOS to resend it, so respond with 2xx code.
				rw.WriteHeader(200)
			} else {
				// Temporary problem.
				rw.WriteHeader(500)
			}
			return
		}

		switch cb := cb.(type) {
		case zooz.AuthorizationCallback:
			fmt.Printf("authorization: %+v", cb.Data)
		case zooz.CaptureCallback:
			fmt.Printf("capture: %+v", cb.Data)
		}

		rw.WriteHeader(200)
	})
}

func ExampleCalculateWebhookSignature() {
	keyProvider := zooz.FixedPrivateKeyProvider{"app-id": []byte("my-private-key")}

	_ = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		signature, _ := zooz.CalculateWebhookSignature(r.Context(), body, r.Header, keyProvider)
		if "sig1="+signature == r.Header.Get("signature") {
			fmt.Print("Request signature is valid")
		} else {
			fmt.Print("Request signature is invalid")
		}
	})
}
