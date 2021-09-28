package cotton

import (
	"io/ioutil"

	"github.com/tonny-zhang/cotton/binding"
)

// ShouldBindWith bind
func (ctx *Context) ShouldBindWith(obj interface{}, b binding.IBinding) error {
	return b.Bind(ctx.Request, obj)
}

// BodyBytesKey indicates a default body bytes key.
const BodyBytesKey = "cotton/bbk"

// ShouldBindBodyWith bind body
func (ctx *Context) ShouldBindBodyWith(obj interface{}, bb binding.IBindingBody) (err error) {
	var body []byte
	if v, ok := ctx.Get(BodyBytesKey); ok {
		if vv, ok := v.([]byte); ok {
			body = vv
		}
	}
	if body == nil {
		body, err = ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return
		}
		ctx.Set(BodyBytesKey, body)
	}
	return bb.BindBody(body, obj)
}

// ShouldBindWithJSON bind with json
func (ctx *Context) ShouldBindWithJSON(obj interface{}) error {
	return ctx.ShouldBindWith(obj, binding.JSON)
}

// ShouldBindBodyWithJSON bind body with json
func (ctx *Context) ShouldBindBodyWithJSON(obj interface{}) (err error) {
	return ctx.ShouldBindBodyWith(obj, binding.JSON)
}
