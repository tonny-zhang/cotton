package cotton

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tonny-zhang/cotton/utils"
)

type groupArr []*Router

func (s groupArr) Len() int      { return len(s) }
func (s groupArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s groupArr) Less(i, j int) bool {
	return len(strings.Split(s[i].prefix, "/")) > len(strings.Split(s[j].prefix, "/"))
}

// Group get group router
func (router *Router) Group(path string, handler ...HandlerFunc) *Router {
	if len(path) == 0 || path[0] != '/' {
		panic(fmt.Errorf("group [%s] must start with /", path))
	}
	if strings.Index(path, "*") > -1 || strings.Index(path, ":") > -1 {
		panic(fmt.Errorf("group path [%s] can not has parameter", path))
	}
	prefix := utils.CleanPath(path + "/")
	hasGroup := router.hasGroup(prefix)
	if hasGroup {
		panic(fmt.Errorf("group [%s] is setted", prefix))
	}
	if router.prefix != "" {
		prefix = utils.CleanPath(router.prefix + "/" + prefix)
	}
	middlewares := append([]HandlerFunc{}, router.middlewares...)
	middlewares = append(middlewares, handler...)
	r := &Router{
		prefix:           prefix,
		domain:           router.domain,
		trees:            router.trees,
		middlewares:      middlewares,
		notfoundHandlers: router.notfoundHandlers,
	}

	router.groups = append(router.groups, r)
	return r
}
func (router *Router) sort() {
	if len(router.groups) > 0 {
		for _, g := range router.groups {
			g.sort()
		}
		sort.Sort(router.groups)
	}
}
func (router *Router) hasGroup(path string) bool {
	for _, g := range router.groups {
		if len(g.groups) > 0 {
			has := g.hasGroup(path)
			if has {
				return true
			}
		}
		if g.prefix == path {
			return true
		}
	}
	return false
}
func (router *Router) matchGroup(path string) *Router {
	for _, g := range router.groups {
		if len(g.groups) > 0 {
			matchedGroup := g.matchGroup(path)
			if matchedGroup != nil {
				return matchedGroup
			}
		}
		if matchGroup(g, path) {
			return g
		}
	}
	return nil
}
func matchGroup(router *Router, path string) bool {
	if len(router.prefix) > 0 {
		if strings.HasPrefix(path, router.prefix) {
			return true
		}
		arrRP := strings.Split(router.prefix, "/")
		arrPath := strings.Split(path, "/")
		if len(arrPath) < len(arrRP) {
			return false
		}

		for i, j := 0, len(arrRP); i < j; i++ {
			if i == j-1 && arrRP[i] == "" {
				return true
			}
			if arrRP[i] != arrPath[i] {
				return false
			}
		}
	}
	return false
}
