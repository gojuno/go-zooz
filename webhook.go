package zooz

import "time"

type Webhook struct {
	ID        string      `json:"id"`
	Created   time.Time   `json:"created"`
	PaymentID string      `json:"payment_id"`
	AccountID string      `json:"account_id"`
	AppID     string      `json:"app_id"`
	Data      interface{} `json:"data"`
}
