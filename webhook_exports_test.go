package zooz

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func WebhookRequestSignature(t *testing.T, keyProvider privateKeyProvider, reqBody []byte, reqHeader http.Header) string {
	signature, err := calculateWebhookSignature(keyProvider, reqBody, reqHeader)
	require.NoError(t, err)
	return "sig1=" + signature
}

func CalculateWebhookSignature(keyProvider privateKeyProvider, reqBody []byte, reqHeader http.Header) (string, error) {
	return calculateWebhookSignature(keyProvider, reqBody, reqHeader)
}
