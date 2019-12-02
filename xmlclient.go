package httpclient

import (
	"context"
	"encoding/xml"

	"go.uber.org/zap"
)

// XMLClient is an wrapper of *Client, which talks in XML
type XMLClient struct {
	*Client
}

// NewXML create a XML http client instance with specified options
func NewXML(logger *zap.Logger, opts ...ClientOption) *XMLClient {
	client := New(logger, opts...)
	return &XMLClient{client}
}

// Options sends the OPTIONS request
func (client *XMLClient) Options(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "OPTIONS", url, body, result, reqOpts...)
}

// Head sends the HEAD request
func (client *XMLClient) Head(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "HEAD", url, body, result, reqOpts...)
}

// Get sends the GET request
func (client *XMLClient) Get(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "GET", url, body, result, reqOpts...)
}

// Post sends the POST request
func (client *XMLClient) Post(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "POST", url, body, result, reqOpts...)
}

// Patch sends the PATCH request
func (client *XMLClient) Patch(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "PATCH", url, body, result, reqOpts...)
}

// Put sends the PUT request
func (client *XMLClient) Put(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "PUT", url, body, result, reqOpts...)
}

// Delete sends the DELETE request
func (client *XMLClient) Delete(ctx context.Context, url string, body, result interface{}, reqOpts ...RequestOption) error {
	return client.Do(ctx, "DELETE", url, body, result, reqOpts...)
}

// Do sends a custom METHOD request
func (client *XMLClient) Do(ctx context.Context, method, url string, body, result interface{}, reqOpts ...RequestOption) error {
	var (
		bodyData  []byte
		resultStr string
		err       error
	)

	if body != nil {
		switch bodyValue := body.(type) {
		case string:
			bodyData = []byte(bodyValue)
		case []byte:
			bodyData = bodyValue
		default:
			if bodyData, err = xml.Marshal(body); err != nil {
				client.logger.Error("marshal request body", zap.Error(err))
				return err
			}
		}
	}

	reqOpts = append([]RequestOption{SetTypeXML()}, reqOpts...)

	if resultStr, err = client.Client.Do(ctx, method, url, string(bodyData), reqOpts...); err != nil {
		return err
	}

	if result != nil && resultStr != "" {
		if err = xml.Unmarshal([]byte(resultStr), result); err != nil {
			client.logger.Error("unmarshal response body", zap.Error(err))
			return err
		}
	}
	return nil
}
