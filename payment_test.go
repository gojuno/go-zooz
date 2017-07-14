package zooz

import (
	"testing"
	"context"
)

func TestPaymentClient_New(t *testing.T) {
	caller := &callerMock{
		t: t,
		expectedMethod: "POST",
		expectedPath: "payments",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		expectedReqObj: &PaymentParams{
			Amount: 100.0,
		},
		returnRespObj: &Payment{
			ID: "id",
		},
	}

	c := &PaymentClient{Caller: caller}

	payment, err := c.New(
		context.Background(),
		"idempotency_key",
		&PaymentParams{
			Amount: 100.0,
		},
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if payment == nil {
		t.Errorf("Payment is nil")
	}
	if payment.ID != "id" {
		t.Errorf("Payment is not as expected: %+v", payment)
	}
}

func TestPaymentClient_Get(t *testing.T) {
	caller := &callerMock{
		t: t,
		expectedMethod: "GET",
		expectedPath: "payments/id?expand=authorizations&expand=captures",
		returnRespObj: &Payment{
			ID: "id",
		},
	}

	c := &PaymentClient{Caller: caller}

	payment, err := c.Get(
		context.Background(),
		"id",
		PaymentExpandAuthorizations,
		PaymentExpandCaptures,
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if payment == nil {
		t.Errorf("Payment is nil")
	}
	if payment.ID != "id" {
		t.Errorf("Payment is not as expected: %+v", payment)
	}
}

func TestPaymentClient_Update(t *testing.T) {
	caller := &callerMock{
		t: t,
		expectedMethod: "PUT",
		expectedPath: "payments/id",
		expectedReqObj: &PaymentParams{
			Amount: 100.0,
		},
		returnRespObj: &Payment{
			ID: "id",
		},
	}

	c := &PaymentClient{Caller: caller}

	payment, err := c.Update(
		context.Background(),
		"id",
		&PaymentParams{
			Amount: 100.0,
		},
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if payment == nil {
		t.Errorf("Payment is nil")
	}
	if payment.ID != "id" {
		t.Errorf("Payment is not as expected: %+v", payment)
	}
}
