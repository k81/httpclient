package httpclient

import (
	"net/http"
	"net/url"
)

// RequestOption defines the request option to customize the request
type RequestOption func(*http.Request) error

// SetHeader sets the request header
func SetHeader(key, value string) RequestOption {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
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
	return func(req *http.Request) error {
		q := req.URL.Query()
		for k, v := range values {
			for _, vv := range v {
				q.Add(k, vv)
			}
		}
		req.URL.RawQuery = q.Encode()
		return nil
	}
}
