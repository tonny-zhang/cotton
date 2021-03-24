package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tonny-zhang/cotton"
)

func main() {
	fmt.Println(os.Args)
	fmt.Println()
	r := cotton.NewRouter()
	// Content-Type is "application/x-www-form-urlencoded" or "multipart/form-data"
	r.Post("/post", func(ctx *cotton.Context) {
		q := ctx.GetQuery("q")
		str := ctx.GetPostForm("str")
		ids := ctx.GetPostFormArray("ids")
		m, _ := ctx.GetPostFormMap("info")

		ctx.String(http.StatusOK, fmt.Sprintf("q = %s, str = %s, ids = %v, info = %v", q, str, ids, m))
	})

	dirSave := filepath.Dir(os.Args[0])
	r.Post("/upload", func(ctx *cotton.Context) {
		name := ctx.GetPostForm("name")
		fileheader := ctx.GetPostFormFile("file")
		files := ctx.GetPostFormFileArray("filelist")

		infos := make([]string, len(files))
		for i, file := range files {
			saveto := filepath.Join(dirSave, "upload", file.Filename)
			e := ctx.SavePostFormFile(file, saveto)
			if e == nil {
				infos[i] = file.Filename + " save to " + saveto
			} else {
				infos[i] = file.Filename + " has err: " + e.Error()
			}
		}

		ctx.String(http.StatusOK, "name = "+name+", file.Name = "+fileheader.Filename+"\n"+strings.Join(infos, "\n"))
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
