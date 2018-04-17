package fthttp

import (
	"net/http"
	"time"

)

const defaultClientTimeout = 8 * time.Second

// NewHttpClient returns an http client with provided timeout and the FT specific transport used within.
// platform and systemCode are used to be construct the FT user-agent header.
//
// @Deprecated, use the NewClientBuilder instead
func NewClient(timeout time.Duration, platform string, systemCode string) *http.Client {
	return NewClientBuilder().
		WithTimeout(timeout).
		WithSysInfo(platform, systemCode).
		Build()
}

// NewClientWithDefaultTimeout returns and http client with the timeout already set.
//
// @Deprecated, use the NewClientBuilder instead
func NewClientWithDefaultTimeout(platform string, systemCode string) *http.Client {
	return NewClientBuilder().
		WithTimeout(defaultClientTimeout).
		WithSysInfo(platform, systemCode).
		Build()
}

