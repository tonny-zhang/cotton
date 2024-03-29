package cotton

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	htmlContentType = "text/html; charset=utf-8"
	jsonContentType = "application/json; charset=utf-8"
)

const (
	defaultMultipartMemory = 32 << 20 // 32 MB
)

var ctxPool sync.Pool

func init() {
	ctxPool.New = func() interface{} {
		return &Context{}
	}
}

// Context context for request
type Context struct {
	Request  *http.Request
	Response responseWriter
	// statusCode int
	handlers   []HandlerFunc
	index      int8
	indexAbort int8

	paramCache    map[string]string
	queryCache    url.Values
	postFormCache url.Values

	values map[string]interface{}

	router *Router
}

func newContext(w http.ResponseWriter, r *http.Request, router *Router) *Context {
	// use sync.Pool
	ctx := ctxPool.Get().(*Context)

	ctx.reset()

	ctx.Request = r
	ctx.Response = &resWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	ctx.router = router

	return ctx
}
func (ctx *Context) reset() {
	ctx.Request = nil
	ctx.Response = nil
	ctx.handlers = ctx.handlers[0:0]
	ctx.index = -1
	ctx.indexAbort = -1

	ctx.paramCache = nil
	ctx.queryCache = nil
	ctx.postFormCache = nil

	ctx.values = nil

	ctx.router = nil
}
func (ctx *Context) destroy() {
	ctx.reset()
	ctxPool.Put(ctx)
}
func (ctx *Context) initQueryCache() {
	if nil == ctx.queryCache {
		if nil != ctx.Request {
			ctx.queryCache = ctx.Request.URL.Query()
		} else {
			ctx.queryCache = url.Values{}
		}
	}
}

// Next fn
func (ctx *Context) Next() {
	if nil != ctx.handlers {
		ctx.index++
		num := int8(len(ctx.handlers))
		for ctx.index < num {
			// NOTICE: there ctx will escape
			ctx.handlers[ctx.index](ctx)
			ctx.index++
		}
	}
}

// Abort fn
func (ctx *Context) Abort() {
	ctx.index = int8(len(ctx.handlers) + 1)
}

// NotFound for 404
func (ctx *Context) NotFound() {
	// ctx.StatusCode(http.StatusNotFound)
	ctx.Response.WriteHeader(http.StatusNotFound)
	// http.NotFound(ctx.Response, ctx.Request)
	ctx.Response.Write([]byte("404 page not found"))
	ctx.Abort()
}

// Set set value
func (ctx *Context) Set(key string, val interface{}) {
	if ctx.values == nil {
		ctx.values = make(map[string]interface{})
	}
	ctx.values[key] = val
}

// Get get value
func (ctx *Context) Get(key string) (interface{}, bool) {
	v, ok := ctx.values[key]
	return v, ok
}

// Cookie get cookie
func (ctx *Context) Cookie(key string) (string, error) {
	c, e := ctx.Request.Cookie(key)

	if e == nil {
		return c.Value, nil
	}
	return "", e
}

// GetQuery for Request.URL.Query().Get
func (ctx *Context) GetQuery(key string) string {
	return ctx.GetDefaultQuery(key, "")
}

// GetDefaultQuery get default query
func (ctx *Context) GetDefaultQuery(key, defaultVal string) string {
	ctx.initQueryCache()
	if v, ok := ctx.queryCache[key]; ok {
		return v[0]
	}
	return defaultVal
}

// GetQueryArray get query array
// url?list[]=1&list[]=2
// 		GetQueryArray("list[]")	=> ["1", "2"]
func (ctx *Context) GetQueryArray(key string) (list []string) {
	ctx.initQueryCache()
	if v, ok := ctx.queryCache[key]; ok {
		return v
	}
	return
}
func getValue(m map[string][]string, key string) (dicts map[string]string, exists bool) {
	dicts = make(map[string]string)
	for k, v := range m {
		if i := strings.IndexByte(k, '['); i > 0 && k[:i] == key {
			if j := strings.IndexByte(k, ']'); j > 2 {
				dicts[k[i+1:j]] = v[0]
				exists = true
			}
		}
	}
	return
}

// GetQueryMap get query map
func (ctx *Context) GetQueryMap(key string) (dicts map[string]string, exists bool) {
	ctx.initQueryCache()
	return getValue(ctx.queryCache, key)
}

// GetAllQuery get all query value
func (ctx Context) GetAllQuery() url.Values {
	ctx.initQueryCache()
	return ctx.queryCache
}

// Param returns the value of the URL param.
//     router.GET("/user/:id", func(c *gin.Context) {
//         // a GET request to /user/john
//         id := c.Param("id") // id == "john"
//     })
func (ctx *Context) Param(key string) string {
	if ctx.paramCache != nil {
		val, ok := ctx.paramCache[key]
		if ok {
			return val
		}
	}

	return ""
}

// response with string
func (ctx *Context) String(code int, content string) {
	ctx.Response.WriteHeader(code)
	ctx.Response.Write([]byte(content))
}

// JSON response with json
func (ctx *Context) JSON(code int, val M) {
	b, e := json.Marshal(val)
	if e != nil {
		panic(e)
	}

	ctx.Response.Header().Add("Content-Type", jsonContentType)
	ctx.Response.WriteHeader(code)
	ctx.Response.Write(b)
}

// HTML response with html
func (ctx *Context) HTML(code int, html string) {
	ctx.Response.Header().Add("Content-Type", htmlContentType)
	ctx.Response.WriteHeader(code)
	ctx.Response.Write([]byte(html))
}

// GetRequestHeader get request header, short for ctx.Request.Header.Get
func (ctx *Context) GetRequestHeader(key string) string {
	if nil != ctx.Request {
		return ctx.Request.Header.Get(key)
	}
	return ""
}

// ClientIP get client ip
func (ctx *Context) ClientIP() string {
	clientIP := ctx.GetRequestHeader("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(ctx.GetRequestHeader("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}

	if addr := ctx.GetRequestHeader("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(ctx.Request.RemoteAddr)); err == nil {
		return ip
	}

	return "no-ip"
}

// Redirect short for http.Redirect
func (ctx *Context) Redirect(code int, location string) {
	http.Redirect(ctx.Response, ctx.Request, location, code)
}
