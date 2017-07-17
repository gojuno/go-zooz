package zooz

import (
	"context"
	"testing"
)

func TestRedirectionClient_Get(t *testing.T) {
	caller := &callerMock{
		t:               t,
		expectedMethod:  "GET",
		expectedPath:    "payments/payment_id/redirections/id",
		expectedHeaders: map[string]string{},
		returnRespObj: &Redirection{
			ID: "id",
		},
	}

	c := &RedirectionClient{Caller: caller}

	redirection, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if redirection == nil {
		t.Errorf("Redirection is nil")
	}
	if redirection.ID != "id" {
		t.Errorf("Redirection is not as expected: %+v", redirection)
	}
}

func TestRedirectionClient_GetList(t *testing.T) {
	caller := &callerMock{
		t:               t,
		expectedMethod:  "GET",
		expectedPath:    "payments/payment_id/redirections",
		expectedHeaders: map[string]string{},
		returnRespObj: &[]Redirection{
			{
				ID: "id1",
			},
			{
				ID: "id2",
			},
		},
	}

	c := &RedirectionClient{Caller: caller}

	redirections, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	if err != nil {
		t.Error("Error must be nil")
	}
	if redirections == nil {
		t.Errorf("Redirections is nil")
	}
	if len(redirections) != 2 {
		t.Errorf("Count of redirections is wrong: %d", len(redirections))
	}
	if redirections[0].ID != "id1" {
		t.Errorf("Redirection is not as expected: %+v", redirections[0])
	}
	if redirections[1].ID != "id2" {
		t.Errorf("Redirection is not as expected: %+v", redirections[1])
	}
}
