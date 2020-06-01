package fthttp

import (
	"net/http"
	"time"

	"github.com/Financial-Times/go-ft-http/transport"
	"github.com/Financial-Times/go-logger/v2"
)

type Option func(c *config)

type config struct {
	logger     *logger.UPPLogger
	timeout    time.Duration
	platform   string
	systemCode string
	userAgent  string
}

// WithLogging instruments the client to start producing log entries for outgoing requests
func WithLogging(logger *logger.UPPLogger) Option {
	return func(c *config) {
		c.logger = logger
	}
}

// WithTimeout sets the http.Client Timeout to the provided duration.
func WithTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.timeout = timeout
	}
}

// WithSysInfo initializes the User-Agent header in a standard "{platform}-{systemCode}/Version-{version}" format
// When both `WithSysInfo` and `WithUserAgent` options are provided, `WithSysInfo` takes precedent.
func WithSysInfo(platform string, systemCode string) Option {
	return func(c *config) {
		c.systemCode = systemCode
		c.platform = platform
	}
}

// WithUserAgent initializes the User-Agent header with the provided string
// When both `WithSysInfo` and `WithUserAgent` options are provided, `WithSysInfo` takes precedent.
func WithUserAgent(user string) Option {
	return func(c *config) {
		c.userAgent = user
	}
}

// NewClient creates a http.Client object with the provided options
func NewClient(options ...Option) *http.Client {
	const defaultClientTimeout = 8 * time.Second
	c := &config{
		timeout: defaultClientTimeout,
	}

	for _, fn := range options {
		fn(c)
	}
	ops := make([]transport.Option, 0)
	if c.logger != nil {
		ops = append(ops, transport.WithLogger(c.logger))
	}
	if c.platform != "" {
		ops = append(ops, transport.WithStandardUserAgent(c.platform, c.systemCode))
	} else if c.userAgent != "" {
		ops = append(ops, transport.WithUserAgent(c.userAgent))
	}

	dt := transport.NewTransport(ops...)

	client := &http.Client{
		Transport: dt,
		Timeout:   c.timeout,
	}

	return client
}
