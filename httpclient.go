package httpclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"context"

	"github.com/eapache/go-resiliency/retrier"
)

var (
	// DefaultTimeout is the default client request timeout if not specified
	DefaultTimeout  = 15 * time.Second
	DebugSetCookies = false
)

// Client is the http client handle
type Client struct {
	*http.Client
	Retry   *Retry
	reqOpts []RequestOption
	ctx     context.Context
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

// NewJSON return a JSON client wrapper
func (client *Client) NewJSON() *JSONClient {
	return &JSONClient{client}
}

// NewXML return a XML client wrapper
func (client *Client) NewXML() *XMLClient {
	return &XMLClient{client}
}

// SetDefaultReqOpts set the default request options, applied before each request.
func (client *Client) SetDefaultReqOpts(reqOpts ...RequestOption) {
	client.reqOpts = reqOpts[:len(reqOpts):len(reqOpts)]
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
			"method", method,
			"url", url,
			"body", body,
			"error", err,
		)
		return "", err
	}

	ctx := client.ctx

	reqOpts = append(client.reqOpts, reqOpts...)

	for _, reqOpt := range reqOpts {
		if err = reqOpt(req); err != nil {
			logger.Error(client.ctx, "set request option",
				"method", method,
				"url", url,
				"body", body,
				"error", err,
			)
			return "", err
		}
	}

	if client.Transport == nil {
		client.Transport = NewLogTransport(client.ctx, http.DefaultTransport)
	}

	if client.Timeout == 0 {
		client.Timeout = DefaultTimeout
	}

	begin := time.Now()
	resp, err = client.Client.Do(req)
	procTime := time.Since(begin)

	if err != nil {
		logger.Error(ctx, "do http request",
			"method", method,
			"url", url,
			"body", body,
			"error", err,
		)
		return "", err
	}
	// nolint: errcheck
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = &HTTPError{resp.StatusCode, resp.Status}
		logger.Error(ctx, "bad http status code",
			"method", method,
			"url", url,
			"body", body,
			"error", err,
		)
		return "", err
	}

	if respData, err = ioutil.ReadAll(resp.Body); err != nil {
		logger.Error(ctx, "read response body",
			"method", method,
			"url", url,
			"body", body,
			"error", err,
		)
		return "", err
	}

	result = string(respData)

	var kvs []interface{}
	if DebugSetCookies {
		buf := &bytes.Buffer{}
		for _, cookie := range resp.Cookies() {
			buf.WriteString(fmt.Sprintf("%v=%v|", cookie.Name, cookie.Value))
		}

		if buf.Len() > 0 {
			buf.Truncate(buf.Len() - 1)
		}
		kvs = []interface{}{
			"method", method,
			"url", url,
			"body", body,
			"result", result,
			"set_cookies", buf.String(),
			"proc_time", procTime,
		}
	} else {
		kvs = []interface{}{
			"method", method,
			"url", url,
			"body", body,
			"result", result,
			"proc_time", procTime,
		}
	}
	logger.Debug(ctx, "request success", kvs...)

	return result, nil
}
