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
			"token": "token1",
			"expiration_date": "12/2051"
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
		Token:          "token1",
		ExpirationDate: "12/2051",
	}, paymentMethod)
}

func TestPaymentMethodClient_Get(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/customers/customer_id/payment-methods/token1",
		responseBody: `{
			"token": "token1",
			"expiration_date": "12/2051"
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
		Token:          "token1",
		ExpirationDate: "12/2051",
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

func Test_ParseCardExpirationDate(t *testing.T) {
	ttc := []struct {
		d   string
		m   int
		y   int
		err string
	}{
		{
			d: "01/1994",
			m: 1, y: 1994,
		},
		{
			d: "01.1984",
			m: 1, y: 1984,
		},
		{
			d: "01-1984",
			m: 1, y: 1984,
		},
		{
			d: "01 1984",
			m: 1, y: 1984,
		},
		{
			d: "01/0000",
			m: 1, y: 0,
		},
		{
			d: "01/01",
			m: 1, y: 2001,
		},
		{
			d: "01.00",
			m: 1, y: 2000,
		},
		{
			d: "01 99",
			m: 1, y: 2099,
		},
		{
			d:   "13/9999",
			err: "month value out of range: 13",
		},
		{
			d:   "00/1984",
			err: "month value out of range: 0",
		},
		{
			d:   "invalid",
			err: "unexpected expiration date format",
		},
		{
			d:   "1/2021",
			err: "unexpected expiration date format",
		},
		{
			d:   "01/200",
			err: "unexpected expiration date format",
		},
	}

	for _, tc := range ttc {
		tc := tc
		t.Run(tc.d, func(t *testing.T) {
			m, y, err := ParseCardExpirationDate(tc.d)
			if tc.err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err)
			} else {
				require.Equal(t, tc.m, m)
				require.Equal(t, tc.y, y)
			}
		})
	}
}
