package httpclient

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"context"

	"github.com/eapache/go-resiliency/retrier"
)

var (
	// DefaultTransport is the default transport to be used when not specified
	DefaultTransport = &http.Transport{
		MaxIdleConnsPerHost: 16,
	}

	// DefaultTimeout is the default client request timeout if not specified
	DefaultTimeout = 15 * time.Second
)

// Client is the http client handle
type Client struct {
	*http.Client
	Retry *Retry
	ctx   context.Context
}

// New creates a new http client with specified client options
func New(ctx context.Context, opts ...ClientOption) *Client {
	client := &Client{
		Client: &http.Client{},
		ctx:    ctx,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// Options sends the OPTIONS request
func (client *Client) Options(url, body string, reqOpts ...RequestOption) (result string, err error) {
	return client.Do("OPTIONS", url, body, reqOpts...)
}

// Head sends the HEAD request
func (client *Client) Head(url, body string, reqOpts ...RequestOption) (result string, err error) {
	return client.Do("HEAD", url, body, reqOpts...)
}

// Get sends the GET request
func (client *Client) Get(url, body string, reqOpts ...RequestOption) (result string, err error) {
	return client.Do("GET", url, body, reqOpts...)
}

// Post sends the POST request
func (client *Client) Post(url, body string, reqOpts ...RequestOption) (result string, err error) {
	return client.Do("POST", url, body, reqOpts...)
}

// Patch sends the PATCH request
func (client *Client) Patch(url, body string, reqOpts ...RequestOption) (result string, err error) {
	return client.Do("PATCH", url, body, reqOpts...)
}

// Put sends the PUT request
func (client *Client) Put(url, body string, reqOpts ...RequestOption) (result string, err error) {
	return client.Do("PUT", url, body, reqOpts...)
}

// Delete sends the DELETE request
func (client *Client) Delete(url, body string, reqOpts ...RequestOption) (result string, err error) {
	return client.Do("DELETE", url, body, reqOpts...)
}

// Do sends a custom METHOD request
func (client *Client) Do(method, url, body string, reqOpts ...RequestOption) (result string, err error) {
	if client.Retry == nil {
		return client.do(method, url, body, reqOpts...)
	}

	retry := retrier.New(client.Retry.BackOffs, client.Retry)

	err = retry.Run(func() error {
		if result, err = client.do(method, url, body, reqOpts...); err != nil {
			return err
		}
		return nil
	})

	return result, err
}

// do the internal request sending implementation
func (client *Client) do(method, url, body string, reqOpts ...RequestOption) (result string, err error) {
	var (
		req      *http.Request
		resp     *http.Response
		respData []byte
	)

	if req, err = http.NewRequest(method, url, strings.NewReader(body)); err != nil {
		logger.Error(client.ctx, "create http request",
			"http_method", method,
			"http_url", url,
			"http_body", body,
			"error", err,
		)
		return "", err
	}

	ctx := client.ctx

	for _, reqOpt := range reqOpts {
		if err = reqOpt(req); err != nil {
			logger.Error(client.ctx, "set request option",
				"http_method", method,
				"http_url", url,
				"http_body", body,
				"error", err,
			)
			return "", err
		}
	}

	if client.Transport == nil {
		client.Transport = DefaultTransport
	}

	if client.Timeout == 0 {
		client.Timeout = DefaultTimeout
	}

	begin := time.Now()
	resp, err = client.Client.Do(req)
	procTime := time.Since(begin)

	if err != nil {
		logger.Error(ctx, "do http request",
			"http_method", method,
			"http_url", url,
			"http_body", body,
			"error", err,
		)
		return "", err
	}
	// nolint: errcheck
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = &HTTPError{resp.StatusCode, resp.Status}
		logger.Error(ctx, "bad http status code",
			"http_method", method,
			"http_url", url,
			"http_body", body,
			"error", err,
		)
		return "", err
	}

	if respData, err = ioutil.ReadAll(resp.Body); err != nil {
		logger.Error(ctx, "read response body",
			"http_method", method,
			"http_url", url,
			"http_body", body,
			"error", err,
		)
		return "", err
	}

	result = string(respData)

	logger.Trace(ctx, "http call ok",
		"http_method", method,
		"http_url", url,
		"http_body", body,
		"result", result,
		"proc_time", procTime,
	)

	return result, nil
}
