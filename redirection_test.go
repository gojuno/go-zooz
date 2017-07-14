package zooz

import "testing"

func TestRedirectionClient_redirectionPath(t *testing.T) {
	c := &RedirectionClient{}
	p := c.redirectionPath("payment_id", "redirection_id")
	if p != "payments/payment_id/redirections/redirection_id" {
		t.Errorf("Invalid redirection path: %s", p)
	}
}
