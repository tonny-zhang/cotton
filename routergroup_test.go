package httpserver

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	router := NewRouter()
	g1 := router.Group("/v1")
	passed := false
	g1.Get("/a", func(c *Context) {
		passed = true
	})

	w := doRequest(&router, http.MethodGet, "/v1/a")

	fmt.Println(w.Code, w.Body.String())
	fmt.Println(passed)
	assert.True(t, passed)
}
