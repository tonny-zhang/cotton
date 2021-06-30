package cotton

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getContext() *Context {
	urlReq := "/get?a=1&b[]=11&b[]=12&c[a]=ca&c[b]=cb&list=1&list=2"
	return &Context{
		Request: httptest.NewRequest(http.MethodGet, urlReq, nil),
	}
}
func postContext() *Context {
	urlReq := "/post"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "abc")
	writer.WriteField("ids", "1")
	writer.WriteField("ids", "2")

	writer.Close()
	req := httptest.NewRequest(http.MethodPost, urlReq, body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=--------------------------611593185451210078804896")

	return &Context{
		Request: req,
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

func TestGetPostForm(t *testing.T) {
	body := bytes.NewBufferString("foo=bar&ids=1&ids=2&info[name]=test&info[age]=10")
	req := httptest.NewRequest(http.MethodPost, "/post?get=1", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := &Context{
		Request: req,
	}

	assert.Equal(t, "1", c.GetQuery("get"))
	assert.Equal(t, "", c.GetPostForm("get"))

	assert.Equal(t, "bar", c.GetPostForm("foo"))
	assert.Equal(t, "", c.GetPostForm("nokey"))

	assert.Equal(t, []string{"1", "2"}, c.GetPostFormArray("ids"))

	m, ok := c.GetPostFormMap("info")
	assert.True(t, ok)
	assert.Equal(t, map[string]string{
		"name": "test",
		"age":  "10",
	}, m)

	_, ok = c.GetPostFormMap("nokey")
	assert.False(t, ok)
}

func TestGetMultipartPostForm(t *testing.T) {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	mw.SetBoundary(boundary)
	mw.WriteField("foo", "bar")
	mw.WriteField("ids", "1")
	mw.WriteField("ids", "2")
	mw.WriteField("info[name]", "test")
	mw.WriteField("info[age]", "10")
	req, _ := http.NewRequest("POST", "/?get=1", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	mw.Close()

	c := &Context{
		Request: req,
	}

	assert.Equal(t, "1", c.GetQuery("get"))
	assert.Equal(t, "", c.GetPostForm("get"))

	assert.Equal(t, "bar", c.GetPostForm("foo"))
	assert.Equal(t, "", c.GetPostForm("nokey"))

	assert.Equal(t, []string{"1", "2"}, c.GetPostFormArray("ids"))

	m, ok := c.GetPostFormMap("info")
	assert.True(t, ok)
	assert.Equal(t, map[string]string{
		"name": "test",
		"age":  "10",
	}, m)

	_, ok = c.GetPostFormMap("nokey")
	assert.False(t, ok)
}

func TestSyncPool(t *testing.T) {
	c1 := newContext(nil, nil, nil)
	c2 := newContext(nil, nil, nil)

	c1.index = 1
	c2.index = 2

	assert.NotEqual(t, c1.index, c2.index)
	// assert.False(t, true)
}
