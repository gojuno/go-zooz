package zooz

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type httpClientMock struct {
	t *testing.T

	expectedMethod   string
	expectedURL      string            // can be relative i.e. "/somepath?k=v", in this case will be prepended with ApiURL
	expectedHeaders  map[string]string // headers expected to be included in request
	expectedBodyJSON string

	responseCode int
	responseBody string
	error        error
}

func (c *httpClientMock) Do(r *http.Request) (*http.Response, error) {
	require.Equal(c.t, c.expectedMethod, r.Method)
	if len(c.expectedURL) > 0 && c.expectedURL[0] == '/' {
		c.expectedURL = ApiURL + c.expectedURL
	}
	require.Equal(c.t, c.expectedURL, r.URL.String())

	for k, v := range c.expectedHeaders {
		require.Equal(c.t, v, r.Header.Get(k))
	}

	var body []byte
	if r.Body != nil {
		var err error
		body, err = ioutil.ReadAll(r.Body)
		require.NoError(c.t, err)
		defer func() { _ = r.Body.Close() }()
	}
	if c.expectedBodyJSON != "" {
		require.JSONEq(c.t, c.expectedBodyJSON, string(body))
	} else {
		require.Empty(c.t, body)
	}

	if c.error != nil {
		return nil, c.error
	}

	return &http.Response{
		StatusCode: c.responseCode,
		Body:       ioutil.NopCloser(bytes.NewBufferString(c.responseBody)),
	}, nil
}

type request struct {
	Field string `json:"field"`
}

func TestNew(t *testing.T) {
	c := New(
		OptApiURL("http://xxx.com"),
		OptAppID("app_id"),
		OptEnv(EnvLive),
		OptPrivateKey("private_key"),
		OptHTTPClient(&httpClientMock{}),
	)

	if c == nil {
		t.Errorf("Client is nil")
	}
	if c.apiURL != "http://xxx.com/" {
		t.Errorf("Invalid apiURL: %s", c.apiURL)
	}
	if c.appID != "app_id" {
		t.Errorf("Invalid appID: %s", c.appID)
	}
	if c.env != EnvLive {
		t.Errorf("Invalid env: %s", c.env)
	}
	if c.privateKey != "private_key" {
		t.Errorf("Invalid privateKey: %s", c.privateKey)
	}
	if _, ok := c.httpClient.(*httpClientMock); !ok {
		t.Errorf("Invalid httpClient: %T", c.httpClient)
	}
}

func TestNew_Defaults(t *testing.T) {
	c := New()

	if c == nil {
		t.Errorf("Client is nil")
	}
	if c.apiURL != "https://api.paymentsos.com/" {
		t.Errorf("Invalid apiURL: %s", c.apiURL)
	}
	if c.appID != "" {
		t.Errorf("Invalid appID: %s", c.appID)
	}
	if c.env != EnvTest {
		t.Errorf("Invalid env: %s", c.env)
	}
	if c.privateKey != "" {
		t.Errorf("Invalid privateKey: %s", c.privateKey)
	}
	if c.httpClient != http.DefaultClient {
		t.Errorf("Invalid httpClient: %[1]T %#[1]v", c.httpClient)
	}
}

func TestCall_WithApiResponse(t *testing.T) {
	httpClientMock := &httpClientMock{
		t:              t,
		expectedMethod: "POST",
		expectedURL:    "http://xxx.com/somepath?testk=testv",

		expectedHeaders: map[string]string{
			headerAPIVersion: "1.3.0",
			headerEnv:        "test",
			headerAppID:      "app_id_test",
			headerPrivateKey: "private_key_test",
			"Content-Type":   "application/json",
			"test-header":    "test-header-value",
		},

		expectedBodyJSON: `{"field":"request_value"}`,

		responseCode: http.StatusOK,
		responseBody: `{"field":"response_value"}`,
	}

	req := request{
		Field: "request_value",
	}

	response := struct {
		Field string `json:"field"`
	}{}

	client := New(
		OptApiURL("http://xxx.com"),
		OptHTTPClient(httpClientMock),
		OptAppID("app_id_test"),
		OptPrivateKey("private_key_test"),
		OptEnv(EnvTest),
	)

	err := client.Call(
		context.Background(),
		"POST",
		"somepath?testk=testv",
		map[string]string{
			"test-header": "test-header-value",
		},
		&req,
		&response,
	)
	require.NoError(t, err)
	require.Equal(t, "response_value", response.Field)
}

func TestCall_WithApiError(t *testing.T) {
	httpClientMock := &httpClientMock{
		t:                t,
		expectedMethod:   "POST",
		expectedURL:      "/somepath",
		expectedBodyJSON: `{"field":"request_value"}`,

		responseCode: http.StatusBadRequest,
		responseBody: `{"category":"category_test"}`,
	}

	req := request{
		Field: "request_value",
	}

	client := New(OptHTTPClient(httpClientMock))

	err := client.Call(
		context.Background(),
		"POST",
		"somepath",
		nil,
		&req,
		nil,
	)

	require.Error(t, err)
	require.Equal(t, &Error{
		StatusCode: http.StatusBadRequest,
		APIError: APIError{
			Category: "category_test",
		},
	}, err)
}

func TestCall_WithTransportError(t *testing.T) {
	httpClientMock := &httpClientMock{
		t:                t,
		expectedMethod:   "POST",
		expectedURL:      "/somepath",
		expectedBodyJSON: `{"field":"request_value"}`,
		error:            errors.New("do_error"),
	}

	req := request{
		Field: "request_value",
	}

	client := New(OptHTTPClient(httpClientMock))

	err := client.Call(
		context.Background(),
		"POST",
		"somepath",
		nil,
		&req,
		nil,
	)

	require.Error(t, err)
	require.EqualError(t, err, "failed to do request: do_error")
}
