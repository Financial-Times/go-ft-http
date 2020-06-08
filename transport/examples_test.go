package transport_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Financial-Times/go-ft-http/transport"
	"github.com/Financial-Times/go-logger/v2"
	tidutils "github.com/Financial-Times/transactionid-utils-go"
)

func ExampleNewTransport() {
	log := logger.NewUPPInfoLogger("systemName")

	// create new ExtensibleTransport with the provided options
	// the options are applied in the order they are provided
	// ExtensibleTransport uses http.DefaultRoundTripper as underlining transport object
	// and adds TIDFromContextExtension by default
	rt := transport.NewTransport(transport.WithLogger(log))

	// options can be applied after ExtensibleTransport object is created
	rt.AddOptions(transport.WithUserAgent("userAgent"))

	// some options could operate on the same request functionality so be mindful ot the order you apply them
	// here the User-Agent header would already be set by the previous option so setting StandardUserAgent would fail silently.
	rt.AddOptions(transport.WithStandardUserAgent("PLATFORM", "system-code"))

	// new custom extensions could be added to ExtensibleTransport
	// IMPORTANT: Please read the documentation for http.RoundTripper before implementing new HttpRequestExtensions.
	rt.AddExtension(&authRequestExtension{usr: "user", pass: "pass"})

	client := http.Client{
		Transport: rt,
		Timeout:   time.Second,
	}

	req, cleanup := getDummyRequest()
	defer cleanup()

	ctx := tidutils.TransactionAwareContext(context.TODO(), "tid_1234")
	req = req.WithContext(ctx)

	resp, _ := client.Do(req)
	defer resp.Body.Close() //nolint:govet

	fmt.Printf("transaction_id: %s; user-agent: %s", req.Header.Get(tidutils.TransactionIDHeader), req.Header.Get("User-Agent"))
	// Output: transaction_id: tid_1234; user-agent: userAgent
}

type authRequestExtension struct {
	usr  string
	pass string
}

func (e *authRequestExtension) ExtendRequest(req *http.Request) {
	req.SetBasicAuth(e.usr, e.pass)
}

func getDummyRequest() (*http.Request, func()) {
	srv := httptest.NewServer(http.NotFoundHandler())

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		srv.Close()
		return nil, nil
	}
	return req, srv.Close
}
