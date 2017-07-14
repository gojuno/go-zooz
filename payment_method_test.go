package zooz

import "testing"

func TestPaymentClient_paymentPath(t *testing.T) {
	c := &PaymentClient{}
	p := c.paymentPath("payment_id", PaymentExpandAuthorizations, PaymentExpandCaptures)
	if p != "payments/payment_id?expand=authorizations&expand=captures" {
		t.Errorf("Invalid payment path: %s", p)
	}
}
