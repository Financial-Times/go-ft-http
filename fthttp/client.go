package fthttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/Financial-Times/go-ft-http/transport"
	"github.com/Financial-Times/go-logger/v2"
)

var ErrWrongTransport = errors.New("expected ExtensibleTransport round tripper")

type Option func(c *http.Client) error

// WithLogging instruments the client to start producing log entries for outgoing requests
func WithLogging(logger *logger.UPPLogger) Option {
	return func(c *http.Client) error {
		tr, ok := c.Transport.(*transport.ExtensibleTransport)
		if !ok {
			return ErrWrongTransport
		}
		tr.AddOptions(transport.WithLogger(logger))
		return nil
	}
}

// WithTimeout sets the http.Client Timeout to the provided duration.
func WithTimeout(timeout time.Duration) Option {
	return func(c *http.Client) error {
		c.Timeout = timeout
		return nil
	}
}

// WithSysInfo initializes the User-Agent header in a standard "{platform}-{systemCode}/Version-{version}" format
// When both `WithSysInfo` and `WithUserAgent` options are provided, `WithSysInfo` takes precedent.
func WithSysInfo(platform string, systemCode string) Option {
	return func(c *http.Client) error {
		tr, ok := c.Transport.(*transport.ExtensibleTransport)
		if !ok {
			return ErrWrongTransport
		}
		tr.AddOptions(transport.WithStandardUserAgent(platform, systemCode))
		return nil
	}
}

// WithUserAgent initializes the User-Agent header with the provided string
// When both `WithSysInfo` and `WithUserAgent` options are provided, `WithSysInfo` takes precedent.
func WithUserAgent(user string) Option {
	return func(c *http.Client) error {
		tr, ok := c.Transport.(*transport.ExtensibleTransport)
		if !ok {
			return ErrWrongTransport
		}
		tr.AddOptions(transport.WithUserAgent(user))
		return nil
	}
}

// NewClient creates a http.Client object with the provided options
func NewClient(options ...Option) (*http.Client, error) {
	const defaultClientTimeout = 8 * time.Second
	client := &http.Client{
		Transport: transport.NewTransport(),
		Timeout:   defaultClientTimeout,
	}

	for _, opt := range options {
		err := opt(client)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}
