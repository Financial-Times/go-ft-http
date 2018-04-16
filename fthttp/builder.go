package fthttp

import (
	"time"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/Financial-Times/go-ft-http/transport"
)

// NewClientBuilder provides an http client (step) builder implementation
// which includes both mandatory and optional steps
func NewClientBuilder() TimeoutStep {
	return &builder{}
}

// ClientBuilder mandatory finalizer step, embeds all possible/optional steps
type ClientBuilder interface {
	Build() *http.Client
	LoggingStep
}

// TimeoutStep mandatory, entry step
type TimeoutStep interface {
	WithTimeout(timeout time.Duration) SysInfoStep
}

// SysInfoStep mandatory, intermediate step
type SysInfoStep interface {
	WithSysInfo(platform string, systemCode string) ClientBuilder
}

// LoggingStep optional, intermediate step
type LoggingStep interface {
	WithLogging(logger *logrus.Logger) ClientBuilder
}

type builder struct {
	logger     *logrus.Logger
	timeout    time.Duration
	client     *http.Client
	platform   string
	systemCode string
}

func (cb *builder) WithLogging(logger *logrus.Logger) ClientBuilder {
	cb.logger = logger
	return cb
}

func (cb *builder) WithTimeout(timeout time.Duration) SysInfoStep {
	cb.timeout = timeout
	return cb
}

func (cb *builder) WithSysInfo(platform string, systemCode string) ClientBuilder {
	cb.platform = platform
	cb.systemCode = systemCode
	return cb
}

func (cb *builder) Build() *http.Client {

	if cb.client != nil {
		return cb.client
	}

	var dt *transport.DelegatingTransport

	if cb.logger != nil {
		dt = transport.NewLoggingTransport(cb.logger)
	} else {
		dt = transport.NewTransport()
	}

	dt = dt.WithStandardUserAgent(cb.platform, cb.systemCode)

	cb.client = &http.Client{
		Transport: dt,
		Timeout:   cb.timeout,
	}

	return cb.client
}
