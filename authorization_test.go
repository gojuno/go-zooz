package zooz

import (
	"context"
	"testing"
)

func TestAuthorizationClient_New(t *testing.T) {
	caller := &callerMock{
		t:              t,
		expectedMethod: "POST",
		expectedPath:   "payments/payment_id/authorizations",
		expectedHeaders: map[string]string{
			headerIdempotencyKey:  "idempotency_key",
			headerClientIPAddress: "ip",
			headerClientUserAgent: "ua",
		},
		expectedReqObj: &AuthorizationParams{
			PaymentMethod: PaymentMethodDetails{
				Type:  "tokenized",
				Token: "token",
			},
		},
		returnRespObj: &Authorization{
			ID: "id",
		},
	}

	c := &AuthorizationClient{Caller: caller}

	authorization, err := c.New(
		context.Background(),
		"idempotency_key",
		"payment_id",
		&AuthorizationParams{
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
	if authorization == nil {
		t.Errorf("Authorization is nil")
	}
	if authorization.ID != "id" {
		t.Errorf("Authorization is not as expected: %+v", authorization)
	}
}

func TestAuthorizationClient_Get(t *testing.T) {
	caller := &callerMock{
		t:               t,
		expectedMethod:  "GET",
		expectedPath:    "payments/payment_id/authorizations/id",
		expectedHeaders: map[string]string{},
		returnRespObj: &Authorization{
			ID: "id",
		},
	}

	c := &AuthorizationClient{Caller: caller}

	authorization, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if authorization == nil {
		t.Errorf("Authorization is nil")
	}
	if authorization.ID != "id" {
		t.Errorf("Authorization is not as expected: %+v", authorization)
	}
}

func TestAuthorizationClient_GetList(t *testing.T) {
	caller := &callerMock{
		t:               t,
		expectedMethod:  "GET",
		expectedPath:    "payments/payment_id/authorizations",
		expectedHeaders: map[string]string{},
		returnRespObj: &[]Authorization{
			{
				ID: "id1",
			},
			{
				ID: "id2",
			},
		},
	}

	c := &AuthorizationClient{Caller: caller}

	authorizations, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if authorizations == nil {
		t.Errorf("Authorizations is nil")
	}
	if len(authorizations) != 2 {
		t.Errorf("Count of authorizations is wrong: %d", len(authorizations))
	}
	if authorizations[0].ID != "id1" {
		t.Errorf("Authorization is not as expected: %+v", authorizations[0])
	}
	if authorizations[1].ID != "id2" {
		t.Errorf("Authorization is not as expected: %+v", authorizations[1])
	}
}
