package httpclient

import (
	"context"
	"encoding/json"
)

// JSONClient is an wrapper of *Client, which talks in JSON
type JSONClient struct {
	*Client
}

// NewJSON create a JSON http client instance with specified options
func NewJSON(ctx context.Context, opts ...ClientOption) *JSONClient {
	client := New(ctx, opts...)
	return &JSONClient{client}
}

// Options sends the OPTIONS request
func (client *JSONClient) Options(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("OPTIONS", url, body, result, reqOpts...)
}

// Head sends the HEAD request
func (client *JSONClient) Head(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("HEAD", url, body, result, reqOpts...)
}

// Get sends the GET request
func (client *JSONClient) Get(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("GET", url, body, result, reqOpts...)
}

// Post sends the POST request
func (client *JSONClient) Post(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("POST", url, body, result, reqOpts...)
}

// Patch sends the PATCH request
func (client *JSONClient) Patch(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("PATCH", url, body, result, reqOpts...)
}

// Put sends the PUT request
func (client *JSONClient) Put(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("PUT", url, body, result, reqOpts...)
}

// Delete sends the DELETE request
func (client *JSONClient) Delete(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("DELETE", url, body, result, reqOpts...)
}

// Do sends a custom METHOD request
func (client *JSONClient) Do(method, url string, body, result interface{}, reqOpts ...RequestOption) error {
	var (
		bodyData  []byte
		resultStr string
		err       error
	)

	if body != nil {
		if bodyData, err = json.Marshal(body); err != nil {
			logger.Error(client.ctx, "marshal request body", "error", err)
			return err
		}
	}

	reqOpts = append(reqOpts, SetTypeJSON())

	if resultStr, err = client.Client.Do(method, url, string(bodyData), reqOpts...); err != nil {
		return err
	}

	if result != nil && resultStr != "" {
		if err = json.Unmarshal([]byte(resultStr), result); err != nil {
			logger.Error(client.ctx, "unmarshal response body", "error", err)
			return err
		}
	}
	return nil
}
