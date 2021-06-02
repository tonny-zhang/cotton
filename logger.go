package cotton

import (
	"fmt"
	"io"
	"os"
	"time"
)

// LoggerConf config for logger
type LoggerConf struct {
	Formatter func(param LoggerFormatterParam, ctx *Context) string
	Writer    io.Writer
}

// LoggerFormatterParam param to formatter
type LoggerFormatterParam struct {
	Host       string
	Method     string
	StatusCode int
	TimeStamp  time.Time
	Latency    time.Duration
	ClientIP   string
	Path       string
}

var defaultLogFormatter = func(param LoggerFormatterParam, ctx *Context) string {
	return fmt.Sprintf("[INFO] %v\t%13s %6s %3d %10v %s \n",
		param.TimeStamp.Format("2006/01/02 15:04:05"),
		param.ClientIP,
		param.Method,
		param.StatusCode,
		param.Latency,
		param.Host+param.Path,
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
	debugPrint("use logger middleware")
	return func(ctx *Context) {
		timeStart := time.Now()
		ctx.Next()

		param := LoggerFormatterParam{
			Method:     ctx.Request.Method,
			Path:       ctx.Request.RequestURI,
			TimeStamp:  time.Now(),
			StatusCode: ctx.Response.GetStatusCode(),
			ClientIP:   ctx.ClientIP(),
			Host:       ctx.Request.Host,
		}

		param.Latency = param.TimeStamp.Sub(timeStart)
		fmt.Fprint(writer, formatter(param, ctx))
	}
}
