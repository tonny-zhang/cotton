package cotton

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// IsDebugging is in debug mode
func IsDebugging() bool {
	return modeRuning == ModeDebug
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(defaultWriter, "[HS-Debug] "+format, values...)
	}
}

func debugPrintRoute(httpMethod, absolutePath string, handler HandlerFunc) {
	if IsDebugging() {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		debugPrint("%-6s %-25s --> %s\n", httpMethod, absolutePath, handlerName)
	}
}
