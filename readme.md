# go http service

# how to use
```go
package main

import (
	"encoding/json"
	"fmt"
	"httpserver"
)

func main() {
	r := httpserver.NewRouter()
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
			return fmt.Sprintf("[info] %s %s\t%d %s\n",
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

	r.Run(":5000")
}

```

# todo list
