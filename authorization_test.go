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

func TestAuthorizationClient_New_3DS_Internal(t *testing.T) {
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
			},
			"three_d_secure_attributes": {
				"internal": {
					"device_channel": "02",
					"browser_header": "text/html",
					"challenge_window_size": "05"
				}
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
			ThreeDSecureAttributes: &ThreeDSecureAttributes{
				Internal: &ThreeDSecureAttributesInternal{
					DeviceChannel:       "02",
					BrowserHeader:       "text/html",
					ChallengeWindowSize: "05",
				}},
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

func TestAuthorizationClient_ContinueAuthentication(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "POST",
		expectedURL:    "/payments/payment_id/authorizations/authorization_id/authentication-flows",
		expectedHeaders: map[string]string{
			headerIdempotencyKey:  "idempotency_key",
			headerClientIPAddress: "ip",
			headerClientUserAgent: "ua",
		},
		expectedBodyJSON: `{
			"reconciliation_id": "reconciliation_id",
			"three_d_secure_attributes" : {
				"internal": {
					"data_collection_completed_ind": "Y"
				}
			}
		}`,
		responseBody: `{
			"related_resources": {
				"authorizations": [
					{
						"id": "authorization_id"
					},
					{
						"id": "not_valid"
					}
				]
			}
		}`,
	}

	c := &AuthorizationClient{Caller: New(OptHTTPClient(cli))}

	auth, err := c.ContinueAuthentication(
		context.Background(),
		"idempotency_key",
		ContinueAuthenticationParams{
			PaymentID:                  "payment_id",
			AuthorizationID:            "authorization_id",
			ReconciliationID:           "reconciliation_id",
			DataCollectionCompletedInd: AuthenticationDataCollectionCompleted,
		},
		&ClientInfo{
			IPAddress: "ip",
			UserAgent: "ua",
		},
	)

	require.Equal(t, "authorization_id", auth.ID)
	require.NoError(t, err)
}
