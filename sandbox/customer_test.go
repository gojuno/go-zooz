package sandbox

import (
	"context"
	"github.com/gtforge/go-zooz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCustomer(t *testing.T) {
	t.Parallel()
	client := GetClient(t)

	t.Run("new (all fields) & get & get-by-reference", func(t *testing.T) {
		t.Parallel()

		customerReference := randomString(32)

		customerCreated, err := client.Customer().New(context.Background(), randomString(32), &zooz.CustomerParams{
			CustomerReference: customerReference,
			FirstName:         "first name",
			LastName:          "last name",
			Email:             "customer@email.com",
			AdditionalDetails: zooz.AdditionalDetails{
				"customer detail 1": "value 1",
				"customer detail 2": "value 2",
			},
			ShippingAddress: &zooz.Address{
				Country:   "RUS",
				State:     "shipping state",
				City:      "shipping city",
				Line1:     "shipping line1",
				Line2:     "shipping line2",
				ZipCode:   "shipping zip code",
				Title:     "shipping title",
				FirstName: "shipping first name",
				LastName:  "shipping last name",
				Phone:     "shipping phone",
				Email:     "shipping-address@email.com",
			},
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, customerCreated.ID)
			assert.NotEmpty(t, customerCreated.Created)
			assert.NotEmpty(t, customerCreated.Modified)
			assert.Equal(t, &zooz.Customer{
				CustomerParams: zooz.CustomerParams{
					CustomerReference: customerReference,
					FirstName:         "first name",
					LastName:          "last name",
					Email:             "customer@email.com",
					AdditionalDetails: zooz.AdditionalDetails{
						"customer detail 1": "value 1",
						"customer detail 2": "value 2",
					},
					ShippingAddress: &zooz.Address{
						Country:   "RUS",
						State:     "shipping state",
						City:      "shipping city",
						Line1:     "shipping line1",
						Line2:     "shipping line2",
						ZipCode:   "shipping zip code",
						Title:     "shipping title",
						FirstName: "shipping first name",
						LastName:  "shipping last name",
						Phone:     "shipping phone",
						Email:     "shipping-address@email.com",
					},
				},
				ID:             customerCreated.ID,       // ignore
				Created:        customerCreated.Created,  // ignore
				Modified:       customerCreated.Modified, // ignore
				PaymentMethods: nil,
				Href:           "", // this field exists only when customer is a part of payment object
			}, customerCreated)
		})

		customerRetrieved, err := client.Customer().Get(context.Background(), customerCreated.ID)
		require.NoError(t, err)
		require.Equal(t, customerCreated, customerRetrieved)

		customerRetrievedByReference, err := client.Customer().GetByReference(context.Background(), customerReference)
		require.NoError(t, err)
		require.Equal(t, customerRetrieved, customerRetrievedByReference)
	})

	t.Run("new (required fields) & update & get", func(t *testing.T) {
		t.Parallel()

		customerReference := randomString(32)
		customerReference2 := randomString(32) // it is not possible to change customer reference

		customerCreated, err := client.Customer().New(context.Background(), randomString(32), &zooz.CustomerParams{
			CustomerReference: customerReference,
			FirstName:         "",
			LastName:          "",
			Email:             "",
			AdditionalDetails: nil,
			ShippingAddress:   nil,
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, customerCreated.ID)
			assert.NotEmpty(t, customerCreated.Created)
			assert.NotEmpty(t, customerCreated.Modified)
			assert.Equal(t, &zooz.Customer{
				CustomerParams: zooz.CustomerParams{
					CustomerReference: customerReference,
					FirstName:         "",
					LastName:          "",
					Email:             "",
					AdditionalDetails: nil,
					ShippingAddress:   nil,
				},
				ID:             customerCreated.ID,       // ignore
				Created:        customerCreated.Created,  // ignore
				Modified:       customerCreated.Modified, // ignore
				PaymentMethods: nil,
				Href:           "", // this field exists only when customer is a part of payment object
			}, customerCreated)
		})

		customerUpdated, err := client.Customer().Update(context.Background(), customerCreated.ID, &zooz.CustomerParams{
			CustomerReference: customerReference2,
			FirstName:         "first name",
			LastName:          "last name",
			Email:             "customer@email.com",
			AdditionalDetails: zooz.AdditionalDetails{
				"customer detail 1": "value 1",
				"customer detail 2": "value 2",
			},
			ShippingAddress: &zooz.Address{
				Country:   "USA",
				State:     "shipping state",
				City:      "shipping city",
				Line1:     "shipping line1",
				Line2:     "shipping line2",
				ZipCode:   "shipping zip code",
				Title:     "shipping title",
				FirstName: "shipping first name",
				LastName:  "shipping last name",
				Phone:     "shipping phone",
				Email:     "shipping-address@email.com",
			},
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEqual(t, customerCreated.Modified, customerUpdated.Modified)
			assert.Equal(t, &zooz.Customer{
				CustomerParams: zooz.CustomerParams{
					CustomerReference: customerReference, // it is not possible to change customer reference
					FirstName:         "first name",
					LastName:          "last name",
					Email:             "customer@email.com",
					AdditionalDetails: zooz.AdditionalDetails{
						"customer detail 1": "value 1",
						"customer detail 2": "value 2",
					},
					ShippingAddress: &zooz.Address{
						Country:   "USA",
						State:     "shipping state",
						City:      "shipping city",
						Line1:     "shipping line1",
						Line2:     "shipping line2",
						ZipCode:   "shipping zip code",
						Title:     "shipping title",
						FirstName: "shipping first name",
						LastName:  "shipping last name",
						Phone:     "shipping phone",
						Email:     "shipping-address@email.com",
					},
				},
				ID:             customerCreated.ID,
				Created:        customerCreated.Created,
				Modified:       customerUpdated.Modified, // ignore
				PaymentMethods: nil,
				Href:           "", // this field exists only when customer is a part of payment object
			}, customerUpdated)
		})

		customerRetrieved, err := client.Customer().Get(context.Background(), customerCreated.ID)
		require.NoError(t, err)
		require.Equal(t, customerUpdated, customerRetrieved)

		customerRetrievedByReference, err := client.Customer().GetByReference(context.Background(), customerReference)
		require.NoError(t, err)
		require.Equal(t, customerUpdated, customerRetrievedByReference)
	})

	t.Run("update (all fields) & get", func(t *testing.T) {
		t.Parallel()

		customer := PrepareCustomer(t, client)
		customerReference2 := randomString(32) // it is not possible to change customer reference

		customerUpdated, err := client.Customer().Update(context.Background(), customer.ID, &zooz.CustomerParams{
			CustomerReference: customerReference2,
			FirstName:         "first name 2",
			LastName:          "last name 2",
			Email:             "customer-2@email.com",
			AdditionalDetails: zooz.AdditionalDetails{
				"customer detail 1": "value 1-2",
				"customer detail 3": "value 3",
			},
			ShippingAddress: &zooz.Address{
				Country:   "USA",
				State:     "shipping state 2",
				City:      "shipping city 2",
				Line1:     "shipping line1 2",
				Line2:     "shipping line2 2",
				ZipCode:   "shipping zip code 2",
				Title:     "shipping title 2",
				FirstName: "shipping first name 2",
				LastName:  "shipping last name 2",
				Phone:     "shipping phone 2",
				Email:     "shipping-address-2@email.com",
			},
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEqual(t, customer.Modified, customerUpdated.Modified)
			assert.Equal(t, &zooz.Customer{
				CustomerParams: zooz.CustomerParams{
					CustomerReference: customer.CustomerReference, // it is not possible to change customer reference
					FirstName:         "first name 2",
					LastName:          "last name 2",
					Email:             "customer-2@email.com",
					AdditionalDetails: zooz.AdditionalDetails{
						"customer detail 1": "value 1-2",
						"customer detail 3": "value 3",
					},
					ShippingAddress: &zooz.Address{
						Country:   "USA",
						State:     "shipping state 2",
						City:      "shipping city 2",
						Line1:     "shipping line1 2",
						Line2:     "shipping line2 2",
						ZipCode:   "shipping zip code 2",
						Title:     "shipping title 2",
						FirstName: "shipping first name 2",
						LastName:  "shipping last name 2",
						Phone:     "shipping phone 2",
						Email:     "shipping-address-2@email.com",
					},
				},
				ID:             customer.ID,
				Created:        customer.Created,
				Modified:       customerUpdated.Modified, // ignore
				PaymentMethods: nil,
				Href:           "", // this field exists only when customer is a part of payment object
			}, customerUpdated)
		})

		customerRetrieved, err := client.Customer().Get(context.Background(), customer.ID)
		require.NoError(t, err)
		require.Equal(t, customerUpdated, customerRetrieved)

		customerRetrievedByReference, err := client.Customer().GetByReference(context.Background(), customer.CustomerReference)
		require.NoError(t, err)
		require.Equal(t, customerUpdated, customerRetrievedByReference)
	})

	t.Run("delete & get", func(t *testing.T) {
		t.Parallel()

		customer := PrepareCustomer(t, client)

		err := client.Customer().Delete(context.Background(), customer.ID)
		require.NoError(t, err)

		_, err = client.Customer().Get(context.Background(), customer.ID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "Customer not found.",
		})
	})

	t.Run("delete & get-by-reference (unexpected error)", func(t *testing.T) {
		t.Parallel()

		customer := PrepareCustomer(t, client)

		err := client.Customer().Delete(context.Background(), customer.ID)
		require.NoError(t, err)

		_, err = client.Customer().GetByReference(context.Background(), customer.CustomerReference)
		require.EqualError(t, err, "PaymentsOS returned empty array") // wow, that is unexpected
	})

	t.Run("get - unknown customer", func(t *testing.T) {
		t.Parallel()

		_, err := client.Customer().Get(context.Background(), UnknownUUID)
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "",
		})
	})

	t.Run("get-by-reference - unknown customer", func(t *testing.T) {
		t.Parallel()

		_, err := client.Customer().GetByReference(context.Background(), randomString(32))
		requireZoozError(t, err, http.StatusNotFound, zooz.APIError{
			Category:    "api_request_error",
			Description: "The resource was not found.",
			MoreInfo:    "",
		})
	})
}

