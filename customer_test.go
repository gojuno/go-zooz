package zooz

import "testing"

func TestCustomerClient_customerPath(t *testing.T) {
	c := &CustomerClient{}
	p := c.customerPath("customer_id")
	if p != "customers/customer_id" {
		t.Errorf("Invalid customer path: %s", p)
	}
}
