package fthttp

import (
	"net/http"
	"time"

	"github.com/Financial-Times/go-ft-http/transport"
	"github.com/Sirupsen/logrus"
)

const defaultClientTimeout = 8 * time.Second

// NewHttpClient returns an http client with provided timeout and the FT specific transport used within.
// platform and systemCode are used to be construct the FT user-agent header.
func NewClient(timeout time.Duration, platform string, systemCode string) *http.Client {
	return &http.Client{
		Transport: transport.NewTransport().WithStandardUserAgent(platform, systemCode),
		Timeout:   timeout,
	}
}

// NewHttpClient returns an http client with provided timeout and the FT specific transport used within.
// platform and systemCode are used to be construct the FT user-agent header.
func NewLoggingClient(timeout time.Duration, platform string, systemCode string, logger *logrus.Logger) *http.Client {
	return &http.Client{
		Transport: transport.NewLoggingTransport(logger).WithStandardUserAgent(platform, systemCode),
		Timeout:   timeout,
	}
}

// NewClientWithDefaultTimeout returns and http client with the timeout already set.
func NewClientWithDefaultTimeout(platform string, systemCode string) *http.Client {
	return NewClient(defaultClientTimeout, platform, systemCode)
}

// NewClientWithDefaultTimeout returns and http client with the timeout already set.
func NewLoggingClientWithDefaultTimeout(platform string, systemCode string, logger *logrus.Logger) *http.Client {
	return NewLoggingClient(defaultClientTimeout, platform, systemCode, logger)
}
