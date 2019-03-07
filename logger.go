package httpclient

import "github.com/k81/log"

// logger is the default logger used
var logger *log.Logger = log.GetLogger()

func SetLogger(l *log.Logger) {
	logger = l
}
