package logging

import (
	"github.com/Financial-Times/transactionid-utils-go"
	"github.com/Sirupsen/logrus"
	"net/http"
	"time"
)

type RoundTripper struct {
	L  *logrus.Logger
	Rt http.RoundTripper
}

func (lrt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	username := "-"
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			username = name
		}
	}

	transactionID := ""
	trxId := req.Header.Get(transactionidutils.TransactionIDHeader)

	if trxId != "" {
		transactionID = trxId
	}

	requestUri := "/"

	if req.URL.Path != "" {
		requestUri = req.URL.Path
	}

	t := time.Now()
	response, err := lrt.Rt.RoundTrip(req)
	elapsed := time.Since(t)

	withFields := lrt.L.WithFields(logrus.Fields{
		"responsetime":   int64(elapsed.Seconds() * 1000),
		"username":       username,
		"method":         req.Method,
		"transaction_id": transactionID,
		"uri":            requestUri,
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
