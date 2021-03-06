package httpserver

import (
	"fmt"
	"httpserver/utils"
	"net/http"
	"sort"
	"strings"
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
	if indexStart >= 0 {
		middlewares := router.middlewares[:indexStart+1]
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
	}
	ctx.Next()
}
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ruleMatchResult := router.match(r.Method, r.URL.Path)

	if nil != ruleMatchResult {
		ctx := Context{
			Request:         r,
			writer:          w,
			index:           -1,
			ruleMatchResult: *ruleMatchResult,
		}
		handler := *ruleMatchResult.rule.handler
		if nil != handler {
			ctx.handlers = append(ctx.handlers, handler)
			router.runMiddleWare(ctx, *&ruleMatchResult.rule.middlewareHandlersIndex)
		} else {
			router.runMiddleWare(ctx, 0)
			fmt.Println("warning: no handler for [" + ruleMatchResult.rule.rule + "]")
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
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
