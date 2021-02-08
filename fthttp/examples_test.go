package fthttp_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Financial-Times/go-ft-http/fthttp"
	"github.com/Financial-Times/go-logger/v2"
)

func ExampleNewClient() {
	log := logger.NewUPPInfoLogger("systemName")
	timeout := time.Second

	// a new client can be created by calling NewClient with desired options
	// see fthttp/client.go for full list of supported standard options
	client, err := fthttp.NewClient(
		fthttp.WithLogging(log),
		fthttp.WithTimeout(timeout),
		fthttp.WithSysInfo("PLATFORM", "system-code"))

	if err != nil {
		log.Fatal(err)
	}

	req, cleanup := getDummyRequest()
	defer cleanup()

	resp, _ := client.Do(req)
	defer resp.Body.Close() //nolint:govet

	fmt.Println(resp.Status)
	// Output: 404 Not Found
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
