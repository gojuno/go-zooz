package zooz

import "testing"

func TestCaptureClient_capturePath(t *testing.T) {
	c := &CaptureClient{}
	p := c.capturePath("payment_id", "capture_id")
	if p != "payments/payment_id/captures/capture_id" {
		t.Errorf("Invalid capture path: %s", p)
	}
}
