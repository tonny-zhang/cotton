package httpserver

import (
	"fmt"
	"io"
	"os"
	"time"
)

// LoggerConf config for logger
type LoggerConf struct {
	Formatter func(param LoggerFormatterParam) string
	Writer    io.Writer
}

// LoggerFormatterParam param to formatter
type LoggerFormatterParam struct {
	Method     string
	StatusCode int
	TimeStamp  time.Time
	ClientIP   string
	Path       string
}

var defaultLogFormatter = func(param LoggerFormatterParam) string {
	return fmt.Sprintf("[HS-INFO] %v\t%13s %6s %3d %s\n",
		param.TimeStamp.Format("2006/01/02 15:04:05"),
		param.ClientIP,
		param.Method,
		param.StatusCode,
		param.Path,
	)
}

var defaultWriter io.Writer = os.Stdout

// Logger logger middleware
func Logger() HandlerFunc {
	return LoggerWidthConf(LoggerConf{})
}

// LoggerWidthConf get logger with config
func LoggerWidthConf(conf LoggerConf) HandlerFunc {
	formatter := conf.Formatter
	writer := conf.Writer
	if nil == formatter {
		formatter = defaultLogFormatter
	}
	if nil == writer {
		writer = defaultWriter
	}
	return func(ctx *Context) {
		ctx.Next()

		param := LoggerFormatterParam{
			Method:     ctx.Request.Method,
			Path:       ctx.Request.RequestURI,
			TimeStamp:  time.Now(),
			StatusCode: ctx.statusCode,
			ClientIP:   ctx.ClientIP(),
		}
		fmt.Fprint(writer, formatter(param))
	}
}
