package zooz

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefundClient_New(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "POST",
		expectedURL:    "/payments/payment_id/refunds",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		expectedBodyJSON: `{
			"reconciliation_id": "reconciliation_id",
			"amount": 42
		}`,
		responseBody: `{
			"id": "id",
			"reconciliation_id": "reconciliation_id",
			"amount": 42
		}`,
	}

	c := &RefundClient{Caller: New(OptHTTPClient(cli))}

	refundParams := RefundParams{
		ReconciliationID: "reconciliation_id",
		Amount:           42,
	}

	refund, err := c.New(
		context.Background(),
		"idempotency_key",
		"payment_id",
		&refundParams,
	)

	require.NoError(t, err)
	require.Equal(t, &Refund{
		ID:           "id",
		RefundParams: refundParams,
	}, refund)
}

func TestRefundClient_Get(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/refunds/id",
		responseBody: `{
			"id": "id",
			"reconciliation_id": "reconciliation_id",
			"amount": 42
		}`,
	}

	c := &RefundClient{Caller: New(OptHTTPClient(cli))}

	refund, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	require.NoError(t, err)
	require.Equal(t, &Refund{
		ID: "id",
		RefundParams: RefundParams{
			ReconciliationID: "reconciliation_id",
			Amount:           42,
		},
	}, refund)
}

func TestRefundClient_GetList(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/refunds",
		responseBody: `[
			{
				"id": "id1"
			},
			{
				"id": "id2"
			}
		]`,
	}

	c := &RefundClient{Caller: New(OptHTTPClient(cli))}

	refunds, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	require.NoError(t, err)
	require.Equal(t, []Refund{
		{
			ID: "id1",
		},
		{
			ID: "id2",
		},
	}, refunds)
}
