package utils

import (
	"reflect"
	"runtime"
	"strings"
)

// GetHandlerName get handler name
func GetHandlerName(handler interface{}) string {
	if strings.HasPrefix(reflect.TypeOf(handler).String(), "func") {
		return runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	}
	return ""
}
