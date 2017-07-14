package zooz

import "testing"

func TestRefundClient_refundPath(t *testing.T) {
	c := &RefundClient{}
	p := c.refundPath("payment_id", "refund_id")
	if p != "payments/payment_id/refunds/refund_id" {
		t.Errorf("Invalid refund path: %s", p)
	}
}
