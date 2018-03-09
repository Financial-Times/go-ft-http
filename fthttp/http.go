package fthttp

import (
	"net/http"
	"time"

	"github.com/Financial-Times/go-ft-http/transport"
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

// NewClientWithDefaultTimeout returns and http client with the timeout already set.
func NewClientWithDefaultTimeout(platform string, systemCode string) *http.Client {
	return NewClient(defaultClientTimeout, platform, systemCode)
}
