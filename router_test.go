package cotton

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var handlerEmpty = func(c *Context) {}

func TestDiffentParamNameInSamePos(t *testing.T) {
	// same postion param name must be same
	assert.PanicsWithError(t, "[:name] in path [/v1/abc/abc/:name/:id] conflicts with [/v1/abc/abc/:id]", func() {
		router := NewRouter()
		g := router.Group("/v1/abc/")
		g.Post("/abc/:id", handlerEmpty)
		g.Post("/abc/:name/:id", handlerEmpty)
	})
}
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
