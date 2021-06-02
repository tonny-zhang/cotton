package cotton

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainAfterGroup(t *testing.T) {
	assert.PanicsWithError(t, "group can not call Domain", func() {
		router := NewRouter()
		router.Group("/hello", nil).Domain("www", nil)
	})
}
func TestDomainAfterDomain(t *testing.T) {
	assert.PanicsWithError(t, "Domain can not call Domain", func() {
		router := NewRouter()
		router.Domain("/hello", nil).Domain("www", nil)
	})
}
func TestDomainExists(t *testing.T) {
	assert.PanicsWithError(t, "domain [www] is exists", func() {
		router := NewRouter()
		router.Domain("www", nil)
		router.Domain("www", nil)
	})
}

func TestDomain(t *testing.T) {
	router := NewRouter()

	d1 := router.Domain("a.test.com")
	d1.Get("/test", func(ctx *Context) {
		ctx.String(http.StatusOK, "d1 test")
	})

	d2 := router.Domain("b.test.com")
	d2.Get("/test", func(ctx *Context) {
		ctx.String(http.StatusOK, "d2 test")
	})

	w := doRequest(router, http.MethodGet, "http://a.test.com/test")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "d1 test", w.Body.String())

	w = doRequest(router, http.MethodGet, "http://b.test.com/test")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "d2 test", w.Body.String())
}
