package zooz

import (
	"encoding/json"
)

// Dispute is an entity model.
type Dispute struct {
	ID                 string             `json:"id"`
	RelatedTransaction RelatedTransaction `json:"related_transaction"`
	Result             Result             `json:"result"`
	Amount             float32            `json:"amount"`
	Currency           string             `json:"currency"`
	ProviderData       ProviderData       `json:"provider_data"`
	Created            json.Number        `json:"created"`
}

// RelatedTransaction is a set of params for creating entity.
type RelatedTransaction struct {
	ID   string `json:"id"`
	Href string `json:"href,omitempty"`
}
