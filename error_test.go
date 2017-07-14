package zooz

import "testing"

func TestError(t *testing.T) {
	e := &Error{
		StatusCode: 500,
		RequestID:  "request_id",
		APIError: APIError{
			Category:    "cat",
			Description: "desc",
			MoreInfo:    "info",
		},
	}

	var err error = e

	if err.Error() != `request: request_id, status: 500, error: {"category":"cat","description":"desc","more_info":"info"}` {
		t.Errorf("Invalid error: %s", err.Error())
	}
}
