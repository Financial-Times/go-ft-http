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

func WithLogging(logger *logger.UPPLogger) Option {
	return func(c *config) {
		c.logger = logger
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.timeout = timeout
	}
}

func WithSysInfo(platform string, systemCode string) Option {
	return func(c *config) {
		c.systemCode = systemCode
		c.platform = platform
	}
}

func WithUserAgent(user string) Option {
	return func(c *config) {
		c.userAgent = user
	}
}

func NewClient(options ...Option) *http.Client {
	const defaultClientTimeout = 8 * time.Second
	c := &config{
		timeout: defaultClientTimeout,
	}

	for _, fn := range options {
		fn(c)
	}
	ops := make([]transport.DelegateOpt, 0)
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
