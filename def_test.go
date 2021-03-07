package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func doRequest(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, nil)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}

func TestMain(m *testing.M) {
	SetMode(ModeTest)
	m.Run()
}
