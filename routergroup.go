package httpserver

// RouterGroup group router
type RouterGroup struct {
	Router
	router *Router
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
