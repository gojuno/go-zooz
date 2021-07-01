package zooz

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthorizationClient_New(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "POST",
		expectedURL:    "/payments/payment_id/authorizations",
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

	c := &AuthorizationClient{Caller: New(OptHTTPClient(cli))}

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

	require.NoError(t, err)
	require.Equal(t, &Authorization{
		ID: "id",
		PaymentMethod: PaymentMethod{
			Type:  "tokenized",
			Token: "token",
		},
	}, authorization)
}

func TestAuthorizationClient_Get(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/authorizations/id",
		responseBody: `{
			"id": "id"
		}`,
	}

	c := &AuthorizationClient{Caller: New(OptHTTPClient(cli))}

	authorization, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	require.NoError(t, err)
	require.Equal(t, &Authorization{
		ID: "id",
	}, authorization)
}

func TestAuthorizationClient_GetList(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/authorizations",
		responseBody: `[
			{
				"id": "id1"
			},
			{
				"id": "id2"
			}
		]`,
	}

	c := &AuthorizationClient{Caller: New(OptHTTPClient(cli))}

	authorizations, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	require.NoError(t, err)
	require.Equal(t, []Authorization{
		{
			ID: "id1",
		},
		{
			ID: "id2",
		},
	}, authorizations)
}
