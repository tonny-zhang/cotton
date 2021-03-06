package httpserver

import (
	"bytes"
	"fmt"
	"httpserver/utils"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// HandlerFuncRecover handler func for recover
type HandlerFuncRecover func(ctx *Context, err interface{})

var defaultHandlerRecover = func(ctx *Context, err interface{}) {
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.writer.WriteHeader(http.StatusInternalServerError)
}

// Recover recover middleware
func Recover() HandlerFunc {
	return RecoverWithWriter(nil)
}

// RecoverWithWriter recover with wirter
func RecoverWithWriter(writer io.Writer, handler ...HandlerFuncRecover) HandlerFunc {
	if nil == writer {
		writer = defaultWriter
	}
	if len(handler) == 0 {
		handler = append(handler, defaultHandlerRecover)
	}
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				for _, handler := range handler {
					handler(ctx, err)
				}

				stack := stack(3)

				fmt.Fprintf(writer, "[HS-RECOVER] %s\n%s\n", utils.TimeFormat(time.Now()), stack)
			}
		}()

		ctx.Next()
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "\t%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
