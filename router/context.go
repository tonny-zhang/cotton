package router

import "net/http"

// Context context for request
type Context struct {
	Request         *http.Request
	writer          http.ResponseWriter
	ruleMatchResult matchResult
}

// GetQuery for Request.URL.Query().Get
func (ctx *Context) GetQuery(key string) string {
	return ctx.Request.URL.Query().Get(key)
}

// GetQuery1 for Request.URL.Query().Get
func (ctx *Context) GetQuery1(key string) string {
	val, ok := ctx.ruleMatchResult.Params[key]
	if ok {
		return val
	}
	return ""
}

func (ctx *Context) String(content string) {
	ctx.writer.Write([]byte(content + "\n"))
}
