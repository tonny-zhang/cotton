package main

import (
	"fmt"
	"net/http"

	"github.com/tonny-zhang/cotton"
)

func main() {
	r := cotton.NewRouter()
	r.Get("/get", func(ctx *cotton.Context) {
		name := ctx.GetQuery("name")
		first := ctx.GetDefaultQuery("first", "first default value")

		ids := ctx.GetQueryArray("ids[]")
		m, _ := ctx.GetQueryMap("info")
		ctx.String(http.StatusOK, fmt.Sprintf("name = %s, first = %s, ids = %v, info = %v", name, first, ids, m))
	})

	r.Run("")
}
