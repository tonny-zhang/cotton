package main

import (
	"httpserver/router"
)

func main() {
	r := router.New()
	r.Get("/hello", func(ctx *router.Context) {
		ctx.String("hello get")
	})
	r.Get("/hello/", func(ctx *router.Context) {
		ctx.String("hello get2")
	})
	r.Get("/user/", func(ctx *router.Context) {
		ctx.String("/user")
	})
	r.Get("/user/:name", func(ctx *router.Context) {
		ctx.String("user name = " + ctx.GetQuery1("name"))
	})
	r.Get("/user/:id/:name", func(ctx *router.Context) {
		ctx.String("user id = " + ctx.GetQuery1("id") + " name =" + ctx.GetQuery1("name"))
	})
	r.Post("/hello", func(ctx *router.Context) {
		ctx.String("hello post")
	})

	r.Run(":5000")
}
