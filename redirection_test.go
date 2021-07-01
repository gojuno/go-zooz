package zooz

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedirectionClient_Get(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/redirections/id",
		responseBody: `{
			"id": "id"
		}`,
	}

	c := &RedirectionClient{Caller: New(OptHTTPClient(cli))}

	redirection, err := c.Get(
		context.Background(),
		"payment_id",
		"id",
	)

	require.NoError(t, err)
	require.Equal(t, &Redirection{
		ID: "id",
	}, redirection)
}

func TestRedirectionClient_GetList(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/payments/payment_id/redirections",
		responseBody: `[
			{
				"id": "id1"
			},
			{
				"id": "id2"
			}
		]`,
	}

	c := &RedirectionClient{Caller: New(OptHTTPClient(cli))}

	redirections, err := c.GetList(
		context.Background(),
		"payment_id",
	)

	require.NoError(t, err)
	require.Equal(t, []Redirection{
		{
			ID: "id1",
		},
		{
			ID: "id2",
		},
	}, redirections)
}
