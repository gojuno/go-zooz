package zooz

import "testing"

func TestAuthorizationClient_authorizationPath(t *testing.T) {
	c := &AuthorizationClient{}
	p := c.authorizationPath("payment_id", "authorization_id")
	if p != "payments/payment_id/authorizations/authorization_id" {
		t.Errorf("Invalid authorization path: %s", p)
	}
}
