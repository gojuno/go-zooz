package zooz

import (
	"context"
	"testing"
)

func TestPaymentMethodClient_New(t *testing.T) {
	caller := &callerMock{
		t: t,
		expectedMethod: "POST",
		expectedPath: "customers/customer_id/payment-methods/token",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		returnRespObj: &PaymentMethod{
			Token: "token",
		},
	}

	c := &PaymentMethodClient{Caller: caller}

	paymentMethod, err := c.New(
		context.Background(),
		"idempotency_key",
		"customer_id",
		"token",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if paymentMethod == nil {
		t.Errorf("PaymentMethod is nil")
	}
	if paymentMethod.Token != "token" {
		t.Errorf("PaymentMethod is not as expected: %+v", paymentMethod)
	}
}

func TestPaymentMethodClient_Get(t *testing.T) {
	caller := &callerMock{
		t: t,
		expectedMethod: "GET",
		expectedPath: "customers/customer_id/payment-methods/token",
		expectedHeaders: map[string]string{},
		returnRespObj: &PaymentMethod{
			Token: "token",
		},
	}

	c := &PaymentMethodClient{Caller: caller}

	paymentMethod, err := c.Get(
		context.Background(),
		"customer_id",
		"token",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if paymentMethod == nil {
		t.Errorf("PaymentMethod is nil")
	}
	if paymentMethod.Token != "token" {
		t.Errorf("PaymentMethod is not as expected: %+v", paymentMethod)
	}
}

func TestPaymentMethodClient_GetList(t *testing.T) {
	caller := &callerMock{
		t: t,
		expectedMethod: "GET",
		expectedPath: "customers/customer_id/payment-methods",
		expectedHeaders: map[string]string{},
		returnRespObj: &[]PaymentMethod{
			{
				Token: "token1",
			},
			{
				Token: "token2",
			},
		},
	}

	c := &PaymentMethodClient{Caller: caller}

	paymentMethods, err := c.GetList(
		context.Background(),
		"customer_id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if paymentMethods == nil {
		t.Errorf("PaymentMethods is nil")
	}
	if len(paymentMethods) != 2 {
		t.Errorf("Count of paymentMethods is wrong: %d", len(paymentMethods))
	}
	if paymentMethods[0].Token != "token1" {
		t.Errorf("PaymentMethod is not as expected: %+v", paymentMethods[0])
	}
	if paymentMethods[1].Token != "token2" {
		t.Errorf("PaymentMethod is not as expected: %+v", paymentMethods[1])
	}
}