package cotton

import (
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
	assert.PanicsWithError(t, "path [/:name] conflicts with [/hello]", func() {
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

func TestCustomNotFound(t *testing.T) {
	router := NewRouter()

	infoCustomNotFound := "not found from custom"
	router.NotFound(func(ctx *Context) {
		ctx.String(http.StatusNotFound, infoCustomNotFound)
	})

	w := doRequest(router, http.MethodGet, "/path404")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, infoCustomNotFound, w.Body.String())
}

func TestCustomGroupNotFound(t *testing.T) {
	router := NewRouter()

	infoCustomNotFound := "not found from custom"
	infoCustomGroupNotFound := "not found from custom group"
	router.NotFound(func(ctx *Context) {
		ctx.String(http.StatusNotFound, infoCustomNotFound)
	})

	g := router.Group("/v1")
	g.NotFound(func(ctx *Context) {
		ctx.String(http.StatusNotFound, infoCustomGroupNotFound)
	})

	w := doRequest(router, http.MethodGet, "/path404")
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, infoCustomNotFound, w.Body.String())

	w = doRequest(router, http.MethodGet, "/v1/path404")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, infoCustomGroupNotFound, w.Body.String())
}
func TestGroup(t *testing.T) {
	router := NewRouter()
	g1 := router.Group("/v1")
	g1.Get("/a", func(c *Context) {
		c.String(http.StatusOK, "g1 a")
	})
	g1.Get("/b", func(c *Context) {
		c.String(http.StatusBadGateway, "g1 b")
	})

	w := doRequest(router, http.MethodGet, "/v1/a")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "g1 a", w.Body.String())

	w = doRequest(router, http.MethodGet, "/v1/b")

	assert.Equal(t, http.StatusBadGateway, w.Code)
	assert.Equal(t, "g1 b", w.Body.String())
	// assert.True(t, false)
}
func TestGroupPanic(t *testing.T) {
	router := NewRouter()
	assert.PanicsWithError(t, "group [] must start with /", func() {
		router.Group("")
	})
	assert.PanicsWithError(t, "group [abc] must start with /", func() {
		router.Group("abc")
	})

	assert.PanicsWithError(t, "group path [/:test] can not has parameter", func() {
		router.Group("/:test")
	})

	assert.PanicsWithError(t, "group [/a/] can not group again", func() {
		router.Group("/a").Group("/a")
	})
}
