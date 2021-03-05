package httpserver

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

// Context context for request
type Context struct {
	Request         *http.Request
	writer          http.ResponseWriter
	ruleMatchResult matchResult
	statusCode      int
	handlers        []HandlerFunc

	queryCache url.Values
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
		for _, handler := range ctx.handlers {
			handler(ctx)
		}
	}
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
	val, ok := ctx.ruleMatchResult.Params[key]
	if ok {
		return val
	}
	return ""
}

// StatusCode set status code
func (ctx *Context) StatusCode(statusCode int) {
	ctx.statusCode = statusCode
}
func (ctx *Context) String(content string) {
	ctx.StatusCode(200)
	ctx.writer.Write([]byte(content + "\n"))
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