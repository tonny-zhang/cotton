package cotton

import (
	"net/http"
	"sync"
	"testing"
	"time"

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
	assert.PanicsWithError(t, "group [] must start with /", func() {
		router := NewRouter()
		router.Group("")
	})
	assert.PanicsWithError(t, "group [abc] must start with /", func() {
		router := NewRouter()
		router.Group("abc")
	})

	assert.PanicsWithError(t, "group [/a/] can not group again", func() {
		router := NewRouter()
		router.Group("/a").Group("/a")
	})

	assert.PanicsWithError(t, "group [/a/] conflicts with [/a/]", func() {
		router := NewRouter()
		router.Group("/a")
		router.Group("/a")
	})
	assert.PanicsWithError(t, "group [/b/] conflicts with [/:method/]", func() {
		router := NewRouter()
		router.Group("/:method")
		router.Group("/b")
	})
	assert.PanicsWithError(t, "group [/:method/] conflicts with [/a/]", func() {
		router := NewRouter()
		router.Group("/a")
		router.Group("/:method")
	})
	assert.PanicsWithError(t, "group [/:id/] conflicts with [/:method/]", func() {
		router := NewRouter()
		router.Group("/:method")
		router.Group("/:id")
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

	assert.True(t, matchGroup(&Router{
		prefix: "/v1/:method/",
	}, "/v1/test/"))
	assert.True(t, matchGroup(&Router{
		prefix: "/v1/:method/",
	}, "/v1/test/abc"))
	assert.False(t, matchGroup(&Router{
		prefix: "/v1/:method/",
	}, "/v1/test"))
	assert.False(t, matchGroup(&Router{
		prefix: "/v1/:method/",
	}, "/v2/test/"))
}

func TestMultipleEOP(t *testing.T) {
	router := NewRouter()
	content := "router a"
	router.Get("/a", func(c *Context) {
		c.String(http.StatusOK, content)
	})

	w := doRequest(router, http.MethodGet, "////a")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, content, w.Body.String())
}

func TestMergeRequest(t *testing.T) {
	var wg sync.WaitGroup
	router := NewRouter()
	contentA := "router a"
	contentB := "router b"
	router.Get("/a", func(c *Context) {
		time.Sleep(time.Second * 3)
		c.String(http.StatusOK, contentA)
	})
	router.Get("/b", func(c *Context) {
		time.Sleep(time.Second * 1)
		c.String(http.StatusOK, contentB)
	})

	wg.Add(2)
	var resultA, resultB string
	go func() {
		defer wg.Done()

		w := doRequest(router, http.MethodGet, "/a")
		resultA = w.Body.String()
	}()
	time.Sleep(time.Second)

	go func() {
		defer wg.Done()
		w := doRequest(router, http.MethodGet, "/b")
		resultB = w.Body.String()
	}()
	wg.Wait()

	assert.Equal(t, contentA, resultA)
	assert.Equal(t, contentB, resultB)
}
