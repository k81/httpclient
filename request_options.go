package httpclient

import (
	"context"
	"net/http"
	"net/url"
)

// RequestOption defines the request option to customize the request
type RequestOption func(ctx context.Context, req *http.Request) (newctx context.Context, err error)

// SetHeader sets the request header
func SetHeader(key, value string) RequestOption {
	return func(ctx context.Context, req *http.Request) (context.Context, error) {
		req.Header.Set(key, value)
		return ctx, nil
	}
}

// SetTypeXML sets the Content-Type to `application/xml`
func SetTypeXML() RequestOption {
	return SetHeader("Content-Type", "application/xml; charset=UTF-8")
}

// SetTypeJSON sets the Content-Type to `application/json`
func SetTypeJSON() RequestOption {
	return SetHeader("Content-Type", "application/json; charset=UTF-8")
}

// SetTypeForm sets the Content-Type to `application/x-www-form-urlencoded`
func SetTypeForm() RequestOption {
	return SetHeader("Content-Type", "application/x-www-form-urlencoded")
}

// SetQuery sets the query params
func SetQuery(values url.Values) RequestOption {
	return func(ctx context.Context, req *http.Request) (context.Context, error) {
		q := req.URL.Query()
		for k, v := range values {
			for _, vv := range v {
				q.Add(k, vv)
			}
		}
		req.URL.RawQuery = q.Encode()
		return ctx, nil
	}
}
