package cotton

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddHandleFunc(t *testing.T) {
	handler := func(c *Context) {}
	assert.PanicsWithError(t, "[hello] shold absolute path", func() {
		router := NewRouter()
		router.Get("hello", nil)
	})
	assert.PanicsWithError(t, "[/:name] conflicts with [/hello]", func() {
		router := NewRouter()
		router.Get("/hello", handler)
		router.Get("/:name", handler)
	})

	assert.NotPanics(t, func() {
		router := NewRouter()
		router.Get("/a", handler)
		router.Get("/b", handler)
	})
	assert.NotPanics(t, func() {
		router := NewRouter()
		router.Get("/a", handler)
		router.Get("/a/", handler)
	})

	assert.NotPanics(t, func() {
		router := NewRouter()
		router.Get("/a/b/", handler)
		router.Get("/b/:name/", handler)
	})

	assert.NotPanics(t, func() {
		router := NewRouter()
		router.Get("/a/", handler)
		router.Get("/a/:name", handler)
	})
}

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
