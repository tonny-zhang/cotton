package cotton

import (
	"net/http"
)

type responseWriter interface {
	http.ResponseWriter
	GetStatusCode() int
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
