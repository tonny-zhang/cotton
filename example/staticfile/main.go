package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tonny-zhang/cotton"
)

// //go:embed ui/*
// var uiEmbed embed.FS

func main() {
	dir, _ := os.Getwd()
	r := cotton.NewRouter()

	r.Use(cotton.Logger())
	// use custom static file
	r.Get("/v1/*file", func(ctx *cotton.Context) {
		file := filepath.Join(dir, ctx.Param("file"))

		http.ServeFile(ctx.Response, ctx.Request, file)
	})

	// use router.StaticFile
	r.StaticFile("/s/", dir, true)  // list dir
	r.StaticFile("/m/", dir, false) // 403 on list dir

	g := r.Group("/g/", func(ctx *cotton.Context) {
		fmt.Printf("status = %d param = %s, abspath = %s\n", ctx.Response.GetStatusCode(), ctx.Param("filepath"), filepath.Join(dir, ctx.Param("filepath")))
	})
	g.StaticFile("/", dir, true)

	// custom use for embed
	// r.Get("/ui/*file", func(ctx *cotton.Context) {
	// 	fileServer := http.StripPrefix("", http.FileServer(http.FS(uiEmbed)))
	// 	fileServer.ServeHTTP(ctx.Response, ctx.Request)
	// })

	r.Run("")
}
