package transport

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Financial-Times/go-logger/v2"
	"github.com/Financial-Times/service-status-go/buildinfo"
)

// HTTPRequestExtension allows access to the request prior to it being executed against the delegated http.RoundTripper.
// IMPORTANT: Please read the documentation for http.RoundTripper before implementing new HttpRequestExtensions.
type HTTPRequestExtension interface {
	ExtendRequest(req *http.Request)
}

// DelegatingTransport pre-processes requests with the configured Extensions, and then delegates to the provided http.RoundTripper implementation
type DelegatingTransport struct {
	delegate   http.RoundTripper
	extensions []HTTPRequestExtension
}

type DelegateOpt func(d *DelegatingTransport)

// NewTransport returns a delegating transport which uses the http.DefaultTransport
func NewTransport(options ...DelegateOpt) *DelegatingTransport {
	tr := &DelegatingTransport{
		delegate: http.DefaultTransport,
		extensions: []HTTPRequestExtension{
			&TIDFromContextExtension{},
		},
	}
	for _, opt := range options {
		opt(tr)
	}
	return tr
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

// NewLoggingTransport returns a delegating transport which creates log entries in the provided logger for every request.
// It adds TIDFromContextExtension to the request handling and uses the http.DefaultTransport as underlining round tripper
func WithLogger(log *logger.UPPLogger) DelegateOpt {
	return func(d *DelegatingTransport) {
		tr := d.delegate
		d.delegate = &loggingRoundTripper{
			log:     log,
			tripper: tr,
		}
	}
}

// WithUserAgent appends the provided value as the User-Agent header for all requests
func WithUserAgent(userAgent string) DelegateOpt {
	return func(d *DelegatingTransport) {
		ext := NewUserAgentExtension(userAgent)
		d.extensions = append(d.extensions, ext)
	}
}

// WithStandardUserAgent receives the platform and system code and appends a User-Agent header of `PLATFORM-system-code/x.x.x`. Version is retrieved from the buildinfo package.
func WithStandardUserAgent(platform string, systemCode string) DelegateOpt {
	return func(d *DelegatingTransport) {
		ext := NewUserAgentExtension(standardUserAgent(platform, systemCode))
		d.extensions = append(d.extensions, ext)
	}
}

// WithTransactionIDFromContext checks the request.Context() for a transaction id, and sets the corresponding X-Request-Id header if one is not already set.
func WithTransactionIDFromContext() DelegateOpt {
	return func(d *DelegatingTransport) {
		d.extensions = append(d.extensions, &TIDFromContextExtension{})
	}
}

func standardUserAgent(platform string, systemCode string) string {
	return removeWhitespace(strings.ToUpper(platform) + "-" + strings.ToLower(systemCode) + "/" + buildinfo.GetBuildInfo().Version)
}

func removeWhitespace(old string) string {
	return strings.Replace(old, " ", "-", -1)
}
