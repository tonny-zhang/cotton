package cotton

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
func TestGroupPrefix(t *testing.T) {
	router := NewRouter()
	g1 := router.Group("/v1")
	g1.Get("/a", func(c *Context) {
		c.String(http.StatusOK, "g1 a")
	})

	w := doRequest(router, http.MethodGet, "/v1/v1/a")

	assert.Equal(t, http.StatusNotFound, w.Code)
}
func TestGroupPanic(t *testing.T) {
	assert.PanicsWithError(t, "group [] must start with /", func() {
		router := NewRouter()
		router.Group("")
	})
	assert.PanicsWithError(t, "group [abc] must start with /", func() {
		router := NewRouter()
		router.Group("abc")
	})

	assert.PanicsWithError(t, "group [/a/] is setted", func() {
		router := NewRouter()
		router.Group("/a")
		router.Group("/a")
	})

	assert.PanicsWithError(t, "group [/a/b/] is setted", func() {
		router := NewRouter()
		router.Group("/a").Group("/b")
		router.Group("/a//b")
	})

	assert.PanicsWithError(t, "group path [/:method] can not has parameter", func() {
		router := NewRouter()
		router.Group("/:method")
	})

	assert.NotPanics(t, func() {
		router := NewRouter()
		router.Group("/s")
		router.Group("/static")
	})
}

func TestMatchGroup(t *testing.T) {
	assert.True(t, matchGroup(&Router{
		prefix: "/v1/",
	}, "/v1/test"))

	assert.False(t, matchGroup(&Router{
		prefix: "/v1/",
	}, "/v2/test"))
}

func TestCustomGroupNotFound(t *testing.T) {
	router := NewRouter()

	infoCustomNotFound := "not found from custom"
	infoCustomGroupNotFound := "not found from custom group"
	infoCustomGroupUserNotFound := "not found from custom group user"
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

	gUser := g.Group("/user")

	w = doRequest(router, http.MethodGet, "/v1/user/path404")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, infoCustomGroupNotFound, w.Body.String())

	gUser.NotFound(func(ctx *Context) {
		ctx.String(http.StatusNotFound, infoCustomGroupUserNotFound)
	})

	w = doRequest(router, http.MethodGet, "/v1/user/path404")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, infoCustomGroupUserNotFound, w.Body.String())

	w = doRequest(router, http.MethodGet, "/v1/user1/path404")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, infoCustomGroupNotFound, w.Body.String())

	w = doRequest(router, http.MethodGet, "/v2/user1/path404")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, infoCustomNotFound, w.Body.String())
}

func TestGroupMulty(t *testing.T) {
	router := NewRouter()
	g1 := router.Group("/a")
	g1.addHandleFunc("GET", "/test", func(ctx *Context) {
		ctx.String(http.StatusOK, "/a/test")
	})
	g2 := g1.Group("/b")
	g2.addHandleFunc("GET", "/test", func(ctx *Context) {
		ctx.String(http.StatusOK, "/a/b/test")
	})

	w := doRequest(router, http.MethodGet, "/a/test")
	assert.Equal(t, "/a/test", w.Body.String())

	w = doRequest(router, http.MethodGet, "/a/b/test")
	assert.Equal(t, "/a/b/test", w.Body.String())
}

func TestCustomGroupNotFoundOrder(t *testing.T) {

	infoCustomNotFound := "not found from custom"
	infoCustomGroupNotFound := "not found from custom group"
	infoCustomGroupUserNotFound := "not found from custom group user"

	{
		router := NewRouter()
		router.NotFound(func(ctx *Context) {
			ctx.String(http.StatusNotFound, infoCustomNotFound)
		})
		g := router.Group("/v1")
		g.NotFound(func(ctx *Context) {
			ctx.String(http.StatusNotFound, infoCustomGroupNotFound)
		})
		gUser := router.Group("/v1/user")
		gUser.NotFound(func(ctx *Context) {
			ctx.String(http.StatusNotFound, infoCustomGroupUserNotFound)
		})

		router.sort()
		w := doRequest(router, http.MethodGet, "/v1/user/path404")
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, infoCustomGroupUserNotFound, w.Body.String())
	}

	{
		router := NewRouter()
		router.NotFound(func(ctx *Context) {
			ctx.String(http.StatusNotFound, infoCustomNotFound)
		})
		gUser := router.Group("/v1/user")
		gUser.NotFound(func(ctx *Context) {
			ctx.String(http.StatusNotFound, infoCustomGroupUserNotFound)
		})
		g := router.Group("/v1")
		g.NotFound(func(ctx *Context) {
			ctx.String(http.StatusNotFound, infoCustomGroupNotFound)
		})

		router.sort()
		w := doRequest(router, http.MethodGet, "/v1/user/path404")
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, infoCustomGroupUserNotFound, w.Body.String())
	}
}
