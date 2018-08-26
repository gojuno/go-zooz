package zooz

import (
	"context"
	"testing"
)

func TestCustomerClient_New(t *testing.T) {
	caller := &callerMock{
		t:              t,
		expectedMethod: "POST",
		expectedPath:   "customers",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		expectedReqObj: &CustomerParams{
			CustomerReference: "reference",
		},
		returnRespObj: &Customer{
			ID: "id",
		},
	}

	c := &CustomerClient{Caller: caller}

	customer, err := c.New(
		context.Background(),
		"idempotency_key",
		&CustomerParams{
			CustomerReference: "reference",
		},
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if customer == nil {
		t.Errorf("Customer is nil")
	}
	if customer.ID != "id" {
		t.Errorf("Customer is not as expected: %+v", customer)
	}
}

func TestCustomerClient_Get(t *testing.T) {
	caller := &callerMock{
		t:              t,
		expectedMethod: "GET",
		expectedPath:   "customers/id",
		returnRespObj: &Customer{
			ID: "id",
		},
	}

	c := &CustomerClient{Caller: caller}

	customer, err := c.Get(
		context.Background(),
		"id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if customer == nil {
		t.Errorf("Customer is nil")
	}
	if customer.ID != "id" {
		t.Errorf("Customer is not as expected: %+v", customer)
	}
}

func TestCustomerClient_Update(t *testing.T) {
	caller := &callerMock{
		t:              t,
		expectedMethod: "PUT",
		expectedPath:   "customers/id",
		expectedReqObj: &CustomerParams{
			CustomerReference: "reference",
		},
		returnRespObj: &Customer{
			ID: "id",
		},
	}

	c := &CustomerClient{Caller: caller}

	customer, err := c.Update(
		context.Background(),
		"id",
		&CustomerParams{
			CustomerReference: "reference",
		},
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if customer == nil {
		t.Errorf("Customer is nil")
	}
	if customer.ID != "id" {
		t.Errorf("Customer is not as expected: %+v", customer)
	}
}

func TestCustomerClient_Delete(t *testing.T) {
	caller := &callerMock{
		t:              t,
		expectedMethod: "DELETE",
		expectedPath:   "customers/id",
	}

	c := &CustomerClient{Caller: caller}

	err := c.Delete(
		context.Background(),
		"id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
}

func TestCustomerClient_GetByReference(t *testing.T) {
	caller := &callerMock{
		t:              t,
		expectedMethod: "GET",
		expectedPath:   "customers/customer_reference",
		returnRespObj: &Customer{
			ID: "id",
			CustomerParams: CustomerParams{
				CustomerReference: "customer_reference",
			},
		},
	}

	c := &CustomerClient{Caller: caller}

	customer, err := c.Get(
		context.Background(),
		"customer_reference",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if customer == nil {
		t.Errorf("Customer is nil")
	}
	if customer.CustomerReference != "customer_reference" {
		t.Errorf("Customer is not as expected: %+v", customer)
	}
}
