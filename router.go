package cotton

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tonny-zhang/cotton/utils"
)

// HandlerFunc handler func
type HandlerFunc func(ctx *Context)

// Router router struct
type Router struct {
	prefix      string
	trees       map[string]*tree
	middlewares []HandlerFunc
}

// NewRouter new router
func NewRouter() Router {
	return Router{
		trees:       make(map[string]*tree),
		middlewares: make([]HandlerFunc, 0),
	}
}

// Group get group router
func (router *Router) Group(path string, handler ...HandlerFunc) Router {
	r := Router{
		prefix:      path,
		trees:       router.trees,
		middlewares: router.middlewares,
	}
	r.middlewares = append(r.middlewares, handler...)
	return r
}
func (router *Router) addHandleFunc(method, path string, handler HandlerFunc) {
	if "" != router.prefix {
		path = utils.CleanPath(router.prefix + "/" + path)
	}
	if !strings.HasPrefix(path, "/") {
		panic(fmt.Errorf("[%s] shold absolute path", path))
	}
	if handler == nil {
		panic(fmt.Errorf("%s %s has no handler", method, path))
	}
	if _, ok := router.trees[method]; !ok {
		router.trees[method] = newTree()
	}
	nodeAdded := router.trees[method].Add(path, nil)
	nodeAdded.middleware = append(nodeAdded.middleware, router.middlewares...)
	nodeAdded.handler = handler
	nodeAdded.middleware = append(nodeAdded.middleware, handler)
	debugPrintRoute(method, path, handler)
}

// Get router get method
func (router *Router) Get(path string, handler HandlerFunc) {
	router.addHandleFunc(http.MethodGet, path, handler)
}

// Post router post method
func (router *Router) Post(path string, handler HandlerFunc) {
	router.addHandleFunc(http.MethodPost, path, handler)
}

// Put router put method
func (router *Router) Put(path string, handler HandlerFunc) {
	router.addHandleFunc(http.MethodPut, path, handler)
}

// Options router options method
func (router *Router) Options(path string, handler HandlerFunc) {
	router.addHandleFunc(http.MethodOptions, path, handler)
}

// Delete router delete method
func (router *Router) Delete(path string, handler HandlerFunc) {
	router.addHandleFunc(http.MethodDelete, path, handler)
}

// Patch router patch method
func (router *Router) Patch(path string, handler HandlerFunc) {
	router.addHandleFunc(http.MethodPatch, path, handler)
}

// Head router head method
func (router *Router) Head(path string, handler HandlerFunc) {
	router.addHandleFunc(http.MethodHead, path, handler)
}
func (router *Router) PrintTree(method string) {
	router.trees[method].root.print(0)
}
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := Context{
		Request:  r,
		Response: w,
		index:    -1,
	}
	// ctx.index++
	r.Method = strings.ToUpper(r.Method)
	if tree, ok := router.trees[r.Method]; ok {
		result := tree.Find(r.URL.Path)
		if result != nil {
			ctx.paramCache = result.params
			ctx.handlers = result.node.middleware

			ctx.Next()
			return
		}
	}

	ctx.handlers = router.middlewares
	ctx.NotFound()
}

// Run run for http
func (router *Router) Run(addr string) {
	debugPrint("Listening and serving HTTP on %s\n", addr)
	http.ListenAndServe(addr, router)
}

// Use use for middleware
func (router *Router) Use(handler ...HandlerFunc) {
	router.middlewares = append(router.middlewares, handler...)
}
