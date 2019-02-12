package zooz

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestSupportedPaymentMethodClient_Get(t *testing.T) {
	spms := []SupportedPaymentMethod{
		{
			ConfigurationID:         "conf-id",
			ConfigurationName:       "conf-name",
			ProviderID:              "provider-id",
			ProviderName:            "provider-name",
			SupportedPaymentMethods: []PaymentMethodData{},
			Result: Result{
				Status: StatusFailed,
			},
		},
	}

	caller := &callerMock{
		t:               t,
		expectedMethod:  http.MethodGet,
		expectedPath:    "supported-payment-methods",
		expectedHeaders: nil,
		expectedReqObj:  nil,
		returnRespObj:   &spms,
	}

	c := &SupportedPaymentMethodClient{Caller: caller}

	res, err := c.Get(context.Background())

	if err != nil {
		t.Error("Error must be nil")
	}
	if res == nil {
		t.Errorf("Supported payment methods is nil")
	}

	if !reflect.DeepEqual(res, spms) {
		t.Errorf("Result is not equal to expected: \n%+v, \n%+v", res, spms)
	}
}
