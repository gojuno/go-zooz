package zooz

import "testing"

func TestVoidClient_voidPath(t *testing.T) {
	c := &VoidClient{}
	p := c.voidPath("payment_id", "void_id")
	if p != "payments/payment_id/voids/void_id" {
		t.Errorf("Invalid void path: %s", p)
	}
}
