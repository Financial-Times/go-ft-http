package transport

import (
	"net/http"
	"time"

	"github.com/Financial-Times/go-logger/v2"
	transactionidutils "github.com/Financial-Times/transactionid-utils-go"
)

type loggingRoundTripper struct {
	log     *logger.UPPLogger
	tripper http.RoundTripper
}

func (lrt *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	username := ""
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			username = name
		}
	}

	transactionID := transactionidutils.GetTransactionIDFromRequest(req)

	requestUri := "/"

	if req.URL.Path != "" {
		requestUri = req.URL.Path
	}

	t := time.Now()
	response, err := lrt.tripper.RoundTrip(req)
	elapsed := time.Since(t)

	withFields := lrt.log.WithFields(map[string]interface{}{
		"responsetime":   int64(elapsed.Seconds() * 1000),
		"username":       username,
		"method":         req.Method,
		"transaction_id": transactionID,
		"uri":            requestUri,
		"requestURL":     req.URL.String(),
		"protocol":       req.Proto,
		"referer":        req.Referer(),
		"userAgent":      req.UserAgent(),
	})

	if err == nil {
		withFields = withFields.WithField("status", response.Status)
	}

	withFields.Info()
	return response, err
}
