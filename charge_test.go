package zooz

import (
	"context"
	"testing"
)

func TestChargeClient_New(t *testing.T) {
	caller := &callerMock{
		t:              t,
		expectedMethod: "POST",
		expectedPath:   "payments/payment_id/charges",
		expectedHeaders: map[string]string{
			headerIdempotencyKey:  "idempotency_key",
			headerClientIPAddress: "ip",
			headerClientUserAgent: "ua",
		},
		expectedReqObj: &ChargeParams{
			PaymentMethod: PaymentMethodDetails{
				Type:  "tokenized",
				Token: "token",
			},
		},
		returnRespObj: &Charge{
			ID: "id",
		},
	}

	c := &ChargeClient{Caller: caller}

	charge, err := c.New(
		context.Background(),
		"idempotency_key",
		"payment_id",
		&ChargeParams{
			PaymentMethod: PaymentMethodDetails{
				Type:  "tokenized",
				Token: "token",
			},
		},
		&ClientInfo{
			IPAddress: "ip",
			UserAgent: "ua",
		},
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if charge == nil {
		t.Errorf("Charge is nil")
	}
	if charge.ID != "id" {
		t.Errorf("Charge is not as expected: %+v", charge)
	}
}

func TestChargeClient_Get(t *testing.T) {
	caller := &callerMock{
		t:               t,
		expectedMethod:  "GET",
		expectedPath:    "payments/payment_id/charges/id",
		expectedHeaders: map[string]string{},
		returnRespObj: &Charge{
			ID: "id",
		},
	}

	c := &ChargeClient{Caller: caller}

	charge, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if charge == nil {
		t.Errorf("Charge is nil")
	}
	if charge.ID != "id" {
		t.Errorf("Charge is not as expected: %+v", charge)
	}
}

func TestChargeClient_GetList(t *testing.T) {
	caller := &callerMock{
		t:               t,
		expectedMethod:  "GET",
		expectedPath:    "payments/payment_id/charges",
		expectedHeaders: map[string]string{},
		returnRespObj: &[]Charge{
			{
				ID: "id1",
			},
			{
				ID: "id2",
			},
		},
	}

	c := &ChargeClient{Caller: caller}

	charges, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if charges == nil {
		t.Errorf("Charges is nil")
	}
	if len(charges) != 2 {
		t.Errorf("Count of charges is wrong: %d", len(charges))
	}
	if charges[0].ID != "id1" {
		t.Errorf("Charge is not as expected: %+v", charges[0])
	}
	if charges[1].ID != "id2" {
		t.Errorf("Charge is not as expected: %+v", charges[1])
	}
}
