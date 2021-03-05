package router

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// Router router struct
type Router struct {
	tree     map[string]pathRuleSlice
	isSorted bool
}

// HandleFunc handle func
type HandleFunc func(ctx *Context)

// New new router
func New() Router {
	return Router{
		tree: make(map[string]pathRuleSlice),
	}
}
func (router *Router) addHandleFunc(method, path string, handle HandleFunc) {
	if !strings.HasPrefix(path, "/") {
		panic(fmt.Errorf("[%s] shold absolute path", path))
	}
	if _, ok := router.tree[method]; !ok {
		router.tree[method] = make([]pathRule, 0)
	}

	pr := newPathRule(path, &handle)
	for _, v := range router.tree[method] {
		if pr.isConflictsWith(&v) {
			panic(fmt.Errorf("[%s] conflicts with [%s]", pr.rule, v.rule))
		}
	}

	router.tree[method] = append(router.tree[method], pr)
}

// Get router get method
func (router *Router) Get(path string, handle HandleFunc) {
	router.addHandleFunc(http.MethodGet, path, handle)
}

// Post router post method
func (router *Router) Post(path string, handle HandleFunc) {
	router.addHandleFunc(http.MethodPost, path, handle)
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
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ruleMatchResult := router.match(r.Method, r.URL.Path)

	if nil != ruleMatchResult {
		ctx := Context{
			Request:         r,
			writer:          w,
			ruleMatchResult: *ruleMatchResult,
		}
		handle := *ruleMatchResult.rule.handle
		if nil != handle {
			handle(&ctx)
		} else {
			fmt.Println("warning: no handle for [" + ruleMatchResult.rule.rule + "]")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Run run for http
func (router *Router) Run(addr string) {
	http.ListenAndServe(addr, router)
}
