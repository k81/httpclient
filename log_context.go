package httpclient

import (
	"context"
	"net/http"
)

// LogContextFunc return a new context for the request
type LogContextFunc func(context.Context, *http.Request) context.Context
