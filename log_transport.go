package httpclient

import (
	"context"
	"io/ioutil"
	"net/http"
)

type LogTransport struct {
	http.RoundTripper
	ctx context.Context
}

func NewLogTransport(ctx context.Context, transport http.RoundTripper) *LogTransport {
	return &LogTransport{
		RoundTripper: transport,
		ctx:          ctx,
	}
}

func (tr *LogTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, _ := req.GetBody()
	bodyContent, _ := ioutil.ReadAll(body)
	logger.Debug(tr.ctx, "do request", "method", req.Method, "url", req.URL.String(), "body", string(bodyContent))
	return tr.RoundTripper.RoundTrip(req)
}
