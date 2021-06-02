package cotton

import "fmt"

// Domain support domain
func (router *Router) Domain(domain string, handler ...HandlerFunc) *Router {
	if router.prefix != "" {
		panic(fmt.Errorf("group can not call Domain"))
	}
	if router.domain != "" {
		panic(fmt.Errorf("Domain can not call Domain"))
	}
	if _, ok := router.domains[domain]; ok {
		panic(fmt.Errorf("domain [%s] is exists", domain))
	}
	r := NewRouter()
	r.middlewares = router.middlewares
	r.domain = domain

	for _, h := range handler {
		if h != nil {
			r.middlewares = append(r.middlewares, h)
		}
	}
	router.domains[domain] = r
	return r
}
