package main

import (
	"encoding/json"
	"fmt"
	"httpserver"
	"httpserver/utils"
	"net/http"
)

func main() {
	r := httpserver.NewRouter()

	// writer logger to file
	// f, e := os.OpenFile("1.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	// fmt.Println(f, e)
	// r.Use(httpserver.Logger(), httpserver.LoggerWidthConf(httpserver.LoggerConf{
	// 	Writer: f,
	// }))

	r.Get("/hello", nil)
	// r.Use(httpserver.Recover())
	r.Use(httpserver.RecoverWithWriter(nil, func(ctx *httpserver.Context, err interface{}) {
		strErr := ""
		switch err.(type) {
		case string:
			strErr = err.(string)
		case error:
			strErr = err.(error).Error()
		default:
			if b, err := json.Marshal(err); err == nil {
				strErr = string(b)
			}
		}
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.String("[500 error]" + strErr)
	}))
	r.Use(httpserver.Logger())
	r.Get("/panic", func(ctx *httpserver.Context) {
		// i := 0
		// fmt.Println(1 / i)
		panic([]int{1, 2})
	})
	r.Get("/hello/", func(ctx *httpserver.Context) {
		ctx.String("hello get2")
	})
	r.Use(httpserver.LoggerWidthConf(httpserver.LoggerConf{
		Formatter: func(param httpserver.LoggerFormatterParam) string {
			return fmt.Sprintf("[info] %s %s %s\t%d %s\n",
				utils.TimeFormat(param.TimeStamp),
				param.ClientIP, param.Method, param.StatusCode,
				param.Path,
			)

		},
	}))
	r.Get("/user/", func(ctx *httpserver.Context) {
		ctx.String("/user")
	})
	r.Get("/user/:name", func(ctx *httpserver.Context) {
		ctx.String("user name = " + ctx.Param("name"))
	})
	r.Get("/user/:id/:name", func(ctx *httpserver.Context) {
		ctx.String("user id = " + ctx.Param("id") + " name = " + ctx.Param("name"))
	})
	r.Post("/user/:id", func(ctx *httpserver.Context) {
		ctx.String("hello post " + ctx.Param("id"))
	})

	g1 := r.Group("/v1")
	{
		g1.Get("/a", func(ctx *httpserver.Context) {
			ctx.String("g1 a")
		})

		// <num> is short for {\d+}
		// /v1/b/123 			=>	match
		// /v1/b/113abc			=> 	no
		// /v1/c/113/abctest	=> 	no
		g1.Get("/b/:id<num>", func(ctx *httpserver.Context) {
			ctx.String("g1 b id = " + ctx.Param("id"))
		})
		// /v1/c/123-abc 		=>	match
		// /v1/c/113-abctest	=> 	no
		g1.Get("/c/:rule{\\d+-[a-d]+}", func(ctx *httpserver.Context) {
			ctx.String("g1 c rule = " + ctx.Param("rule"))
		})
		// /v1/c/123-abc/t 		=>	match
		// /v1/c/113-abc/test	=> 	no
		g1.Get("/d/:rule{\\d+-[a-d]+}/t", func(ctx *httpserver.Context) {
			ctx.String("g1 d rule = " + ctx.Param("rule"))
		})
	}
	g2 := r.Group("/v2/:method")
	{
		g2.Get("/a", func(ctx *httpserver.Context) {
			ctx.String("g2 a " + ctx.Param("method"))
		})
		g2.Get("/b", func(ctx *httpserver.Context) {
			ctx.String("g2 b " + ctx.Param("method"))
		})
		g2.Get("/c/:id", func(ctx *httpserver.Context) {
			ctx.String("g2 c " + ctx.Param("method") + " id = " + ctx.Param("id"))
		})
	}

	r.Run(":5000")
}
