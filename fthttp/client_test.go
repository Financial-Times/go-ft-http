package fthttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Financial-Times/go-logger/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	hooks "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

const transactionIDKey string = "transaction_id"

func TestNewClientWithLogging(t *testing.T) {
	log := logger.NewUPPLogger("test", "info")
	h := hooks.NewLocal(log.Logger)
	defer h.Reset()

	client, err := NewClient(WithLogging(log))
	assert.NoError(t, err)

	srv := httptest.NewServer(http.NotFoundHandler())
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	assert.NoError(t, err)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	logEntry := h.LastEntry()
	msg, err := logEntry.String()
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"level":               "info",
		"method":              "GET",
		"protocol":            "HTTP/1.1",
		"requestURL":          srv.URL,
		"service_name":        "test",
		"status":              "404 Not Found",
		"uri":                 "/",
		transactionIDKey:      "ignored",
		logger.DefaultKeyTime: "ignored",
		"responsetime":        "ignored",
	}

	fields := map[string]interface{}{}
	err = json.Unmarshal([]byte(msg), &fields)
	assert.NoError(t, err)

	specialFields := map[string]bool{
		transactionIDKey:      true,
		logger.DefaultKeyTime: true,
		"responsetime":        true,
	}

	for key := range specialFields {
		_, ok := fields[key]
		assert.True(t, ok, "expect to have %s field in log", key)
	}

	ignoreFilter := cmpopts.IgnoreMapEntries(func(key string, val interface{}) bool {
		return specialFields[key]
	})

	if !cmp.Equal(expected, fields, ignoreFilter) {
		diff := cmp.Diff(expected, fields, ignoreFilter)
		t.Errorf("unexpected log output: %s", diff)
	}
}

func TestNewClientWithSysInfo(t *testing.T) {
	client, err := NewClient(WithSysInfo("TEST", "SystemCode"))
	assert.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		userAgent := req.Header.Get("User-Agent")
		// the version is initialized by our build system, so it's not setup correctly in the test
		assert.Equal(t, "TEST-systemcode/Version--is-not-a-semantic-version", userAgent)
	}))

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	assert.NoError(t, err)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	resp.Body.Close()
}

func TestNewClientWithUserAgent(t *testing.T) {
	client, err := NewClient(WithUserAgent("test-user-agent"))
	assert.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		userAgent := req.Header.Get("User-Agent")
		assert.Equal(t, "test-user-agent", userAgent)
	}))

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	assert.NoError(t, err)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	resp.Body.Close()
}

func TestNewClientWithTimeout(t *testing.T) {
	testTime := time.Millisecond * 200
	client, err := NewClient(WithTimeout(testTime))
	assert.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(testTime * 2)
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	assert.NoError(t, err)
	_, err = client.Do(req) // nolint:bodyclose
	urlErr := &url.Error{}
	assert.True(t, errors.As(err, &urlErr))
	assert.True(t, urlErr.Timeout())
}
