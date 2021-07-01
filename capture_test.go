package zooz

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCaptureClient_New(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "POST",
		expectedURL:    "/payments/payment_id/captures",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		expectedBodyJSON: `{
			"reconciliation_id": "reconciliation_id"
		}`,
		responseBody: `{
			"id": "id"
		}`,
	}

	c := &CaptureClient{Caller: New(OptHTTPClient(cli))}

	capture, err := c.New(
		context.Background(),
		"idempotency_key",
		"payment_id",
		&CaptureParams{
			ReconciliationID: "reconciliation_id",
		},
	)

	require.NoError(t, err)
	require.Equal(t, &Capture{
		ID: "id",
	}, capture)
}

func TestCaptureClient_Get(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/captures/id",
		responseBody: `{
			"id": "id"
		}`,
	}

	c := &CaptureClient{Caller: New(OptHTTPClient(cli))}

	capture, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	require.NoError(t, err)
	require.Equal(t, &Capture{
		ID: "id",
	}, capture)
}

func TestCaptureClient_GetList(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/captures",
		responseBody: `[
			{
				"id": "id1"
			},
			{
				"id": "id2"
			}
		]`,
	}

	c := &CaptureClient{Caller: New(OptHTTPClient(cli))}

	captures, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	require.NoError(t, err)
	require.Equal(t, []Capture{
		{
			ID: "id1",
		},
		{
			ID: "id2",
		},
	}, captures)
}
