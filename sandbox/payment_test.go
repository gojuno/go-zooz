package sandbox

import (
	"context"
	"github.com/gtforge/go-zooz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestPayment(t *testing.T) {
	t.Parallel()
	client := GetClient(t)

	t.Run("new (all fields) & get & get-expanded", func(t *testing.T) {
		t.Parallel()

		customer := PrepareCustomer(t, client)

		paymentCreated, err := client.Payment().New(context.Background(), randomString(32), &zooz.PaymentParams{
			Amount:     5000,
			Currency:   "USD",
			CustomerID: customer.ID,
			AdditionalDetails: zooz.AdditionalDetails{
				"payment detail 1": "value 1",
				"payment detail 2": "value 2",
			},
			StatementSoftDescriptor: "statement soft descriptor",
			Order: &zooz.PaymentOrder{
				ID: "order-id",
				AdditionalDetails: zooz.AdditionalDetails{
					"order detail 1": "value 1",
					"order detail 2": "value 2",
				},
				TaxAmount:     30,
				TaxPercentage: 10,
				LineItems: []zooz.PaymentOrderLineItem{
					{
						ID:        "line item 1 id",
						Name:      "line item 1 name",
						Quantity:  5,
						UnitPrice: 6,
					},
				},
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
			BillingAddress: &zooz.Address{
				Country:   "USA",
				State:     "billing state",
				City:      "billing city",
				Line1:     "billing line1",
				Line2:     "billing line2",
				ZipCode:   "billing zip code",
				Title:     "billing title",
				FirstName: "billing first name",
				LastName:  "billing last name",
				Phone:     "billing phone",
				Email:     "billing-address@email.com",
			},
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, paymentCreated.ID)
			assert.NotEmpty(t, paymentCreated.Created)
			assert.NotEmpty(t, paymentCreated.Modified)
			assert.Equal(t, &zooz.Payment{
				PaymentParams: zooz.PaymentParams{
					Amount:     5000,
					Currency:   "USD",
					CustomerID: "", // doesn't exist in response, just in request
					AdditionalDetails: zooz.AdditionalDetails{
						"payment detail 1": "value 1",
						"payment detail 2": "value 2",
					},
					StatementSoftDescriptor: "statement soft descriptor",
					Order: &zooz.PaymentOrder{
						ID: "order-id",
						AdditionalDetails: zooz.AdditionalDetails{
							"order detail 1": "value 1",
							"order detail 2": "value 2",
						},
						TaxAmount:     30,
						TaxPercentage: 10,
						LineItems: []zooz.PaymentOrderLineItem{
							{
								ID:        "line item 1 id",
								Name:      "line item 1 name",
								Quantity:  5,
								UnitPrice: 6,
							},
						},
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
					BillingAddress: &zooz.Address{
						Country:   "USA",
						State:     "billing state",
						City:      "billing city",
						Line1:     "billing line1",
						Line2:     "billing line2",
						ZipCode:   "billing zip code",
						Title:     "billing title",
						FirstName: "billing first name",
						LastName:  "billing last name",
						Phone:     "billing phone",
						Email:     "billing-address@email.com",
					},
				},
				ID:                  paymentCreated.ID,       // ignore
				Created:             paymentCreated.Created,  // ignore
				Modified:            paymentCreated.Modified, // ignore
				Status:              "Initialized",
				PossibleNextActions: paymentCreated.PossibleNextActions, // ignore
				PaymentMethod:       nil,
				Customer:            CustomerOnlyHref(customer), // according to docs should be full customer, but is filled only on Payment Get with expansion
				RelatedResources:    nil,
			}, paymentCreated)
		})

		paymentRetrieved, err := client.Payment().Get(context.Background(), paymentCreated.ID)
		require.NoError(t, err)
		require.Equal(t, paymentCreated, paymentRetrieved)

		paymentRetrievedExpandCustomer, err := client.Payment().Get(context.Background(), paymentCreated.ID, zooz.PaymentExpandCustomer)
		require.NoError(t, err)
		must(t, func() {
			assert.Equal(t, CustomerWithHref(customer), paymentRetrievedExpandCustomer.Customer)
			paymentRetrievedExpandCustomer.Customer = paymentCreated.Customer
			assert.Equal(t, paymentCreated, paymentRetrievedExpandCustomer)
		})
	})

	t.Run("new (required fields) & update & get", func(t *testing.T) {
		t.Parallel()

		const amount1, currency1 = 5000, "USD"
		const amount2, currency2 = 6000, "RUB"

		paymentCreated, err := client.Payment().New(context.Background(), randomString(32), &zooz.PaymentParams{
			Amount:                  amount1,
			Currency:                currency1,
			CustomerID:              "",
			AdditionalDetails:       nil,
			StatementSoftDescriptor: "",
			Order:                   nil,
			ShippingAddress:         nil,
			BillingAddress:          nil,
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEmpty(t, paymentCreated.ID)
			assert.NotEmpty(t, paymentCreated.Created)
			assert.NotEmpty(t, paymentCreated.Modified)
			assert.Equal(t, &zooz.Payment{
				PaymentParams: zooz.PaymentParams{
					Amount:                  5000,
					Currency:                "USD",
					CustomerID:              "", // doesn't exist in response, just in request
					AdditionalDetails:       nil,
					StatementSoftDescriptor: "",
					Order:                   nil,
					ShippingAddress:         nil,
					BillingAddress:          nil,
				},
				ID:                  paymentCreated.ID,       // ignore
				Created:             paymentCreated.Created,  // ignore
				Modified:            paymentCreated.Modified, // ignore
				Status:              "Initialized",
				PossibleNextActions: paymentCreated.PossibleNextActions, // ignore
				PaymentMethod:       nil,
				Customer:            nil,
				RelatedResources:    nil,
			}, paymentCreated)
		})

		customer := PrepareCustomer(t, client)

		paymentUpdated, err := client.Payment().Update(context.Background(), paymentCreated.ID, &zooz.PaymentParams{
			Amount:     amount2,
			Currency:   currency2,
			CustomerID: customer.ID,
			AdditionalDetails: zooz.AdditionalDetails{
				"payment detail 1": "value 1",
				"payment detail 2": "value 2",
			},
			StatementSoftDescriptor: "statement soft descriptor",
			Order: &zooz.PaymentOrder{
				ID: "order-id",
				AdditionalDetails: zooz.AdditionalDetails{
					"order detail 1": "value 1",
					"order detail 2": "value 2",
				},
				TaxAmount:     30,
				TaxPercentage: 10,
				LineItems: []zooz.PaymentOrderLineItem{
					{
						ID:        "line item 1 id",
						Name:      "line item 1 name",
						Quantity:  5,
						UnitPrice: 6,
					},
				},
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
			BillingAddress: &zooz.Address{
				Country:   "USA",
				State:     "billing state",
				City:      "billing city",
				Line1:     "billing line1",
				Line2:     "billing line2",
				ZipCode:   "billing zip code",
				Title:     "billing title",
				FirstName: "billing first name",
				LastName:  "billing last name",
				Phone:     "billing phone",
				Email:     "billing-address@email.com",
			},
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEqual(t, paymentCreated.Modified, paymentUpdated.Modified)
			assert.Equal(t, &zooz.Payment{
				PaymentParams: zooz.PaymentParams{
					Amount:     amount2,
					Currency:   currency2,
					CustomerID: "", // doesn't exist in response, just in request
					AdditionalDetails: zooz.AdditionalDetails{
						"payment detail 1": "value 1",
						"payment detail 2": "value 2",
					},
					StatementSoftDescriptor: "statement soft descriptor",
					Order: &zooz.PaymentOrder{
						ID: "order-id",
						AdditionalDetails: zooz.AdditionalDetails{
							"order detail 1": "value 1",
							"order detail 2": "value 2",
						},
						TaxAmount:     30,
						TaxPercentage: 10,
						LineItems: []zooz.PaymentOrderLineItem{
							{
								ID:        "line item 1 id",
								Name:      "line item 1 name",
								Quantity:  5,
								UnitPrice: 6,
							},
						},
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
					BillingAddress: &zooz.Address{
						Country:   "USA",
						State:     "billing state",
						City:      "billing city",
						Line1:     "billing line1",
						Line2:     "billing line2",
						ZipCode:   "billing zip code",
						Title:     "billing title",
						FirstName: "billing first name",
						LastName:  "billing last name",
						Phone:     "billing phone",
						Email:     "billing-address@email.com",
					},
				},
				ID:                  paymentCreated.ID,
				Created:             paymentCreated.Created,
				Modified:            paymentUpdated.Modified, // ignore
				Status:              "Initialized",
				PossibleNextActions: paymentCreated.PossibleNextActions, // ignore
				PaymentMethod:       nil,
				Customer:            CustomerOnlyHref(customer), // according to docs should be full customer, but is filled only on Payment Get with expansion
				RelatedResources:    nil,
			}, paymentUpdated)
		})

		paymentRetrieved, err := client.Payment().Get(context.Background(), paymentCreated.ID)
		require.NoError(t, err)
		require.Equal(t, paymentUpdated, paymentRetrieved)
	})

	t.Run("update (all fields) & get", func(t *testing.T) {
		t.Parallel()

		customer1 := PrepareCustomer(t, client)
		customer2 := PrepareCustomer(t, client)
		payment := PreparePayment(t, client, customer1.ID)

		paymentUpdated, err := client.Payment().Update(context.Background(), payment.ID, &zooz.PaymentParams{
			Amount:     6000,
			Currency:   "EUR",
			CustomerID: customer2.ID,
			AdditionalDetails: zooz.AdditionalDetails{
				"payment detail 1": "value 1-2",
				"payment detail 3": "value 3",
			},
			StatementSoftDescriptor: "statement soft descriptor 2",
			Order: &zooz.PaymentOrder{
				ID: "order-id-2",
				AdditionalDetails: zooz.AdditionalDetails{
					"order detail 1": "value 1-2",
					"order detail 3": "value 3",
				},
				TaxAmount:     50,
				TaxPercentage: 20,
				LineItems: []zooz.PaymentOrderLineItem{
					{
						ID:        "line item 1 id 2",
						Name:      "line item 1 name 2",
						Quantity:  10,
						UnitPrice: 7,
					},
				},
			},
			ShippingAddress: &zooz.Address{
				Country:   "RUS",
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
			BillingAddress: &zooz.Address{
				Country:   "RUS",
				State:     "billing state 2",
				City:      "billing city 2",
				Line1:     "billing line1 2",
				Line2:     "billing line2 2",
				ZipCode:   "billing zip code 2",
				Title:     "billing title 2",
				FirstName: "billing first name 2",
				LastName:  "billing last name 2",
				Phone:     "billing phone 2",
				Email:     "billing-address-2@email.com",
			},
		})
		require.NoError(t, err)
		must(t, func() {
			assert.NotEqual(t, payment.Modified, paymentUpdated.Modified)
			assert.Equal(t, &zooz.Payment{
				PaymentParams: zooz.PaymentParams{
					Amount:     6000,
					Currency:   "EUR",
					CustomerID: "", // doesn't exist in response, just in request
					AdditionalDetails: zooz.AdditionalDetails{
						"payment detail 1": "value 1-2",
						"payment detail 3": "value 3",
					},
					StatementSoftDescriptor: "statement soft descriptor 2",
					Order: &zooz.PaymentOrder{
						ID: "order-id-2",
						AdditionalDetails: zooz.AdditionalDetails{
							"order detail 1": "value 1-2",
							"order detail 3": "value 3",
						},
						TaxAmount:     50,
						TaxPercentage: 20,
						LineItems: []zooz.PaymentOrderLineItem{
							{
								ID:        "line item 1 id 2",
								Name:      "line item 1 name 2",
								Quantity:  10,
								UnitPrice: 7,
							},
						},
					},
					ShippingAddress: &zooz.Address{
						Country:   "RUS",
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
					BillingAddress: &zooz.Address{
						Country:   "RUS",
						State:     "billing state 2",
						City:      "billing city 2",
						Line1:     "billing line1 2",
						Line2:     "billing line2 2",
						ZipCode:   "billing zip code 2",
						Title:     "billing title 2",
						FirstName: "billing first name 2",
						LastName:  "billing last name 2",
						Phone:     "billing phone 2",
						Email:     "billing-address-2@email.com",
					},
				},
				ID:                  payment.ID,
				Created:             payment.Created,
				Modified:            paymentUpdated.Modified, // ignore
				Status:              "Initialized",
				PossibleNextActions: paymentUpdated.PossibleNextActions, // ignore
				PaymentMethod:       nil,
				Customer:            CustomerOnlyHref(customer2), // according to docs should be full customer, but is filled only on Payment Get with expansion
				RelatedResources:    nil,
			}, paymentUpdated)
		})

		paymentRetrieved, err := client.Payment().Get(context.Background(), payment.ID)
		require.NoError(t, err)
		require.Equal(t, paymentUpdated, paymentRetrieved)
	})

	t.Run("get - unknown payment", func(t *testing.T) {
		t.Parallel()

		_, err := client.Payment().Get(context.Background(), "00000000-0000-1000-8000-000000000000")
		zoozErr := &zooz.Error{}
		require.ErrorAs(t, err, &zoozErr)
		require.Equal(t, &zooz.Error{
			StatusCode: http.StatusNotFound,
			RequestID:  zoozErr.RequestID, // ignore
			APIError: zooz.APIError{
				Category:    "api_request_error",
				Description: "The resource was not found.",
				MoreInfo:    "Payment resource does not exists",
			},
		}, zoozErr)
	})
}

// PreparePayment is a helper to create new payment in zooz.
func PreparePayment(t *testing.T, client *zooz.Client, customerID string) *zooz.Payment {
	payment, err := client.Payment().New(context.Background(), randomString(32), &zooz.PaymentParams{
		Amount:     5000,
		Currency:   "USD",
		CustomerID: customerID,
		AdditionalDetails: zooz.AdditionalDetails{
			"payment detail 1": "value 1",
			"payment detail 2": "value 2",
		},
		StatementSoftDescriptor: "statement soft descriptor",
		Order: &zooz.PaymentOrder{
			ID: "order-id",
			AdditionalDetails: zooz.AdditionalDetails{
				"order detail 1": "value 1",
				"order detail 2": "value 2",
			},
			TaxAmount:     30,
			TaxPercentage: 10,
			LineItems: []zooz.PaymentOrderLineItem{
				{
					ID:        "line item 1 id",
					Name:      "line item 1 name",
					Quantity:  5,
					UnitPrice: 6,
				},
			},
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
		BillingAddress: &zooz.Address{
			Country:   "USA",
			State:     "billing state",
			City:      "billing city",
			Line1:     "billing line1",
			Line2:     "billing line2",
			ZipCode:   "billing zip code",
			Title:     "billing title",
			FirstName: "billing first name",
			LastName:  "billing last name",
			Phone:     "billing phone",
			Email:     "billing-address@email.com",
		},
	})
	require.NoError(t, err)
	return payment
}
