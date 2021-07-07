package zooz

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentMethodClient_New(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "POST",
		expectedURL:    "/customers/customer_id/payment-methods/token1",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		responseBody: `{
			"token": "token1"
		}`,
	}

	c := &PaymentMethodClient{Caller: New(OptHTTPClient(cli))}

	paymentMethod, err := c.New(
		context.Background(),
		"idempotency_key",
		"customer_id",
		"token1",
	)

	require.NoError(t, err)
	require.Equal(t, &PaymentMethod{
		Token: "token1",
	}, paymentMethod)
}

func TestPaymentMethodClient_Get(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/customers/customer_id/payment-methods/token1",
		responseBody: `{
			"token": "token1"
		}`,
	}

	c := &PaymentMethodClient{Caller: New(OptHTTPClient(cli))}

	paymentMethod, err := c.Get(
		context.Background(),
		"customer_id",
		"token1",
	)

	require.NoError(t, err)
	require.Equal(t, &PaymentMethod{
		Token: "token1",
	}, paymentMethod)
}

func TestPaymentMethodClient_GetList(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/customers/customer_id/payment-methods",
		responseBody: `[
			{
				"token": "token1"
			},
			{
				"token": "token2"
			}
		]`,
	}

	c := &PaymentMethodClient{Caller: New(OptHTTPClient(cli))}

	paymentMethods, err := c.GetList(
		context.Background(),
		"customer_id",
	)

	require.NoError(t, err)
	require.Equal(t, []PaymentMethod{
		{
			Token: "token1",
		},
		{
			Token: "token2",
		},
	}, paymentMethods)
}

func TestPaymentMethodClient_Delete(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "DELETE",
		expectedURL:    "/customers/customer_id/payment-methods/token1",
	}

	c := &PaymentMethodClient{Caller: New(OptHTTPClient(cli))}

	err := c.Delete(
		context.Background(),
		"customer_id",
		"token1",
	)

	require.NoError(t, err)
}
