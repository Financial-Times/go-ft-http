package transport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testHandler struct {
	t                   *testing.T
	userAgent           *string
	expectUserAgent     bool
	transactionID       *string
	expectTransactionID bool
}

func newTestHandler(t *testing.T, userAgent *string, transactionID *string) *testHandler {
	return &testHandler{
		t:                   t,
		userAgent:           userAgent,
		expectUserAgent:     userAgent != nil,
		transactionID:       transactionID,
		expectTransactionID: transactionID != nil,
	}
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.expectUserAgent {
		actual := r.Header.Get("User-Agent")
		assert.Equal(h.t, *h.userAgent, actual)
	}

	if h.expectTransactionID {
		actual := r.Header.Get(DefaultTransactionIDHeaderName)
		assert.Equal(h.t, *h.transactionID, actual)
	}
}

func TestUserAgent(t *testing.T) {
	testUserAgent := "PAC/blah"
	d := NewTransport().WithUserAgent(testUserAgent)

	c := http.Client{Transport: d}
	h := newTestHandler(t, &testUserAgent, nil)

	srv := httptest.NewServer(h)

	req, err := http.NewRequest("GET", srv.URL, nil)
	require.NoError(t, err)

	resp, err := c.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUserAgentIsNotOverridden(t *testing.T) {
	testUserAgent := "EXPECTED/found"
	d := NewTransport().WithUserAgent("NOT/found")

	c := http.Client{Transport: d}
	h := newTestHandler(t, &testUserAgent, nil)

	srv := httptest.NewServer(h)

	req, err := http.NewRequest("GET", srv.URL, nil)
	require.NoError(t, err)

	req.Header.Add("User-Agent", testUserAgent)

	resp, err := c.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTransactionIdFromContext(t *testing.T) {
	testTransactionID := "tid_testtttt"
	ctx := context.WithValue(context.Background(), DefaultTransactionIDContextValueKey, testTransactionID)

	d := NewTransport().TransactionIDFromContext()

	c := http.Client{Transport: d}
	h := newTestHandler(t, nil, &testTransactionID)

	srv := httptest.NewServer(h)

	req, err := http.NewRequest("GET", srv.URL, nil)
	require.NoError(t, err)

	resp, err := c.Do(req.WithContext(ctx))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTransactionIdFromContextNoValueInContext(t *testing.T) {
	d := NewTransport().TransactionIDFromContext()

	c := http.Client{Transport: d}
	h := &testHandler{t: t}

	srv := httptest.NewServer(h)

	req, err := http.NewRequest("GET", srv.URL, nil)
	require.NoError(t, err)

	resp, err := c.Do(req.WithContext(context.Background()))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestStandardUserAgent(t *testing.T) {
	testUserAgent := "PAC-example-system-code/Version--is-not-a-semantic-version" // Version--is-not-a-semantic-version is the default version returned from buildinfo. We can't influence the version without building with ldflags, so this will have to do.
	d := NewTransport().WithStandardUserAgent("PAC", "example-system-code")

	c := http.Client{Transport: d}
	h := newTestHandler(t, &testUserAgent, nil)

	srv := httptest.NewServer(h)

	req, err := http.NewRequest("GET", srv.URL, nil)
	require.NoError(t, err)

	resp, err := c.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
