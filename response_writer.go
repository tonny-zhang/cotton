package cotton

import (
	"net/http"
)

type responseWriter interface {
	http.ResponseWriter
	GetStatusCode() int
	GetHTTPResponseWriter() http.ResponseWriter
}

type resWriter struct {
	http.ResponseWriter
	statusCode int
}

var _ responseWriter = &resWriter{}

func (w *resWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
func (w *resWriter) GetStatusCode() int {
	return w.statusCode
}
func (w *resWriter) GetHTTPResponseWriter() http.ResponseWriter {
	return w.ResponseWriter
}
