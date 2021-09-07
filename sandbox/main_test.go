package sandbox

import (
	"fmt"
	"github.com/gtforge/go-zooz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"text/template"
	"time"
)

const (
	EnvTest        = "PAYMENTSOS_SANDBOX_TEST"
	EnvAppID       = "PAYMENTSOS_SANDBOX_APP_ID"
	EnvPrivateKey  = "PAYMENTSOS_SANDBOX_PRIVATE_KEY"
	EnvPublicKey   = "PAYMENTSOS_SANDBOX_PUBLIC_KEY"
	EnvLogRequests = "PAYMENTSOS_SANDBOX_LOG_REQUESTS"
)

const UnknownUUID = "00000000-0000-1000-8000-000000000000"

func TestMain(m *testing.M) {
	runTests := os.Getenv(EnvTest)
	if runTests == "" {
		println("Skip tests with paymentsos sandbox.")
		println("Add env " + EnvTest + "=1 to run them.")
		println("Also following envs are required: " + EnvAppID + ", " + EnvPrivateKey + ", " + EnvPublicKey)
		os.Exit(0)
	}
	rand.Seed(time.Now().UnixNano())
	os.Exit(m.Run())
}

var LogTemplate = template.Must(template.New("").Parse(`
==> HTTP Request:
		==> {{ .Method }} {{ .URL }}
		{{- range $k, $v := .RequestHeader }}
			{{ $k -}}: {{ range $v -}} {{- . -}} {{- end -}}
		{{- end }}
			{{- if .RequestBody }}
			Body:
				 {{ .RequestBody }}
			{{- end }}
		==> END {{ .Method }}
{{ if .Error -}}
<== ERROR ({{ .Duration }}):
		{{ .Error }}
{{- else -}}
<== HTTP Response:
		<== {{ .URL }} ({{ .Duration }})
			Status-Code: {{ .Status }}
		{{- range $k, $v := .ResponseHeader }}
			{{ $k -}}: {{ range $v -}} {{- . -}} {{- end -}}
		{{- end }}
			{{- if .ResponseBody }}
			Body:
				 {{ .ResponseBody }}
			{{- end }}
		<== END HTTP
{{- end }}
=== END ===`))

type Log struct {
	Method         string
	URL            string
	RequestHeader  http.Header
	RequestBody    string
	Duration       time.Duration
	Status         string
	ResponseHeader http.Header
	ResponseBody   string
	Error          string
}

type LoggingHTTPClient struct {
	next zooz.HTTPClient
	t    *testing.T
}

func (c LoggingHTTPClient) Do(request *http.Request) (*http.Response, error) {
	log := Log{
		Method:        request.Method,
		URL:           request.URL.String(),
		RequestHeader: request.Header,
	}
	defer func() {
		builder := &strings.Builder{}
		err := LogTemplate.Execute(builder, log)
		assert.NoError(c.t, err)
		c.t.Log(builder.String())
	}()

	if request.Body != nil {
		reqBody, err := io.ReadAll(request.Body)
		require.NoError(c.t, err)
		require.NoError(c.t, request.Body.Close())
		log.RequestBody = string(reqBody)
		request.Body = io.NopCloser(strings.NewReader(string(reqBody)))
	}

	start := time.Now()
	response, err := c.next.Do(request)
	log.Duration = time.Since(start)
	if err != nil {
		log.Error = fmt.Sprintf("%+v", err)
		return nil, err
	}

	log.Status = response.Status
	log.ResponseHeader = response.Header

	if response.Body != nil {
		respBody, err := io.ReadAll(response.Body)
		require.NoError(c.t, err)
		require.NoError(c.t, response.Body.Close())
		log.ResponseBody = string(respBody)
		response.Body = io.NopCloser(strings.NewReader(string(respBody)))
	}

	return response, nil
}

func GetClient(t *testing.T) *zooz.Client {
	var errs []string

	appID := os.Getenv(EnvAppID)
	if appID == "" {
		errs = append(errs, "missing env "+EnvAppID)
	}
	privateKey := os.Getenv(EnvPrivateKey)
	if privateKey == "" {
		errs = append(errs, "missing env "+EnvPrivateKey)
	}
	publicKey := os.Getenv(EnvPublicKey)
	if publicKey == "" {
		errs = append(errs, "missing env "+EnvPublicKey)
	}

	require.Empty(t, errs)

	var httpClient zooz.HTTPClient = http.DefaultClient
	if os.Getenv(EnvLogRequests) != "" {
		httpClient = LoggingHTTPClient{
			next: httpClient,
			t:    t,
		}
	}

	return zooz.New(
		zooz.OptHTTPClient(httpClient),
		zooz.OptAppID(appID),
		zooz.OptEnv(zooz.EnvTest),
		zooz.OptPrivateKey(privateKey),
		zooz.OptPublicKey(publicKey),
	)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = randChar()
	}
	return string(bytes)
}

// randChar returns random character in range 0-9 or A-Z or a-z
func randChar() byte {
	r := rand.Intn(10 + 26 + 26)
	if r < 10 {
		// 48 - 57 -> digits
		return byte(48 + r)
	}

	r -= 10
	if r < 26 {
		// 65 - 90 -> capital letters
		return byte(65 + r)
	}

	r -= 26
	// 97 - 122 -> letters
	return byte(97 + r)
}

// must is a helper function to implement soft assertions from Java world.
// The idea is to do several assertions at once:
//  * All assertions inside do func will be executed no matter what.
//  * But the test will be stopped just after if any of them failed.
// Note: Do not use require package in do func, use assert package instead.
// Example:
// 	must(t, func() {
// 		assert.Equal(t, ...)
// 		assert.Equal(t, ...)
// 		assert.Equal(t, ...)
// 	})
func must(t *testing.T, do func()) {
	do()
	if t.Failed() {
		t.FailNow()
	}
}

func requireZoozError(t *testing.T, err error, statusCode int, expected zooz.APIError) {
	zoozErr := &zooz.Error{}
	require.ErrorAs(t, err, &zoozErr)
	require.Equal(t, &zooz.Error{
		StatusCode: statusCode,
		RequestID:  zoozErr.RequestID, // ignore
		APIError:   expected,
	}, zoozErr)
}

func normalizeExpirationDate(date string) string {
	xxx := regexp.MustCompile(`^(0[1-9]|1[0-2])\D(\d{2,4})$`).FindStringSubmatch(date)
	month, year := xxx[1], xxx[2]
	if len(year) == 2 {
		year = "20" + year
	}
	return month + "/" + year
}

func last4(cardNumber string) string {
	return cardNumber[len(cardNumber)-4:]
}
