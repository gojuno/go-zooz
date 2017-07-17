package zooz

import (
	"context"
	"testing"
)

func TestCaptureClient_New(t *testing.T) {
	caller := &callerMock{
		t:              t,
		expectedMethod: "POST",
		expectedPath:   "payments/payment_id/captures",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		expectedReqObj: &CaptureParams{
			ReconciliationID: "reconcilation_id",
		},
		returnRespObj: &Capture{
			ID: "id",
		},
	}

	c := &CaptureClient{Caller: caller}

	capture, err := c.New(
		context.Background(),
		"idempotency_key",
		"payment_id",
		&CaptureParams{
			ReconciliationID: "reconcilation_id",
		},
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if capture == nil {
		t.Errorf("Capture is nil")
	}
	if capture.ID != "id" {
		t.Errorf("Capture is not as expected: %+v", capture)
	}
}

func TestCaptureClient_Get(t *testing.T) {
	caller := &callerMock{
		t:               t,
		expectedMethod:  "GET",
		expectedPath:    "payments/payment_id/captures/id",
		expectedHeaders: map[string]string{},
		returnRespObj: &Capture{
			ID: "id",
		},
	}

	c := &CaptureClient{Caller: caller}

	capture, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if capture == nil {
		t.Errorf("Capture is nil")
	}
	if capture.ID != "id" {
		t.Errorf("Capture is not as expected: %+v", capture)
	}
}

func TestCaptureClient_GetList(t *testing.T) {
	caller := &callerMock{
		t:               t,
		expectedMethod:  "GET",
		expectedPath:    "payments/payment_id/captures",
		expectedHeaders: map[string]string{},
		returnRespObj: &[]Capture{
			{
				ID: "id1",
			},
			{
				ID: "id2",
			},
		},
	}

	c := &CaptureClient{Caller: caller}

	captures, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if captures == nil {
		t.Errorf("Captures is nil")
	}
	if len(captures) != 2 {
		t.Errorf("Count of captures is wrong: %d", len(captures))
	}
	if captures[0].ID != "id1" {
		t.Errorf("Capture is not as expected: %+v", captures[0])
	}
	if captures[1].ID != "id2" {
		t.Errorf("Capture is not as expected: %+v", captures[1])
	}
}
