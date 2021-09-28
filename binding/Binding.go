package binding

import "net/http"

// IBinding Binding interface
type IBinding interface {
	Bind(*http.Request, interface{}) error
}

// IBindingBody bind body
type IBindingBody interface {
	IBinding
	BindBody([]byte, interface{}) error
}

var (
	// JSON json binding
	JSON = jsonBinding{}
)