// PrepareCustomer is a helper to create new customer in zooz.
func PrepareCustomer(t *testing.T, client *zooz.Client) *zooz.Customer {
	customer, err := client.Customer().New(context.Background(), randomString(32), &zooz.CustomerParams{
		CustomerReference: randomString(32),
		FirstName:         "customer first name",
		LastName:          "customer last name",
		Email:             "customer@email.com",
		AdditionalDetails: zooz.AdditionalDetails{
			"customer detail 1": "customer detail 1 value",
			"customer detail 2": "customer detail 2 value",
		},
		ShippingAddress: &zooz.Address{
			Country: "RUS",
		},
	})
	require.NoError(t, err)
	return customer
}

// CustomerOnlyHref is a helper to get customer object that contains only href.
// Such object returned if no expansion requested.
func CustomerOnlyHref(customer *zooz.Customer) *zooz.Customer {
	return &zooz.Customer{
		Href: zooz.ApiURL + "/customers/" + customer.ID,
	}
}

// CustomerWithHref is a helper to add href field to customer.
func CustomerWithHref(customer *zooz.Customer) *zooz.Customer {
	return &zooz.Customer{
		CustomerParams: customer.CustomerParams,
		ID:             customer.ID,
		Created:        customer.Created,
		Modified:       customer.Modified,
		PaymentMethods: customer.PaymentMethods,
		Href:           zooz.ApiURL + "/customers/" + customer.ID,
	}
}
