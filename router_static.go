package cotton

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tonny-zhang/cotton/utils"
)

type onlyFilesFS struct {
	fs http.FileSystem
}
type noReaddirFile struct {
	http.File
}

func (file *noReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (fs *onlyFilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	info, _ := f.Stat()
	if info.IsDir() {
		return nil, os.ErrPermission
	}
	return &noReaddirFile{f}, nil
}

// StaticFile static file handler
//
// panic when prefix has ':' or '*'; and when open root error
//
// you can use `ctx.Param("filepath")` to get relativepath
func (router *Router) StaticFile(prefix, root string, listDir bool) {
	if strings.Index(prefix, "*") > -1 || strings.Index(prefix, ":") > -1 {
		panic(fmt.Errorf("static file prefix [%s] is illegal", prefix))
	}
	_, e := os.Open(root)
	if e != nil {
		panic(fmt.Errorf("static file root error, %s", e.Error()))
	}
	var fs http.FileSystem = http.Dir(root)
	if !listDir {
		fs = &onlyFilesFS{fs}
	}
	fileServer := http.StripPrefix(filepath.Join(router.prefix, prefix), http.FileServer(fs))

	router.Get(utils.CleanPath(prefix+"/*filepath"), func(ctx *Context) {
		fileServer.ServeHTTP(ctx.Response, ctx.Request)
	})
}
