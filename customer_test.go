package zooz

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustomerClient_New(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "POST",
		expectedURL:    "/customers",
		expectedHeaders: map[string]string{
			headerIdempotencyKey: "idempotency_key",
		},
		expectedBodyJSON: `{
			"customer_reference": "reference",
			"first_name": "john",
			"last_name": "doe",
			"email": "john@doe.com"
		}`,
		responseBody: `{
			"id": "id",
			"customer_reference": "reference",
			"first_name": "john",
			"last_name": "doe",
			"email": "john@doe.com"
		}`,
	}

	c := &CustomerClient{Caller: New(OptHTTPClient(cli))}

	customerParams := CustomerParams{
		CustomerReference: "reference",
		FirstName:         "john",
		LastName:          "doe",
		Email:             "john@doe.com",
	}

	customer, err := c.New(
		context.Background(),
		"idempotency_key",
		&customerParams,
	)

	require.NoError(t, err)
	require.Equal(t, &Customer{
		ID:             "id",
		CustomerParams: customerParams,
	}, customer)
}

func TestCustomerClient_Get(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "GET",
		expectedURL:    "/customers/id",
		responseBody: `{
			"id": "id",
			"customer_reference": "reference",
			"first_name": "john",
			"last_name": "doe",
			"email": "john@doe.com"
		}`,
	}

	c := &CustomerClient{Caller: New(OptHTTPClient(cli))}

	customer, err := c.Get(
		context.Background(),
		"id",
	)

	require.NoError(t, err)
	require.Equal(t, &Customer{
		ID: "id",
		CustomerParams: CustomerParams{
			CustomerReference: "reference",
			FirstName:         "john",
			LastName:          "doe",
			Email:             "john@doe.com",
		},
	}, customer)
}

func TestCustomerClient_GetByReference(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cli := &httpClientMock{
			t:              t,
			expectedMethod: "GET",
			expectedURL:    "/customers?customer_reference=john+doe%3F",
			responseBody: `[{
        		"id": "id",
        		"customer_reference": "john doe?",
				"first_name": "john",
				"last_name": "doe",
				"email": "john@doe.com"
    		}]`,
		}

		c := &CustomerClient{Caller: New(OptHTTPClient(cli))}

		customer, err := c.GetByReference(context.Background(), "john doe?")

		require.NoError(t, err)
		require.Equal(t, &Customer{
			ID: "id",
			CustomerParams: CustomerParams{
				CustomerReference: "john doe?",
				FirstName:         "john",
				LastName:          "doe",
				Email:             "john@doe.com",
			},
		}, customer)
	})

	t.Run("empty customer list", func(t *testing.T) {
		cli := &httpClientMock{
			t:              t,
			expectedMethod: "GET",
			expectedURL:    "/customers?customer_reference=john",
			responseBody:   `[]`,
		}

		c := &CustomerClient{Caller: New(OptHTTPClient(cli))}

		customer, err := c.GetByReference(context.Background(), "john")
		require.EqualError(t, err, "PaymentsOS returned empty array")
		require.Nil(t, customer)
	})

	t.Run("more than one item in customer list", func(t *testing.T) {
		cli := &httpClientMock{
			t:              t,
			expectedMethod: "GET",
			expectedURL:    "/customers?customer_reference=john",
			responseBody:   `[{}, {}]`,
		}

		c := &CustomerClient{Caller: New(OptHTTPClient(cli))}

		customer, err := c.GetByReference(context.Background(), "john")
		require.EqualError(t, err, "PaymentsOS returned array with more than one item")
		require.Nil(t, customer)
	})
}

func TestCustomerClient_Update(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "PUT",
		expectedURL:    "/customers/id",
		expectedBodyJSON: `{
			"customer_reference": "reference",
			"first_name": "john",
			"last_name": "doe",
			"email": "john@doe.com"
		}`,
		responseBody: `{
			"id": "id",
			"customer_reference": "reference",
			"first_name": "john",
			"last_name": "doe",
			"email": "john@doe.com"
		}`,
	}

	c := &CustomerClient{Caller: New(OptHTTPClient(cli))}

	customerParams := CustomerParams{
		CustomerReference: "reference",
		FirstName:         "john",
		LastName:          "doe",
		Email:             "john@doe.com",
	}

	customer, err := c.Update(
		context.Background(),
		"id",
		&customerParams,
	)

	require.NoError(t, err)
	require.Equal(t, &Customer{
		ID:             "id",
		CustomerParams: customerParams,
	}, customer)
}

func TestCustomerClient_Delete(t *testing.T) {
	cli := &httpClientMock{
		t:              t,
		expectedMethod: "DELETE",
		expectedURL:    "/customers/id",
	}

	c := &CustomerClient{Caller: New(OptHTTPClient(cli))}

	err := c.Delete(
		context.Background(),
		"id",
	)

	require.NoError(t, err)
}
