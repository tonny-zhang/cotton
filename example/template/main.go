package main

import (
	"net/http"
	"os"

	"github.com/tonny-zhang/cotton"
)

func main() {
	root, _ := os.Getwd()

	router := cotton.NewRouter()
	router.Use(cotton.Logger(), cotton.RecoverWithWriter(os.Stdout, func(ctx *cotton.Context, err interface{}) {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		switch err.(type) {
		case error:
			ctx.Response.Write([]byte(err.(error).Error()))
		case string:
			ctx.Response.Write([]byte(err.(string)))
		}
	}))
	router.LoadTemplates(root, map[string]interface{}{
		"md5": func(str string) string {
			return str + "_md5"
		},
	})

	router.Get("/", func(ctx *cotton.Context) {
		ctx.Render("index", map[string]interface{}{
			"title": "index",
			"msg":   "hello from index",
		})
	})

	router.Run("")
}
