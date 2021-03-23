package cotton

import "fmt"

func (ctx *Context) Render(tplName string, params map[string]interface{}) {
	t := ctx.router.globalTemplate.Lookup(tplName)
	if t == nil {
		panic(fmt.Errorf("no template [%s]", tplName))
	} else {
		e := t.Execute(ctx.Response, params)
		if e != nil {
			panic(e)
		}
	}
	ctx.Next()
}
