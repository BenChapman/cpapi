package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func NewLogger(logger *logrus.Logger) negroni.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		// Start timer
		start := time.Now()
		path := request.URL.Path

		// Process request
		writer := negroni.NewResponseWriter(w)
		next(writer, request)

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := request.RemoteAddr
		method := request.Method
		statusCode := writer.Status()
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)

		msg := fmt.Sprintf("%v |%s %3d %s| %13v | %s |%s  %s %-7s %s",
			end.Format("2006/01/02 - 15:04:05"),
			statusColor, statusCode, reset,
			latency,
			clientIP,
			methodColor, reset, method,
			path,
		)

		logger.Info(msg)
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}
