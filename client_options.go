package httpclient

import (
	"net/http"
	"time"
)

// ClientOption defines the client option to customize the client
type ClientOption func(*Client)

// DisableRedirect disables to follow 3xx redirection
func DisableRedirect(client *Client) {
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}

// Timeout set the client request timeout
func Timeout(timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.Timeout = timeout
	}
}

// SetTransport set the transport of client
func SetTransport(transport http.RoundTripper) ClientOption {
	return func(client *Client) {
		client.Transport = NewLogTransport(client.ctx, transport)
	}
}

// SetCookieJar set the cookie jar of client
func SetCookieJar(cookieJar http.CookieJar) ClientOption {
	return func(client *Client) {
		client.Jar = cookieJar
	}
}
