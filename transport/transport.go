package transport

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Financial-Times/service-status-go/buildinfo"
	tidutils "github.com/Financial-Times/transactionid-utils-go"
)

const (
	// DefaultTransactionIDHeaderName is the standard transaction id http header name
	DefaultTransactionIDHeaderName = "X-Request-Id"
	// DefaultTransactionIDContextValueKey is the standard transaction id context key name
	DefaultTransactionIDContextValueKey = "transaction_id"
)

// DelegatingTransport pre-processes requests with the configured Extensions, and then delegates to the provided http.RoundTripper implementation
type DelegatingTransport struct {
	delegate   http.RoundTripper
	extensions []Extension
}

// HeaderExtension adds the provided header if it has not already been set.
type HeaderExtension struct {
	header string
	value  string
}

// TIDFromContextExtension adds a transaction id request header if there is one available in the request.Context()
type TIDFromContextExtension struct{}

// Extension allows access to the request prior to it being executed against the delegated http.RoundTripper.
// Extensions MUST be side-effect free, and in general should not MODIFY the request so that it could produce an unintended response.
// For example, modifying the Host header could influence the outcome of the request, so should NOT be modified.
// An example of an acceptable modification is adding a suitable User-Agent header if one is not already set by the client.
type Extension interface {
	ExtendRequest(req *http.Request)
}

// NewTransport returns a delegating transport which uses the http.DefaultTransport
func NewTransport() *DelegatingTransport {
	return (&DelegatingTransport{delegate: http.DefaultTransport}).TransactionIDFromContext()
}

// NewUserAgentExtension creates a new HeaderExtension with the provided user agent value.
func NewUserAgentExtension(userAgent string) Extension {
	return &HeaderExtension{header: "User-Agent", value: userAgent}
}

// ExtendRequest adds the provided header if it has not already been set.
func (h *HeaderExtension) ExtendRequest(req *http.Request) {
	val := req.Header.Get(h.header)
	if val == "" {
		req.Header.Set(h.header, h.value)
	}
	log.Println(val)
}

// ExtendRequest retrieves the transaction_id from the http.Request.Context() and sets the corresponding X-Request-Id http.Header
func (h *TIDFromContextExtension) ExtendRequest(req *http.Request) {
	tid, err := tidutils.GetTransactionIDFromContext(req.Context())
	if err != nil {
		return
	}

	if header := req.Header.Get(DefaultTransactionIDHeaderName); header == "" {
		req.Header.Set(DefaultTransactionIDHeaderName, tid)
	}
}

// RoundTrip implementation will run the *http.Request against the configured extensions, and then delegate the request to the provided http.RoundTripper
func (d *DelegatingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header == nil {
		defer req.Body.Close()
		return nil, errors.New("http: nil Request.Header")
	}

	for _, e := range d.extensions {
		e.ExtendRequest(req)
	}

	return d.delegate.RoundTrip(req)
}

// WithUserAgent appends the provided value as the User-Agent header for all requests
func (d *DelegatingTransport) WithUserAgent(userAgent string) *DelegatingTransport {
	ext := NewUserAgentExtension(userAgent)
	d.extensions = append(d.extensions, ext)
	return d
}

// WithStandardUserAgent receives the platform and system code and appends a User-Agent header of `PLATFORM-system-code/x.x.x`. Version is retrieved from the buildinfo package.
func (d *DelegatingTransport) WithStandardUserAgent(platform string, systemCode string) *DelegatingTransport {
	ext := NewUserAgentExtension(standardUserAgent(platform, systemCode))
	d.extensions = append(d.extensions, ext)
	return d
}

// TransactionIDFromContext checks the request.Context() for a transaction id, and sets the corresponding X-Request-Id header if one is not already set.
func (d *DelegatingTransport) TransactionIDFromContext() *DelegatingTransport {
	d.extensions = append(d.extensions, &TIDFromContextExtension{})
	return d
}

func standardUserAgent(platform string, systemCode string) string {
	return strings.ToUpper(platform) + "-" + strings.ToLower(systemCode) + "/" + versionFromBuildInfo()
}

func versionFromBuildInfo() string {
	return strings.Replace(buildinfo.GetBuildInfo().Version, " ", "-", -1)
}
