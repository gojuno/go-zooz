package zooz

import "testing"

func TestPaymentMethodClient_paymentMethodsPath(t *testing.T) {
	c := &PaymentMethodClient{}
	p := c.paymentMethodsPath("customer_id")
	if p != "customers/customer_id/payment-methods" {
		t.Errorf("Invalid payment methods path: %s", p)
	}
}

func TestPaymentMethodClient_tokenPath(t *testing.T) {
	c := &PaymentMethodClient{}
	p := c.tokenPath("customer_id", "token")
	if p != "customers/customer_id/payment-methods/token" {
		t.Errorf("Invalid token path: %s", p)
	}
}
