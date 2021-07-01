package zooz

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentClient_New(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "POST",
		expectedURL:    "/payments",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		expectedBodyJSON: `{
			"amount": 100,
			"currency": "RUB"
		}`,
		responseBody: `{
			"id": "id",
			"amount": 100,
			"currency": "RUB"
		}`,
	}

	c := &PaymentClient{Caller: New(OptHTTPClient(cli))}

	paymentParams := PaymentParams{
		Amount:   100,
		Currency: "RUB",
	}

	payment, err := c.New(
		context.Background(),
		"idempotency_key",
		&paymentParams,
	)

	require.NoError(t, err)
	require.Equal(t, &Payment{
		ID:            "id",
		PaymentParams: paymentParams,
	}, payment)
}

func TestPaymentClient_Get(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/id?expand=authorizations&expand=captures",
		responseBody: `{
			"id": "id"
		}`,
	}

	c := &PaymentClient{Caller: New(OptHTTPClient(cli))}

	payment, err := c.Get(
		context.Background(),
		"id",
		PaymentExpandAuthorizations,
		PaymentExpandCaptures,
	)

	require.NoError(t, err)
	require.Equal(t, &Payment{
		ID: "id",
	}, payment)
}

func TestPaymentClient_Update(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "PUT",
		expectedURL:    "/payments/id",
		expectedBodyJSON: `{
			"amount": 100,
			"currency": "RUB"
		}`,
		responseBody: `{
			"id": "id",
			"amount": 100,
			"currency": "RUB"
		}`,
	}

	c := &PaymentClient{Caller: New(OptHTTPClient(cli))}

	paymentParams := PaymentParams{
		Amount:   100,
		Currency: "RUB",
	}

	payment, err := c.Update(
		context.Background(),
		"id",
		&paymentParams,
	)

	require.NoError(t, err)
	require.Equal(t, &Payment{
		ID:            "id",
		PaymentParams: paymentParams,
	}, payment)
}
