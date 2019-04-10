package httpclient

import (
	"net"

	"github.com/eapache/go-resiliency/retrier"
)

// DefaultRetryClassifer is the default retry error classifer
var DefaultRetryClassifer = &retryClassifer{}

// retryClassifer defines the a retry strategy for network error
type retryClassifer struct{}

// Classify implements the retrier.Classifier interface
func (r *retryClassifer) Classify(err error) retrier.Action {
	if err == nil {
		return retrier.Succeed
	}

	if ne, ok := err.(net.Error); ok && ne.Temporary() {
		return retrier.Retry
	}

	return retrier.Fail
}
