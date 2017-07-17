package zooz

import (
	"context"
	"testing"
)

func TestRefundClient_New(t *testing.T) {
	caller := &callerMock{
		t:              t,
		expectedMethod: "POST",
		expectedPath:   "payments/payment_id/refunds",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		expectedReqObj: &RefundParams{
			ReconciliationID: "reconcilation_id",
		},
		returnRespObj: &Refund{
			ID: "id",
		},
	}

	c := &RefundClient{Caller: caller}

	refund, err := c.New(
		context.Background(),
		"idempotency_key",
		"payment_id",
		&RefundParams{
			ReconciliationID: "reconcilation_id",
		},
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if refund == nil {
		t.Errorf("Refund is nil")
	}
	if refund.ID != "id" {
		t.Errorf("Refund is not as expected: %+v", refund)
	}
}

func TestRefundClient_Get(t *testing.T) {
	caller := &callerMock{
		t:               t,
		expectedMethod:  "GET",
		expectedPath:    "payments/payment_id/refunds/id",
		expectedHeaders: map[string]string{},
		returnRespObj: &Refund{
			ID: "id",
		},
	}

	c := &RefundClient{Caller: caller}

	refund, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if refund == nil {
		t.Errorf("Refund is nil")
	}
	if refund.ID != "id" {
		t.Errorf("Refund is not as expected: %+v", refund)
	}
}

func TestRefundClient_GetList(t *testing.T) {
	caller := &callerMock{
		t:               t,
		expectedMethod:  "GET",
		expectedPath:    "payments/payment_id/refunds",
		expectedHeaders: map[string]string{},
		returnRespObj: &[]Refund{
			{
				ID: "id1",
			},
			{
				ID: "id2",
			},
		},
	}

	c := &RefundClient{Caller: caller}

	refunds, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if refunds == nil {
		t.Errorf("Refunds is nil")
	}
	if len(refunds) != 2 {
		t.Errorf("Count of refunds is wrong: %d", len(refunds))
	}
	if refunds[0].ID != "id1" {
		t.Errorf("Refund is not as expected: %+v", refunds[0])
	}
	if refunds[1].ID != "id2" {
		t.Errorf("Refund is not as expected: %+v", refunds[1])
	}
}
