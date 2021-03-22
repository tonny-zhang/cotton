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

	notfoundHandlers []HandlerFunc

	groups []*Router
}

// NewRouter new router
func NewRouter() *Router {
	return &Router{
		trees:       make(map[string]*tree),
		middlewares: make([]HandlerFunc, 0),
	}
}

// Group get group router
func (router *Router) Group(path string, handler ...HandlerFunc) *Router {
	if router.prefix != "" {
		panic(fmt.Errorf("group [%s] can not group again", router.prefix))
	}
	if len(path) == 0 || path[0] != '/' {
		panic(fmt.Errorf("group [%s] must start with /", path))
	}
	if strings.Index(path, ":") > -1 || strings.Index(path, "*") > -1 {
		panic(fmt.Errorf("group path [%s] can not has parameter", path))
	}
	r := &Router{
		prefix:           utils.CleanPath(path + "/"),
		trees:            router.trees,
		middlewares:      router.middlewares,
		notfoundHandlers: router.notfoundHandlers,
	}
	r.middlewares = append(r.middlewares, handler...)
	router.groups = append(router.groups, r)
	return r
}

// NotFound custom NotFoundHandler
func (router *Router) NotFound(handler ...HandlerFunc) {
	router.notfoundHandlers = handler
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
	nodeAdded := router.trees[method].add(path, nil)
	nodeAdded.middleware = append(nodeAdded.middleware, router.middlewares...)
	// nodeAdded.handler = handler
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

// func (router *Router) PrintTree(method string) {
// 	router.trees[method].root.print(0)
// }

// func (router *Router) Find(method, path string) {
// 	if tree, ok := router.trees[method]; ok {
// 		tree.root.find(path)
// 	}
// }
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(w, r)

	r.Method = strings.ToUpper(r.Method)
	reqURI := r.URL.Path
	if tree, ok := router.trees[r.Method]; ok {
		result := tree.root.find(reqURI)

		if result.node != nil {
			ctx.paramCache = result.params
			ctx.handlers = result.node.middleware

			ctx.Next()
			return
		}
	}

	notfoundHandlers := router.notfoundHandlers
	for _, g := range router.groups {
		if strings.HasPrefix(reqURI, g.prefix) {
			notfoundHandlers = g.notfoundHandlers
			break
		}
	}
	if len(notfoundHandlers) > 0 {
		ctx.handlers = append(router.middlewares, notfoundHandlers...)
		ctx.Next()
	} else {
		ctx.handlers = router.middlewares
		ctx.NotFound()
	}
}

// Run run for http
func (router *Router) Run(addr string) {
	if addr == "" {
		addr = ":5000"
	}
	debugPrint("Listening and serving HTTP on %s\n", addr)
	http.ListenAndServe(addr, router)
}

// Use use for middleware
func (router *Router) Use(handler ...HandlerFunc) {
	router.middlewares = append(router.middlewares, handler...)
}
