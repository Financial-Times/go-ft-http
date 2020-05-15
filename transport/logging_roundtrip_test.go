package transport

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Financial-Times/go-logger/v2"
	"github.com/stretchr/testify/assert"
)

func TestLoggingRoundTripper_RoundTrip(t *testing.T) {

	var logBuffer bytes.Buffer

	log := logger.NewUPPInfoLogger("testSystemCode")
	log.Out = &logBuffer

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	loggingTransport := NewLoggingTransport(log)

	request, _ := http.NewRequest("GET", server.URL, nil)

	_, err := loggingTransport.RoundTrip(request)

	assert.NoError(t, err)
	assert.NotEmpty(t, logBuffer.String())

}
