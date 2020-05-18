package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	transactionidutils "github.com/Financial-Times/transactionid-utils-go"

	"github.com/Financial-Times/go-logger/v2"
	"github.com/stretchr/testify/assert"
)

func TestLoggingRoundTripper_RoundTrip(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tests := map[string]struct {
		expected string
		method   string
		url      string
		headers  map[string]string
	}{
		"Minimum logging data": {
			method:   http.MethodGet,
			url:      "/testing",
			expected: `{"level":"info", "method":"GET", "protocol":"HTTP/1.1", "service_name":"testSystemCode", "status":"200 OK", "uri":"/testing"}`,
		},
		"Extract data from headers": {
			method: http.MethodPost,
			url:    "/testing",
			headers: map[string]string{
				"Referer":      "http://example.com",
				"User-Agent":   "User agent",
				"X-Request-Id": "KnownTransactionId",
			},
			expected: `{"level":"info", "method":"POST", "protocol":"HTTP/1.1", "referer":"http://example.com", "service_name":"testSystemCode", "status":"200 OK", "transaction_id":"KnownTransactionId", "uri":"/testing", "userAgent":"User agent"}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var logBuffer bytes.Buffer

			log := logger.NewUPPInfoLogger("testSystemCode")
			log.Out = &logBuffer

			loggingTransport := NewLoggingTransport(log)

			req, _ := http.NewRequest(test.method, server.URL+test.url, nil)
			for key, val := range test.headers {
				req.Header.Set(key, val)
			}

			resp, err := loggingTransport.RoundTrip(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			fields := map[string]interface{}{}
			err = json.Unmarshal(logBuffer.Bytes(), &fields)
			assert.NoError(t, err)

			_, ok := fields[logger.DefaultKeyTime]
			assert.True(t, ok)
			delete(fields, logger.DefaultKeyTime)

			_, ok = fields[logger.DefaultKeyTransactionID]
			assert.True(t, ok)
			if _, ok = test.headers[transactionidutils.TransactionIDHeader]; !ok {
				delete(fields, logger.DefaultKeyTransactionID)
			}

			respTime, ok := fields["responsetime"]
			assert.True(t, ok, "Missing responsetime in the logs")
			assert.InDelta(t, 1, respTime, 10)
			delete(fields, "responsetime")

			data, err := json.Marshal(fields)
			assert.NoError(t, err)
			assert.JSONEq(t, test.expected, string(data))
		})
	}

}
