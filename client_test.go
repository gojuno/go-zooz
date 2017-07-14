package zooz

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

type httpClientMock struct {
	do func(r *http.Request) (*http.Response, error)
}

type request struct {
	Field string `json:"field"`
}

type callerMock struct {
	t *testing.T
	expectedMethod string
	expectedPath string
	expectedHeaders map[string]string
	expectedReqObj interface{}
	returnRespObj interface{}
	returnError error
}

func (c *httpClientMock) Do(r *http.Request) (*http.Response, error) {
	return c.do(r)
}

func (c *callerMock) Call(ctx context.Context, method, path string, headers map[string]string, reqObj interface{}, respObj interface{}) error {
	if method != c.expectedMethod {
		c.t.Errorf("Invalid method: %s", method)
	}
	if path != c.expectedPath {
		c.t.Errorf("Invalid path: %s", path)
	}
	for k, v := range headers {
		if v != c.expectedHeaders[k] {
			c.t.Errorf("Invalid header %s: %s", k, v)
		}
	}
	if len(headers) != len(c.expectedHeaders) {
		c.t.Errorf("Invalid headers count: %d", len(headers))
	}

	reqBody, err := json.Marshal(reqObj)
	if err != nil {
		c.t.Fatalf("Marshal error: %s", err)
	}

	expectedReqBody, err := json.Marshal(c.expectedReqObj)
	if err != nil {
		c.t.Fatalf("Marshal error: %s", err)
	}

	if string(reqBody) != string(expectedReqBody) {
		c.t.Errorf("Invalid request body: %s", string(reqBody))
	}

	if c.returnRespObj != nil {
		reflect.ValueOf(respObj).Elem().Set(reflect.ValueOf(c.returnRespObj).Elem())
	}

	return c.returnError
}

func TestNew(t *testing.T) {
	c := New(
		OptAppID("app_id"),
		OptEnv(EnvLive),
		OptPrivateKey("private_key"),
		OptHTTPClient(&httpClientMock{}),
	)

	if c == nil {
		t.Errorf("Client is nil")
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

func TestCall_WithApiResponse(t *testing.T) {
	httpClientMock := &httpClientMock{
		do: func(r *http.Request) (*http.Response, error) {
			if r.URL.String() != "https://api.paymentsos.com/somepath" {
				t.Errorf("Invalid request URI: %s", r.RequestURI)
			}
			if r.Method != "POST" {
				t.Errorf("Invalid request method: %s", r.Method)
			}
			if r.Header.Get(headerEnv) != string(EnvTest) {
				t.Errorf("Invalid request env: %s", r.Header.Get(headerEnv))
			}
			if r.Header.Get(headerAppID) != "app_id_test" {
				t.Errorf("Invalid request app ID: %s", r.Header.Get(headerAppID))
			}
			if r.Header.Get(headerPrivateKey) != "private_key_test" {
				t.Errorf("Invalid request private key: %s", r.Header.Get(headerPrivateKey))
			}
			if r.Header.Get("test-header") != "test-header-value" {
				t.Errorf("Invalid request custom header: %s", r.Header.Get("test-header"))
			}
			body, _ := ioutil.ReadAll(r.Body)
			if string(body) != `{"field":"request_value"}` {
				t.Errorf("Invalid request body: %s", string(body))
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"field":"response_value"}`)),
			}, nil
		},
	}

	req := request{
		Field: "request_value",
	}

	response := struct {
		Field string `json:"field"`
	}{}

	client := Client{
		httpClient: httpClientMock,
		appID:      "app_id_test",
		privateKey: "private_key_test",
		env:        EnvTest,
	}

	err := client.Call(
		context.Background(),
		"POST",
		"somepath",
		map[string]string{
			"test-header": "test-header-value",
		},
		&req,
		&response,
	)

	if err != nil {
		t.Errorf("Call returned error: %v", err)
	}

	if response.Field != "response_value" {
		t.Errorf("Response is invalid: %+v", response)
	}
}

func TestCall_WithApiError(t *testing.T) {
	httpClientMock := &httpClientMock{
		do: func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"category":"category_test"}`)),
			}, nil
		},
	}

	req := request{
		Field: "request_value",
	}

	client := Client{httpClient: httpClientMock}

	err := client.Call(
		context.Background(),
		"POST",
		"somepath",
		map[string]string{
			"test-header": "test-header-value",
		},
		&req,
		nil,
	)

	if err == nil {
		t.Error("Call didn't return error")
	}
	if zoozErr, ok := err.(*Error); ok {
		if zoozErr.StatusCode != http.StatusBadRequest {
			t.Errorf("Invalid error status code: %d", zoozErr.StatusCode)
		}
		if zoozErr.APIError.Category != "category_test" {
			t.Errorf("Invalid API error category: %s", zoozErr.APIError.Category)
		}
	} else {
		t.Errorf("Call return invalid error type: %T", err)
	}
}

func TestCall_WithTransportError(t *testing.T) {
	httpClientMock := &httpClientMock{
		do: func(r *http.Request) (*http.Response, error) {
			return nil, errors.New("do_error")
		},
	}

	req := request{
		Field: "request_value",
	}

	client := Client{httpClient: httpClientMock}

	err := client.Call(
		context.Background(),
		"POST",
		"somepath",
		map[string]string{
			"test-header": "test-header-value",
		},
		&req,
		nil,
	)

	if err == nil {
		t.Error("Call didn't return error")
	}
	if errors.Cause(err).Error() != "do_error" {
		t.Errorf("Invalid error cause: %v", errors.Cause(err))
	}
}

