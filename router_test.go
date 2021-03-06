package httpserver

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddHandleFunc(t *testing.T) {
	assert.PanicsWithError(t, "[hello] shold absolute path", func() {
		router := NewRouter()
		router.Get("hello", nil)
	})
	assert.PanicsWithError(t, "[/:name] conflicts with [/hello]", func() {
		router := NewRouter()
		router.Get("/hello", nil)
		router.Get("/:name", nil)
	})

	assert.NotPanics(t, func() {
		router := NewRouter()
		router.Get("/a", nil)
		router.Get("/b", nil)
	})
	assert.NotPanics(t, func() {
		router := NewRouter()
		router.Get("/a", nil)
		router.Get("/a/", nil)
	})

	assert.NotPanics(t, func() {
		router := NewRouter()
		router.Get("/a/b/", nil)
		router.Get("/b/:name/", nil)
	})

	assert.NotPanics(t, func() {
		router := NewRouter()
		router.Get("/a/", nil)
		router.Get("/a/:name", nil)
	})
}

func TestMatch(t *testing.T) {
	router := NewRouter()
	router.Get("/user", nil)
	router.Get("/user/", nil)
	router.Get("/user/:name", nil)

	resultMatch := router.match(http.MethodGet, "/user/test")

	assert.Equal(t, "/user/:name", resultMatch.rule.rule)

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
