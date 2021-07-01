package zooz

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChargeClient_New(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "POST",
		expectedURL:    "/payments/payment_id/charges",
		expectedHeaders: map[string]string{
			headerIdempotencyKey:  "idempotency_key",
			headerClientIPAddress: "ip",
			headerClientUserAgent: "ua",
		},
		expectedBodyJSON: `{
			"payment_method": {
				"type": "tokenized",
				"token": "token"
			}
		}`,
		responseBody: `{
			"id": "id",
			"payment_method": {
				"type": "tokenized",
				"token": "token"
			}
		}`,
	}

	c := &ChargeClient{Caller: New(OptHTTPClient(cli))}

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

	require.NoError(t, err)
	require.Equal(t, &Charge{
		ID: "id",
		PaymentMethod: PaymentMethod{
			Type:  "tokenized",
			Token: "token",
		},
	}, charge)
}

func TestChargeClient_Get(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/charges/id",
		responseBody: `{
			"id": "id",
			"payment_method": {
				"type": "tokenized",
				"token": "token"
			}
		}`,
	}

	c := &ChargeClient{Caller: New(OptHTTPClient(cli))}

	charge, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	require.NoError(t, err)
	require.Equal(t, &Charge{
		ID: "id",
		PaymentMethod: PaymentMethod{
			Type:  "tokenized",
			Token: "token",
		},
	}, charge)
}

func TestChargeClient_GetList(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/charges",
		responseBody: `[
			{
				"id": "id1"
			},
			{
				"id": "id2"
			}
		]`,
	}

	c := &ChargeClient{Caller: New(OptHTTPClient(cli))}

	charges, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	require.NoError(t, err)
	require.Equal(t, []Charge{
		{
			ID: "id1",
		},
		{
			ID: "id2",
		},
	}, charges)
}
