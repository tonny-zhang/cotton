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

	paramCache map[string]string
	queryCache url.Values

	router *Router
}

func newContext(w http.ResponseWriter, r *http.Request, router *Router) *Context {
	// use sync.Pool
	ctx := ctxPool.Get().(*Context)

	// reset all property
	ctx.Request = r
	ctx.Response = &resWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
	ctx.router = router
	ctx.indexAbort = -1
	ctx.index = -1
	ctx.handlers = ctx.handlers[0:0]
	ctx.paramCache = nil
	ctx.queryCache = nil

	ctxPool.Put(ctx)
	return ctx
}
func (ctx *Context) initQueryCache() {
	if nil == ctx.queryCache {
		if nil == ctx.Request {
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
	ctx.Next()
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

// // StatusCode set status code
// func (ctx *Context) StatusCode(statusCode int) {
// 	if ctx.Response.statusCode == 0 {
// 		ctx.statusCode = statusCode
// 		ctx.Response.WriteHeader(statusCode)
// 	} else {
// 		fmt.Printf("warning: alread set statusCode [%d], can't set [%d] again\n", ctx.statusCode, statusCode)
// 	}
// }

// response with string
func (ctx *Context) String(code int, content string) {
	ctx.Response.WriteHeader(code)
	ctx.Response.Write([]byte(content))
}

// response with json
func (ctx *Context) JSON(code int, val M) {
	b, e := json.Marshal(val)
	if e != nil {
		panic(e)
	}

	ctx.Response.Header().Add("Content-Type", jsonContentType)
	ctx.Response.WriteHeader(code)
	ctx.Response.Write(b)
}

// response with html
func (ctx *Context) HTML(code int, html string) {
	ctx.Response.Header().Add("Content-Type", htmlContentType)
	ctx.Response.WriteHeader(code)
	ctx.Response.Write([]byte(html))
}
func (ctx *Context) getRequestHeader(key string) string {
	if nil != ctx.Request {
		return ctx.Request.Header.Get(key)
	}
	return ""
}

// ClientIP get client ip
func (ctx *Context) ClientIP() string {
	clientIP := ctx.getRequestHeader("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(ctx.getRequestHeader("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}

	if addr := ctx.getRequestHeader("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(ctx.Request.RemoteAddr)); err == nil {
		return ip
	}

	return "no-ip"
}
