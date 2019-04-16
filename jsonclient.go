package httpclient

import (
	"context"
	"encoding/json"

	"github.com/k81/log"
)

// JSONClient is an wrapper of *Client, which talks in JSON
type JSONClient struct {
	*Client
}

// NewJSON create a JSON http client instance with specified options
func NewJSON(opts ...ClientOption) *JSONClient {
	client := New(opts...)
	return &JSONClient{client}
}

// Options sends the OPTIONS request
func (client *JSONClient) Options(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "OPTIONS", url, body, result, reqOpts...)
}

// Head sends the HEAD request
func (client *JSONClient) Head(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "HEAD", url, body, result, reqOpts...)
}

// Get sends the GET request
func (client *JSONClient) Get(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "GET", url, body, result, reqOpts...)
}

// Post sends the POST request
func (client *JSONClient) Post(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "POST", url, body, result, reqOpts...)
}

// Patch sends the PATCH request
func (client *JSONClient) Patch(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "PATCH", url, body, result, reqOpts...)
}

// Put sends the PUT request
func (client *JSONClient) Put(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "PUT", url, body, result, reqOpts...)
}

// Delete sends the DELETE request
func (client *JSONClient) Delete(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "DELETE", url, body, result, reqOpts...)
}

// Do sends a custom METHOD request
func (client *JSONClient) Do(ctx context.Context, method, url string, body, result interface{}, reqOpts ...RequestOption) error {
	var (
		bodyData  []byte
		resultStr string
		err       error
	)

	if body != nil {
		switch bodyValue := body.(type) {
		case string:
			bodyData = []byte(bodyValue)
		case json.RawMessage:
			bodyData = []byte(bodyValue)
		case []byte:
			bodyData = bodyValue
		default:
			if bodyData, err = json.Marshal(body); err != nil {
				log.Error(ctx, "marshal request body", "error", err)
				return err
			}
		}
	}

	reqOpts = append([]RequestOption{SetTypeJSON()}, reqOpts...)

	if resultStr, err = client.Client.Do(ctx, method, url, string(bodyData), reqOpts...); err != nil {
		return err
	}

	if result != nil && resultStr != "" {
		if err = json.Unmarshal([]byte(resultStr), result); err != nil {
			log.Error(ctx, "unmarshal response body", "error", err)
			return err
		}
	}
	return nil
}
