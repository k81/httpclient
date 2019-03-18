package httpclient

import (
	"context"
	"encoding/xml"
)

// XMLClient is an wrapper of *Client, which talks in XML
type XMLClient struct {
	*Client
}

// NewXML create a XML http client instance with specified options
func NewXML(ctx context.Context, opts ...ClientOption) *XMLClient {
	client := New(ctx, opts...)
	return &XMLClient{client}
}

// Options sends the OPTIONS request
func (client *XMLClient) Options(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("OPTIONS", url, body, result, reqOpts...)
}

// Head sends the HEAD request
func (client *XMLClient) Head(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("HEAD", url, body, result, reqOpts...)
}

// Get sends the GET request
func (client *XMLClient) Get(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("GET", url, body, result, reqOpts...)
}

// Post sends the POST request
func (client *XMLClient) Post(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("POST", url, body, result, reqOpts...)
}

// Patch sends the PATCH request
func (client *XMLClient) Patch(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("PATCH", url, body, result, reqOpts...)
}

// Put sends the PUT request
func (client *XMLClient) Put(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("PUT", url, body, result, reqOpts...)
}

// Delete sends the DELETE request
func (client *XMLClient) Delete(url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do("DELETE", url, body, result, reqOpts...)
}

// Do sends a custom METHOD request
func (client *XMLClient) Do(method, url string, body, result interface{}, reqOpts ...RequestOption) error {
	var (
		bodyData  []byte
		resultStr string
		err       error
	)

	if body != nil {
		if bodyData, err = xml.Marshal(body); err != nil {
			logger.Error(client.ctx, "marshal request body", "error", err)
			return err
		}
	}

	reqOpts = append(reqOpts, SetTypeXML())

	if resultStr, err = client.Client.Do(method, url, string(bodyData), reqOpts...); err != nil {
		return err
	}

	if result != nil && resultStr != "" {
		if err = xml.Unmarshal([]byte(resultStr), result); err != nil {
			logger.Error(client.ctx, "unmarshal response body", "error", err)
			return err
		}
	}
	return nil
}
