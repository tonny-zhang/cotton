package cotton

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/tonny-zhang/cotton/utils"
)

// HandlerFunc handler func
type HandlerFunc func(ctx *Context)

// Router router struct
type Router struct {
	prefix      string
	tree        map[string]pathRuleSlice
	middlewares []middleware
	countRouter int
	isSorted    bool
}

// NewRouter new router
func NewRouter() Router {
	return Router{
		tree:        make(map[string]pathRuleSlice),
		middlewares: make([]middleware, 0),
	}
}

// Group get group router
func (router *Router) Group(path string) Router {
	r := Router{
		prefix:      path,
		tree:        router.tree,
		middlewares: router.middlewares,
		countRouter: router.countRouter,
		isSorted:    router.isSorted,
	}
	return r
}
func (router *Router) addHandleFunc(method, path string, handler HandlerFunc) {
	if "" != router.prefix {
		path = utils.CleanPath(router.prefix + "/" + path)
	}
	if !strings.HasPrefix(path, "/") {
		panic(fmt.Errorf("[%s] shold absolute path", path))
	}
	if _, ok := router.tree[method]; !ok {
		router.tree[method] = make([]pathRule, 0)
	}

	pr := newPathRule(path, &handler)
	for _, v := range router.tree[method] {
		if pr.isConflictsWith(&v) {
			panic(fmt.Errorf("[%s] conflicts with [%s]", pr.rule, v.rule))
		}
	}

	pr.middlewareHandlersIndex = len(router.middlewares) - 1
	router.countRouter++
	router.tree[method] = append(router.tree[method], pr)
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

func (router *Router) match(method, path string) *matchResult {
	if !router.isSorted {
		for k := range router.tree {
			sort.Sort(router.tree[k])
		}
	}
	rules, ok := router.tree[method]

	if ok && len(rules) > 0 {
		for _, rule := range rules {
			rm := rule.match(path)
			if rm.IsMatch {
				return &rm
			}
		}
	}
	return nil
}
func (router *Router) runMiddleWare(ctx Context, indexStart int) {
	var middlewares []middleware
	if indexStart >= 0 {
		middlewares = router.middlewares[:indexStart+1]
	} else {
		// run all middlewares
		middlewares = router.middlewares
	}
	if len(middlewares) > 0 {
		for i, middleware := range middlewares {
			if i <= middleware.countRouter {
				for _, handler := range middleware.handlers {
					handler(&ctx)
				}
			}
		}
		return
	}
	ctx.Next()
}
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ruleMatchResult := router.match(r.Method, r.URL.Path)

	ctx := Context{
		Request:  r,
		Response: w,
		index:    -1,
	}
	if nil != ruleMatchResult {
		ctx.ruleMatchResult = *ruleMatchResult
		handler := *ruleMatchResult.rule.handler
		if nil != handler {
			ctx.handlers = append(ctx.handlers, handler)
		} else {
			fmt.Println("warning: no handler for [" + ruleMatchResult.rule.rule + "]")
		}

		router.runMiddleWare(ctx, *&ruleMatchResult.rule.middlewareHandlersIndex)
	} else {
		ctx.StatusCode(http.StatusNotFound)
		router.runMiddleWare(ctx, -1)
	}
}

// Run run for http
func (router *Router) Run(addr string) {
	debugPrint("Listening and serving HTTP on %s\n", addr)
	http.ListenAndServe(addr, router)
}

// Use use for middleware
func (router *Router) Use(handler ...HandlerFunc) {
	router.middlewares = append(router.middlewares, middleware{
		handlers:    append(make([]HandlerFunc, 0), handler...),
		countRouter: router.countRouter,
	})
}
