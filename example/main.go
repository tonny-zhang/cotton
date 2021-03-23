package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"net/http"

	"github.com/tonny-zhang/cotton"
)

func main() {

	r := cotton.NewRouter()

	// writer logger to file
	// f, e := os.OpenFile("1.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	// fmt.Println(f, e)
	// r.Use(cotton.Logger(), cotton.LoggerWidthConf(cotton.LoggerConf{
	// 	Writer: f,
	// }))

	// r.Use(cotton.Recover())
	r.Use(cotton.Logger())
	// r.Use(func(ctx *cotton.Context) {
	// 	fmt.Println("first")
	// 	ctx.Abort()
	// })
	r.Use(cotton.RecoverWithWriter(nil, func(ctx *cotton.Context, err interface{}) {
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
		ctx.String(http.StatusInternalServerError, "[500 error]"+strErr)
	}))

	dir, _ := os.Getwd()
	r.Group("/static/", cotton.LoggerWidthConf(cotton.LoggerConf{
		Writer: os.Stdout,
		Formatter: func(param cotton.LoggerFormatterParam, ctx *cotton.Context) string {
			return fmt.Sprintf("[INFO-STATIC] %v\t %d %s\n",
				param.TimeStamp.Format("2006/01/02 15:04:05"),
				param.StatusCode,
				filepath.Join(dir, ctx.Param("file")),
			)
		},
	})).Get("/*file", func(ctx *cotton.Context) {
		// file := filepath.Join(dir, ctx.Param("file"))

		// http.ServeFile(ctx.Response, ctx.Request, file)
		// // ctx.Response.GetStatusCode() for log
		// fmt.Println(ctx.Response.GetStatusCode(), file)

		http.StripPrefix("/static/", http.FileServer(http.Dir(dir))).ServeHTTP(ctx.Response, ctx.Request)
		// http.StripPrefix("", http.FileServer(nil)).ServeHTTP()
	})
	gs := r.Group("/s/")
	gs.StaticFile("/", dir, false)
	r.StaticFile("/m/", dir, true)
	r.Get("/panic", func(ctx *cotton.Context) {
		// i := 0
		// fmt.Println(1 / i)
		panic([]int{1, 2})
	})
	r.Get("/hello/", func(ctx *cotton.Context) {
		ctx.String(http.StatusOK, "hello get2")
	})
	// r.Use(cotton.LoggerWidthConf(cotton.LoggerConf{
	// 	Formatter: func(param cotton.LoggerFormatterParam) string {
	// 		return fmt.Sprintf("[info] %s %s %s\t%d %s\n",
	// 			utils.TimeFormat(param.TimeStamp),
	// 			param.ClientIP, param.Method, param.StatusCode,
	// 			param.Path,
	// 		)

	// 	},
	// }))
	r.Get("/user/", func(ctx *cotton.Context) {
		ctx.String(http.StatusOK, "/user")
	})
	r.Get("/user/:name", func(ctx *cotton.Context) {
		ctx.String(http.StatusOK, "user name = "+ctx.Param("name"))
	})
	r.Get("/user/:name/:id", func(ctx *cotton.Context) {
		ctx.String(http.StatusOK, "user id = "+ctx.Param("id")+" name = "+ctx.Param("name"))
	})
	r.Get("/user/:name/:id/one", func(ctx *cotton.Context) {
		ctx.String(http.StatusOK, "one user id = "+ctx.Param("id")+" name = "+ctx.Param("name"))
	})
	r.Get("/user/:name/:id/two", func(ctx *cotton.Context) {
		ctx.String(http.StatusOK, "two user id = "+ctx.Param("id")+" name = "+ctx.Param("name"))
	})
	r.Get("/info/*file", func(ctx *cotton.Context) {
		ctx.String(http.StatusOK, "info file = "+ctx.Param("file"))
	})
	r.Post("/user/:id", func(ctx *cotton.Context) {
		ctx.String(http.StatusOK, "hello post "+ctx.Param("id"))
	})

	g1 := r.Group("/v1/", func(ctx *cotton.Context) {
		fmt.Println("g1 middleware")
	})
	g1.NotFound(func(ctx *cotton.Context) {
		ctx.String(http.StatusNotFound, "page ["+ctx.Request.RequestURI+"] not found")
	})
	{
		g1.Get("/a", func(ctx *cotton.Context) {
			ctx.String(http.StatusOK, "g1 a")
		})
		g1.Get("/info", func(ctx *cotton.Context) {
			ctx.JSON(http.StatusOK, cotton.M{
				"message": "from g1 info",
			})
		})
	}
	g2 := r.Group("/v2/")
	{
		g2.Get("/a", func(ctx *cotton.Context) {
			ctx.String(http.StatusOK, "g2 a "+ctx.Param("method"))
		})
		g2.Get("/b", func(ctx *cotton.Context) {
			ctx.String(http.StatusOK, "g2 b "+ctx.Param("method"))
		})
		g2.Get("/c/:id", func(ctx *cotton.Context) {
			ctx.String(http.StatusOK, "g2 c "+ctx.Param("method")+" id = "+ctx.Param("id"))
		})
	}

	g3 := r.Group("/v3/:method/")
	g3.Use(func(ctx *cotton.Context) {
		if ctx.Param("method") != "test" {
			ctx.Abort()
			ctx.String(http.StatusBadRequest, "no method test")
		}
	})
	{
		g3.Get("/a", func(ctx *cotton.Context) {
			ctx.String(http.StatusOK, "g3 a "+ctx.Param("method"))
		})
		g3.Get("/b", func(ctx *cotton.Context) {
			ctx.String(http.StatusOK, "g3 b "+ctx.Param("method"))
		})
		g3.Get("/c/:id", func(ctx *cotton.Context) {
			ctx.String(http.StatusOK, "g3 c "+ctx.Param("method")+" id = "+ctx.Param("id"))
		})
	}

	r.Group("/a/b")

	// r.PrintTree(http.MethodGet)
	r.Run(":5000")
}
