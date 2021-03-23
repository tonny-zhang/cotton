package cotton

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)

func getTplFiles(dir string, ext string) (list []string, err error) {
	files, err := ioutil.ReadDir(dir)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		path := filepath.Join(dir, f.Name())
		if f.IsDir() {
			listSub, e := getTplFiles(path, ext)
			if e != nil {
				err = e
				return
			}
			list = append(listSub, list...)
		} else if filepath.Ext(f.Name()) == ext {
			list = append(list, path)
		}
	}
	return
}

// load template files
//
// funcs is functions register to template
// example:
// 	router.LoadTemplates(root, map[string]interface{}{
// 		"md5": func(str string) string {
// 			return str + "_md5"
// 		},
// 	})
func (router *Router) LoadTemplates(tplRoot string, funcs map[string]interface{}) {
	if router.globalTemplate == nil {
		tpl := template.New("global")
		if funcs != nil && len(funcs) > 0 {
			tpl.Funcs(funcs)
		}

		list, e := getTplFiles(tplRoot, ".html")
		if e != nil {
			panic(e)
		}
		tpl.ParseFiles(list...)
		router.globalTemplate = tpl
	}
}
