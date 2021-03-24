package cotton

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getContext() *Context {
	url_req := "/get?a=1&b[]=11&b[]=12&c[a]=ca&c[b]=cb&list=1&list=2"
	return &Context{
		Request: httptest.NewRequest(http.MethodGet, url_req, nil),
	}
}
func TestGetQueryArray(t *testing.T) {
	c := getContext()

	arr := c.GetQueryArray("b[]")
	assert.Equal(t, 2, len(arr))
	assert.Equal(t, "11", arr[0])
	assert.Equal(t, "12", arr[1])

	arr = c.GetQueryArray("list")
	assert.Equal(t, 2, len(arr))
	assert.Equal(t, "1", arr[0])
	assert.Equal(t, "2", arr[1])

	arr = c.GetQueryArray("nokey")
	assert.Equal(t, 0, len(arr))
}
func TestGetQueryMap(t *testing.T) {
	c := getContext()

	m, ok := c.GetQueryMap("c")
	assert.True(t, ok)
	assert.Equal(t, 2, len(m))
	assert.Equal(t, "ca", m["a"])
	assert.Equal(t, "cb", m["b"])

	_, ok = c.GetQueryMap("nokey")
	assert.False(t, ok)
}
