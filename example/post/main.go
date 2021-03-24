package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tonny-zhang/cotton"
)

func main() {
	r := cotton.NewRouter()
	// Content-Type is "application/x-www-form-urlencoded" or "multipart/form-data"
	r.Post("/post", func(ctx *cotton.Context) {
		q := ctx.GetQuery("q")
		str := ctx.GetPostForm("str")
		ids := ctx.GetPostFormArray("ids")
		m, _ := ctx.GetPostFormMap("info")

		ctx.String(http.StatusOK, fmt.Sprintf("q = %s, str = %s, ids = %v, info = %v", q, str, ids, m))
	})

	// Content-Type is "application/json"
	r.Post("/json", func(ctx *cotton.Context) {
		ct := ctx.GetRequestHeader("Content-Type")
		if ct == "application/json" {
			body := ctx.Request.Body
			if body != nil {
				obj := make(map[string]interface{})
				json.NewDecoder(body).Decode(&obj)

				fmt.Println(obj)
			}
		}
	})

	r.Run("")
}
