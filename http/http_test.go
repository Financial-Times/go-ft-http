package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testPlatform                           = "C64"
	testSystemCode                         = "an-awesome-service-as-usual"
	expectedUserAgent                      = "PAC-an-awesome-service-as-usual/Version--is-not-a-semantic-version"
	expectedUserAgentWithUnknownSystemCode = "PAC-unknown/Version--is-not-a-semantic-version"
)

func TestNewClientTimeoutSetting(t *testing.T) {
	client := NewClient(time.Second, testPlatform, testSystemCode)
	assert.Equal(t, time.Second, client.Timeout)
}

func TestNewDefaultClientSettings(t *testing.T) {

	os.Setenv("APP_SYSTEM_CODE", testSystemCode)
	defer os.Unsetenv("APP_SYSTEM_CODE")
	testServer := newHttpTestServer(defaultSettingsTestHandler{false, t})

	defaultClient := NewDefaultClient()
	assert.Equal(t, defaultClientTimeout, defaultClient.Timeout)

	request, _ := http.NewRequest(http.MethodGet, testServer.URL, nil)
	response, _ := defaultClient.Do(request)

	response.Body.Close()
}

func TestNewDefaultClientSettingsForUnknownSystemCode(t *testing.T) {

	testServer := newHttpTestServer(defaultSettingsTestHandler{true, t})

	defaultClient := NewDefaultClient()
	request, _ := http.NewRequest(http.MethodGet, testServer.URL, nil)
	response, _ := defaultClient.Do(request)

	response.Body.Close()
}

func newHttpTestServer(d defaultSettingsTestHandler) *httptest.Server {
	return httptest.NewServer(d)
}

type defaultSettingsTestHandler struct {
	unknownSystemCode bool
	t                 *testing.T
}

func (d defaultSettingsTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// validate defaults in Request
	userAgent := r.Header.Get("User-Agent")

	if d.unknownSystemCode {
		assert.Equal(d.t, expectedUserAgentWithUnknownSystemCode, userAgent)
		return
	}

	assert.Equal(d.t, expectedUserAgent, userAgent)

}
