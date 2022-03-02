package binding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type jsonBinding struct{}

func (jsonBinding) Bind(req *http.Request, obj interface{}) (err error) {
	if req == nil || req.Body == nil {
		err = fmt.Errorf("bad request")
	} else {
		err = decodeJSON(req.Body, obj)
	}
	return
}
func (jsonBinding) BindBody(b []byte, obj interface{}) error {
	return decodeJSON(bytes.NewBuffer(b), obj)
}

func decodeJSON(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(obj)
}

var (
	// JSON json binding
	JSON IBindingBody = jsonBinding{}
)
