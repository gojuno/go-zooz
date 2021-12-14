package zooz

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendRequest(t *testing.T) {
	z := New(
		OptAppID("com.gojuno.gett_stage"),
		OptPrivateKey("fab6f75d-c706-4259-834b-c78e53e299bc"),
		OptEnv(EnvTest),
	)

	paymentParams := &PaymentParams{
		Amount:   100,
		Currency: "RUB",
	}

	paym, err := z.Payment().New(context.Background(), "idempotency_id", paymentParams)

	if err != nil {
		log.Println(err)
	}

	auth, err := z.Authorization().New(context.Background(),
		"idempotency_key",
		paym.ID,
		&AuthorizationParams{
			PaymentMethod: PaymentMethodDetails{
				Type:  "tokenized",
				Token: "bcd4e9b1-4577-43d3-9280-3c54c5356941",
			},
		},
		&ClientInfo{
			IPAddress: "ip",
			UserAgent: "ua",
		})
	if err != nil {
		log.Println(err)
	}

	log.Println(auth)
}

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
	}

	c := &AuthorizationClient{Caller: New(OptHTTPClient(cli))}

	err := c.ContinueAuthentication(
		context.Background(),
		"idempotency_key",
		"payment_id",
		"authorization_id",
		&ContinueAuthenticationParams{
			ReconciliationID: "reconciliation_id",
			ThreeDSecureAttributes: &ThreeDSecureAttributes{ThreeDSecureAttributesInternal{
				DataCollectionCompletedInd: AuthenticationDataCollectionCompleted,
			}},
		},
		&ClientInfo{
			IPAddress: "ip",
			UserAgent: "ua",
		},
	)

	require.NoError(t, err)
}
