package transport

import (
	"testing"
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggingRoundTripper_RoundTrip(t *testing.T) {

	var logBuffer bytes.Buffer

	logger := logrus.New()
	logger.Out = &logBuffer

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	loggingTransport := NewLoggingTransport(logger)

	request, _ := http.NewRequest("GET", server.URL, nil)

	_, err := loggingTransport.RoundTrip(request)

	assert.NoError(t, err)
	assert.NotEmpty(t, logBuffer.String())

}
