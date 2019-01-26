package httpclient

import "fmt"

// HTTPError is the http error status code info, which is not in range [200,300)
type HTTPError struct {
	StatusCode int
	StatusText string
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP Error: %v, %v", e.StatusCode, e.StatusText)
}
