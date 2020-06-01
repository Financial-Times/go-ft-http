package transport

import (
	"net/http"
	"strings"

	"github.com/Financial-Times/go-logger/v2"
	"github.com/Financial-Times/service-status-go/buildinfo"
)

// ExtensibleTransport pre-processes requests with the configured Extensions, and then delegates to the provided http.RoundTripper implementation
type ExtensibleTransport struct {
	delegate   http.RoundTripper
	extensions []HTTPRequestExtension
}

type Option func(d *ExtensibleTransport)

// NewTransport returns a delegating transport which uses the http.DefaultTransport
func NewTransport(options ...Option) *ExtensibleTransport {
	tr := &ExtensibleTransport{
		delegate: http.DefaultTransport,
		extensions: []HTTPRequestExtension{
			&TIDFromContextExtension{},
		},
	}
	tr.AddOptions(options...)
	return tr
}

// RoundTrip implementation will run the *http.Request against the configured extensions, and then delegate the request to the provided http.RoundTripper
func (d *ExtensibleTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, e := range d.extensions {
		e.ExtendRequest(req)
	}

	return d.delegate.RoundTrip(req)
}

// AddOptions provides a way to extends the current transport
func (d *ExtensibleTransport) AddOptions(options ...Option) {
	for _, opt := range options {
		opt(d)
	}
}

func (d *ExtensibleTransport) AddExtension(ext HTTPRequestExtension) {
	d.extensions = append(d.extensions, ext)
}

// NewLoggingTransport returns a delegating transport which creates log entries in the provided logger for every request.
// It adds TIDFromContextExtension to the request handling and uses the http.DefaultTransport as underlining transport
func WithLogger(log *logger.UPPLogger) Option {
	return func(d *ExtensibleTransport) {
		tr := d.delegate
		d.delegate = &loggingTransport{
			log:       log,
			transport: tr,
		}
	}
}

// WithUserAgent appends the provided value as the User-Agent header for all requests
func WithUserAgent(userAgent string) Option {
	return func(d *ExtensibleTransport) {
		ext := NewUserAgentExtension(userAgent)
		d.extensions = append(d.extensions, ext)
	}
}

// WithStandardUserAgent receives the platform and system code and appends a User-Agent header of `PLATFORM-system-code/x.x.x`. Version is retrieved from the buildinfo package.
func WithStandardUserAgent(platform string, systemCode string) Option {
	return func(d *ExtensibleTransport) {
		ext := NewUserAgentExtension(standardUserAgent(platform, systemCode))
		d.extensions = append(d.extensions, ext)
	}
}

func standardUserAgent(platform string, systemCode string) string {
	return removeWhitespace(strings.ToUpper(platform) + "-" + strings.ToLower(systemCode) + "/" + buildinfo.GetBuildInfo().Version)
}

func removeWhitespace(old string) string {
	return strings.Replace(old, " ", "-", -1)
}
