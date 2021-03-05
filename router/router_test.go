package router

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddHandleFunc(t *testing.T) {
	assert.PanicsWithError(t, "[hello] shold absolute path", func() {
		router := New()
		router.Get("hello", nil)
	})
	assert.PanicsWithError(t, "[/:name] conflicts with [/hello]", func() {
		router := New()
		router.Get("/hello", nil)
		router.Get("/:name", nil)
	})

	assert.NotPanics(t, func() {
		router := New()
		router.Get("/a", nil)
		router.Get("/b", nil)
	})
	assert.NotPanics(t, func() {
		router := New()
		router.Get("/a", nil)
		router.Get("/a/", nil)
	})

	assert.NotPanics(t, func() {
		router := New()
		router.Get("/a/b/", nil)
		router.Get("/b/:name/", nil)
	})

	assert.NotPanics(t, func() {
		router := New()
		router.Get("/a/", nil)
		router.Get("/a/:name", nil)
	})
}

func TestMatch(t *testing.T) {
	router := New()
	router.Get("/user", nil)
	router.Get("/user/", nil)
	router.Get("/user/:name", nil)

	resultMatch := router.match(http.MethodGet, "/user/test")

	assert.Equal(t, "/user/:name", resultMatch.rule.rule)

}
