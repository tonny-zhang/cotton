package main

import (
	"httpserver"
)

func main() {
	r := httpserver.New()
	// f, e := os.OpenFile("1.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	// fmt.Println(f, e)
	// r.Use(httpserver.Logger(), httpserver.LoggerWidthConf(httpserver.LoggerConf{
	// 	Writer: f,
	// }))

	r.Get("/hello", nil)
	r.Use(httpserver.Logger())
	r.Get("/hello/", func(ctx *httpserver.Context) {
		ctx.String("hello get2")
	})
	r.Use(httpserver.Logger())
	r.Get("/user/", func(ctx *httpserver.Context) {
		ctx.String("/user")
	})
	r.Get("/user/:name", func(ctx *httpserver.Context) {
		ctx.String("user name = " + ctx.Param("name"))
	})
	r.Get("/user/:id/:name", func(ctx *httpserver.Context) {
		ctx.String("user id = " + ctx.Param("id") + " name =" + ctx.Param("name"))
	})
	r.Post("/user/:id", func(ctx *httpserver.Context) {
		ctx.String("hello post " + ctx.Param("id"))
	})

	r.Run(":5000")
}
