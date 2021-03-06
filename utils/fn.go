package utils

import (
	"reflect"
	"runtime"
)

// GetHandlerName get handler name
func GetHandlerName(handler interface{}) string {
	defer func() {
		// TODO: recover 之后在外层的处理
		recover()
	}()
	// fmt.Println(reflect.TypeOf(handler).String())
	return runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
}
