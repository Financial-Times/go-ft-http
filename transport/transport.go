package transport

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Financial-Times/service-status-go/buildinfo"
	tidutils "github.com/Financial-Times/transactionid-utils-go"
)

// DelegatingTransport pre-processes requests with the configured Extensions, and then delegates to the provided http.RoundTripper implementation
type DelegatingTransport struct {
	delegate   http.RoundTripper
	extensions []HttpRequestExtension
}

// HeaderExtension adds the provided header if it has not already been set.
type HeaderExtension struct {
	header string
	value  string
}

// TIDFromContextExtension adds a transaction id request header if there is one available in the request.Context()
type TIDFromContextExtension struct{}

// HttpRequestExtension allows access to the request prior to it being executed against the delegated http.RoundTripper.
// IMPORTANT: Please read the documentation for http.RoundTripper before implementing new HttpRequestExtensions.
type HttpRequestExtension interface {
	ExtendRequest(req *http.Request)
}

// NewTransport returns a delegating transport which uses the http.DefaultTransport
func NewTransport() *DelegatingTransport {
	return (&DelegatingTransport{delegate: http.DefaultTransport}).WithTransactionIDFromContext()
}

// NewUserAgentExtension creates a new HeaderExtension with the provided user agent value.
func NewUserAgentExtension(userAgent string) HttpRequestExtension {
	return &HeaderExtension{header: "User-Agent", value: userAgent}
}

// ExtendRequest adds the provided header if it has not already been set.
func (h *HeaderExtension) ExtendRequest(req *http.Request) {
	val := req.Header.Get(h.header)
	if val == "" {
		req.Header.Set(h.header, h.value)
	}
}

// ExtendRequest retrieves the transaction_id from the http.Request.Context() and sets the corresponding X-Request-Id http.Header
func (h *TIDFromContextExtension) ExtendRequest(req *http.Request) {
	tid, err := tidutils.GetTransactionIDFromContext(req.Context())
	if err != nil {
		return
	}

	if header := req.Header.Get(tidutils.TransactionIDHeader); header == "" {
		req.Header.Set(tidutils.TransactionIDHeader, tid)
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

// WithTransactionIDFromContext checks the request.Context() for a transaction id, and sets the corresponding X-Request-Id header if one is not already set.
func (d *DelegatingTransport) WithTransactionIDFromContext() *DelegatingTransport {
	d.extensions = append(d.extensions, &TIDFromContextExtension{})
	return d
}

func standardUserAgent(platform string, systemCode string) string {
	return removeWhitespace(strings.ToUpper(platform) + "-" + strings.ToLower(systemCode) + "/" + buildinfo.GetBuildInfo().Version)
}

func removeWhitespace(old string) string {
	return strings.Replace(old, " ", "-", -1)
}
