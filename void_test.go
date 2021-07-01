package zooz

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVoidClient_New(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "POST",
		expectedURL:    "/payments/payment_id/voids",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		responseBody: `{
			"id": "id"
		}`,
	}

	c := &VoidClient{Caller: New(OptHTTPClient(cli))}

	void, err := c.New(
		context.Background(),
		"idempotency_key",
		"payment_id",
	)

	require.NoError(t, err)
	require.Equal(t, &Void{
		ID: "id",
	}, void)
}

func TestVoidClient_Get(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/voids/id",
		responseBody: `{
			"id": "id"
		}`,
	}

	c := &VoidClient{Caller: New(OptHTTPClient(cli))}

	void, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	require.NoError(t, err)
	require.Equal(t, &Void{
		ID: "id",
	}, void)
}

func TestVoidClient_GetList(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/voids",
		responseBody: `[
			{
				"id": "id1"
			},
			{
				"id": "id2"
			}
		]`,
	}

	c := &VoidClient{Caller: New(OptHTTPClient(cli))}

	voids, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	require.NoError(t, err)
	require.Equal(t, []Void{
		{
			ID: "id1",
		},
		{
			ID: "id2",
		},
	}, voids)
}
