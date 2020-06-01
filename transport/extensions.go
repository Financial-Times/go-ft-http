package transport

import (
	"net/http"

	tidutils "github.com/Financial-Times/transactionid-utils-go"
)

// HeaderExtension adds the provided header if it has not already been set.
type HeaderExtension struct {
	header string
	value  string
}

// HTTPRequestExtension allows access to the request prior to it being executed against the delegated http.RoundTripper.
// IMPORTANT: Please read the documentation for http.RoundTripper before implementing new HttpRequestExtensions.
type HTTPRequestExtension interface {
	ExtendRequest(req *http.Request)
}

// NewUserAgentExtension creates a new HeaderExtension with the provided user agent value.
func NewUserAgentExtension(userAgent string) HTTPRequestExtension {
	return &HeaderExtension{header: "User-Agent", value: userAgent}
}

// ExtendRequest adds the provided header if it has not already been set.
func (h *HeaderExtension) ExtendRequest(req *http.Request) {
	val := req.Header.Get(h.header)
	if val == "" {
		req.Header.Set(h.header, h.value)
	}
}

// TIDFromContextExtension adds a transaction id request header if there is one available in the request.Context()
type TIDFromContextExtension struct{}

// ExtendRequest retrieves the transaction_id from the http.Request.Context() and sets the corresponding X-Request-Id http.Header
func (h *TIDFromContextExtension) ExtendRequest(req *http.Request) {
	tid, err := tidutils.GetTransactionIDFromContext(req.Context())
	if err != nil {
		return
	}

	if header := req.Header.Get(tidutils.TransactionIDHeader); header == "" {
		req.Header.Set(tidutils.TransactionIDHeader, tid)
	}
}
