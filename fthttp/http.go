package fthttp

import (
	"net/http"
	"os"
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

// NewHttpClient returns an http client with defaults;
//  - 8 seconds timeout.
//  - Default Platform as PAC if no APP_PLATFORM env variable is defined.
//  - Default System code as unknown if no APP_SYSTEM_CODE env variable is defined.
func NewDefaultClient() *http.Client {

	appSystemCode, present := os.LookupEnv("APP_SYSTEM_CODE")

	if !present {
		appSystemCode = "Unknown"
	}

	appPlatform, present := os.LookupEnv("APP_PLATFORM")

	if !present {
		appPlatform = "PAC"
	}

	return NewClient(defaultClientTimeout, appPlatform, appSystemCode)
}
