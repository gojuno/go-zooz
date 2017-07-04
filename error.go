package zooz

import (
	"encoding/json"
	"fmt"
)

// Error represents possible client error.
type Error struct {
	StatusCode int
	RequestID  string
	ApiError   ApiError
}

// ApiError represents API error response.
// https://developers.paymentsos.com/docs/api#/introduction/responses/errors
type ApiError struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	MoreInfo    string `json:"more_info"`
}

// String implements stringer interface.
func (e ApiError) String() string {
	str, _ := json.Marshal(e)
	return string(str)
}

// Error implements error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("request: %s, status: %d, error: %s", e.RequestID, e.StatusCode, e.ApiError)
}
