package cotton

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/tonny-zhang/cotton/utils"
)

// HandlerFunc handler func
type HandlerFunc func(ctx *Context)

// Router router struct
type Router struct {
	srv         *http.Server
	prefix      string
	domain      string
	hasHandled  bool
	trees       map[string]*tree
	middlewares []HandlerFunc

	notfoundHandlers []HandlerFunc

	groups groupArr

	domains map[string]*Router

	globalTemplate *template.Template
}

// NewRouter new router
func NewRouter() *Router {
	return &Router{
		trees:       make(map[string]*tree),
		domains:     make(map[string]*Router),
		middlewares: make([]HandlerFunc, 0),
	}
}

// Default get default router
func Default() *Router {
	router := NewRouter()

	router.Use(Recover())
	router.Use(Logger())
	return router
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
	result := router.trees[method].root.find(path)
	if result.node != nil && result.node.isRealNode {
		panic(fmt.Errorf("[%s %s] has being setted", method, path))
	}
	nodeAdded := router.trees[method].add(path, nil)
	nodeAdded.middleware = append(nodeAdded.middleware, router.middlewares...)
	nodeAdded.middleware = append(nodeAdded.middleware, handler)
	router.hasHandled = true
	debugPrintRoute(method, router.domain+path, handler)
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

// ServeHTTP serve http handler
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(w, r, router)

	host := strings.Split(r.Host, ":")[0]

	if r, ok := router.domains[host]; ok {
		router = r
	}

	r.Method = strings.ToUpper(r.Method)
	reqURI := r.URL.Path

	if tree, ok := router.trees[r.Method]; ok {

		result := tree.root.find(reqURI)
		if result.node != nil {
			ctx.paramCache = result.params
			ctx.handlers = result.node.middleware

			ctx.Next()
			ctx.destroy()
			return
		}
	}

	routerUse := router.matchGroup(reqURI)
	if routerUse == nil {
		routerUse = router
	}

	notfoundHandlers := routerUse.notfoundHandlers
	middlewares := routerUse.middlewares

	if len(notfoundHandlers) > 0 {
		ctx.handlers = append(middlewares, notfoundHandlers...)
	} else {
		ctx.handlers = append(middlewares, func(ctx *Context) {
			ctx.NotFound()
		})
	}
	ctx.Next()
	ctx.destroy()
}

// Run run for http
func (router *Router) Run(addr string) error {
	if addr == "" {
		addr = ":5000"
	}
	for _, r := range router.domains {
		var groupsNew []*Router
		for _, g := range r.groups {
			if g.hasHandled {
				groupsNew = append(groupsNew, g)
			} else {
				debugPrint("group [%s] has no handler, will be discarded", g.prefix)
			}
		}
		r.groups = groupsNew
	}
	router.sort() // 对group进行排序
	debugPrint("Listening and serving HTTP on %s\n", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	router.srv = srv
	e := srv.ListenAndServe()

	return e
}

// Stop stop http service
func (router *Router) Stop(ctx context.Context) (e error) {
	if router.srv != nil {
		if ctx == nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
		}

		e = router.srv.Shutdown(ctx)

		debugPrint("exit http")
	}
	return
}

// Use use for middleware
func (router *Router) Use(handler ...HandlerFunc) {
	router.middlewares = append(router.middlewares, handler...)
}
