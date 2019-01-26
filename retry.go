package httpclient

import (
	"net"
	"time"

	"github.com/eapache/go-resiliency/retrier"
)

// Retry defines the retry strategy
type Retry struct {
	BackOffs []time.Duration
}

// Classify implements the retrier.Classifier interface
func (r *Retry) Classify(err error) retrier.Action {
	if err == nil {
		return retrier.Succeed
	}

	if ne, ok := err.(net.Error); ok && ne.Temporary() {
		return retrier.Retry
	}

	return retrier.Fail
}
