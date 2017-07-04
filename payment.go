package zooz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// PaymentClient is a client for work with Payment entity.
// https://developers.paymentsos.com/docs/api#/reference/payments
type PaymentClient struct {
	Caller Caller
}

// Payment is a model of entity.
type Payment struct {
	PaymentParams

	ID                  string          `json:"id"`
	Created             json.Number     `json:"created"`
	Modified            json.Number     `json:"modified"`
	Status              paymentStatus   `json:"status"`
	PossibleNextActions []PaymentAction `json:"possible_next_actions"`

	// Expansions
	PaymentMethod    *PaymentMethod           `json:"payment_method"`
	Customer         *Customer                `json:"customer"`
	RelatedResources *PaymentRelatedResources `json:"related_resources"`
}

// PaymentParams is a set of params for creating and updating entity.
type PaymentParams struct {
	Amount                  float64           `json:"amount"`
	Currency                string            `json:"currency"`
	CustomerID              string            `json:"customer_id,omitempty"`
	AdditionalDetails       AdditionalDetails `json:"additional_details,omitempty"`
	StatementSoftDescriptor string            `json:"statement_soft_descriptor,omitempty"`
	Order                   *PaymentOrder     `json:"order,omitempty"`
	ShippingAddress         *Address          `json:"shipping_address,omitempty"`
	BillingAddress          *Address          `json:"billing_address,omitempty"`
}

// PaymentOrder represents order description.
type PaymentOrder struct {
	ID                string                 `json:"id,omitempty"`
	AdditionalDetails AdditionalDetails      `json:"additional_details,omitempty"`
	TaxAmount         float64                `json:"tax_amount,omitempty"`
	TaxPercentage     float64                `json:"tax_percentage,omitempty"`
	LineItems         []PaymentOrderLineItem `json:"line_items,omitempty"`
}

// PaymentOrderLineItem represents one item of order.
type PaymentOrderLineItem struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Quantity  int     `json:"quantity,omitempty"`
	UnitPrice float64 `json:"unit_price,omitempty"`
}

// PaymentAction represents action which may be performed on Payment entity.
type PaymentAction struct {
	Action paymentAction `json:"action"`
	Href   string        `json:"href"`
}

// PaymentRelatedResources is a set of resources related to Payment.
type PaymentRelatedResources struct {
	Authorizations []Authorization `json:"authorizations"`
	Charges        []Charge        `json:"charges"`
	Voids          []Void          `json:"voids"`
	Redirections   []Redirection   `json:"redirections"`
	Captures       []Capture       `json:"captures"`
	Refunds        []Refund        `json:"refunds"`
}

type paymentStatus string

type paymentExpand string

type paymentAction string

const (
	paymentsPath = "payments"

	// List of possible payment status values.
	PaymentStatusInitialized paymentStatus = "Initialized"
	PaymentStatusPending     paymentStatus = "Pending"
	PaymentStatusAuthorized  paymentStatus = "Authorized"
	PaymentStatusCaptured    paymentStatus = "Captured"
	PaymentStatusRefunded    paymentStatus = "Refunded"
	PaymentStatusVoided      paymentStatus = "Voided"

	// List of possible payment expansion values.
	PaymentExpandAuthorizations paymentExpand = "authorizations"
	PaymentExpandRedirections   paymentExpand = "redirections"
	PaymentExpandCaptures       paymentExpand = "captures"
	PaymentExpandRefunds        paymentExpand = "refunds"
	PaymentExpandVoids          paymentExpand = "voids"
	PaymentExpandCredits        paymentExpand = "credits"
	PaymentExpandCustomer       paymentExpand = "customer"
	PaymentExpandPaymentMethod  paymentExpand = "payment_method"
	PaymentExpandAll            paymentExpand = "all"

	// List of possible payment action values.
	PaymentActionAuthorize     paymentAction = "Authorize"
	PaymentActionCharge        paymentAction = "Charge"
	PaymentActionCapture       paymentAction = "Capture"
	PaymentActionRefund        paymentAction = "Refund"
	PaymentActionVoid          paymentAction = "Void"
	PaymentActionUpdatePayment paymentAction = "Update Payment"
)

// New creates new Payment entity.
func (c *PaymentClient) New(ctx context.Context, idempotencyKey string, params *PaymentParams) (*Payment, error) {
	payment := &Payment{}
	if err := c.Caller.Call(ctx, "POST", paymentsPath, map[string]string{headerIdempotencyKey: idempotencyKey}, params, payment); err != nil {
		return nil, err
	}
	return payment, nil
}

// Get returns Payment entity with optional expansions. You may specify any number of expansion or
// use zooz.PaymentExpandAll for expand payments with all expansions.
func (c *PaymentClient) Get(ctx context.Context, id string, expands ...paymentExpand) (*Payment, error) {
	payment := &Payment{}
	if err := c.Caller.Call(ctx, "GET", c.paymentPath(id, expands...), nil, nil, payment); err != nil {
		return nil, err
	}
	return payment, nil
}

// Update changes Payment entity and returned updated entity.
// Payment details can only be updated if no other action has been performed on the Payment resource.
// Note: In addition to the fields that you want to update, you must re-send all the other original argument fields,
// because this operation replaces the Payment resource.
func (c *PaymentClient) Update(ctx context.Context, id string, params *PaymentParams) (*Payment, error) {
	payment := &Payment{}
	if err := c.Caller.Call(ctx, "PUT", c.paymentPath(id), nil, params, payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (c *PaymentClient) paymentPath(id string, expands ...paymentExpand) string {
	values := url.Values{}
	for _, expand := range expands {
		values.Add("expand", string(expand))
	}
	query := values.Encode()
	if query != "" {
		query = "?" + query
	}

	return fmt.Sprintf("%s/%s%s", paymentsPath, id, query)
}
