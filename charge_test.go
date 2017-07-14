package zooz

import "testing"

func TestChargeClient_chargePath(t *testing.T) {
	c := &ChargeClient{}
	p := c.chargePath("payment_id", "charge_id")
	if p != "payments/payment_id/charges/charge_id" {
		t.Errorf("Invalid charge path: %s", p)
	}
}
