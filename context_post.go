package cotton

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func (ctx *Context) initPostFormCache() {
	if nil == ctx.postFormCache {
		if nil != ctx.Request {
			if e := ctx.Request.ParseMultipartForm(defaultMultipartMemory); e != nil {
				if e != http.ErrNotMultipart {
					panic(e)
				}
			}
			ctx.postFormCache = ctx.Request.PostForm
		} else {
			ctx.postFormCache = url.Values{}
		}
	}
}

// GetAllPostForm get all post form value
func (ctx Context) GetAllPostForm() url.Values {
	ctx.initPostFormCache()
	return ctx.postFormCache
}

// GetPostForm get postform param
func (ctx *Context) GetPostForm(key string) string {
	ctx.initPostFormCache()
	if v, ok := ctx.postFormCache[key]; ok {
		return v[0]
	}
	return ""
}

// GetPostFormArray get postform param array
func (ctx *Context) GetPostFormArray(key string) []string {
	ctx.initPostFormCache()
	if v, ok := ctx.postFormCache[key]; ok {
		return v
	}
	return []string{}
}

// GetPostFormMap get postform param map
func (ctx *Context) GetPostFormMap(key string) (dicts map[string]string, exists bool) {
	ctx.initPostFormCache()
	return getValue(ctx.postFormCache, key)
}

// GetPostFormFile get postform file
func (ctx *Context) GetPostFormFile(key string) *multipart.FileHeader {
	list := ctx.GetPostFormFileArray(key)
	if len(list) > 0 {
		return list[0]
	}
	return nil
}

// GetPostFormFileArray get postform files
func (ctx *Context) GetPostFormFileArray(key string) (list []*multipart.FileHeader) {
	ctx.initPostFormCache()
	if ctx.Request.MultipartForm != nil {
		list, _ = ctx.Request.MultipartForm.File[key]
	}
	return
}

// SavePostFormFile save file
func (ctx *Context) SavePostFormFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	os.MkdirAll(filepath.Dir(dst), 0755)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
