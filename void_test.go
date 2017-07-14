package zooz

import (
	"context"
	"testing"
)

func TestVoidClient_New(t *testing.T) {
	caller := &callerMock{
		t: t,
		expectedMethod: "POST",
		expectedPath: "payments/payment_id/voids",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		returnRespObj: &Void{
			ID: "id",
		},
	}

	c := &VoidClient{Caller: caller}

	void, err := c.New(
		context.Background(),
		"idempotency_key",
		"payment_id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if void == nil {
		t.Errorf("Void is nil")
	}
	if void.ID != "id" {
		t.Errorf("Void is not as expected: %+v", void)
	}
}

func TestVoidClient_Get(t *testing.T) {
	caller := &callerMock{
		t: t,
		expectedMethod: "GET",
		expectedPath: "payments/payment_id/voids/id",
		expectedHeaders: map[string]string{},
		returnRespObj: &Void{
			ID: "id",
		},
	}

	c := &VoidClient{Caller: caller}

	void, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if void == nil {
		t.Errorf("Void is nil")
	}
	if void.ID != "id" {
		t.Errorf("Void is not as expected: %+v", void)
	}
}

func TestVoidClient_GetList(t *testing.T) {
	caller := &callerMock{
		t: t,
		expectedMethod: "GET",
		expectedPath: "payments/payment_id/voids",
		expectedHeaders: map[string]string{},
		returnRespObj: &[]Void{
			{
				ID: "id1",
			},
			{
				ID: "id2",
			},
		},
	}

	c := &VoidClient{Caller: caller}

	voids, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if voids == nil {
		t.Errorf("Voids is nil")
	}
	if len(voids) != 2 {
		t.Errorf("Count of voids is wrong: %d", len(voids))
	}
	if voids[0].ID != "id1" {
		t.Errorf("Void is not as expected: %+v", voids[0])
	}
	if voids[1].ID != "id2" {
		t.Errorf("Void is not as expected: %+v", voids[1])
	}
}
