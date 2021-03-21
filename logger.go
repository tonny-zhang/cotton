package cotton

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
	Latency    time.Duration
	ClientIP   string
	Path       string
}

var defaultLogFormatter = func(param LoggerFormatterParam) string {
	return fmt.Sprintf("[INFO] %v\t%13s %6s %3d %10v %s \n",
		param.TimeStamp.Format("2006/01/02 15:04:05"),
		param.ClientIP,
		param.Method,
		param.StatusCode,
		param.Latency,
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
		timeStart := time.Now()
		ctx.Next()

		param := LoggerFormatterParam{
			Method:     ctx.Request.Method,
			Path:       ctx.Request.RequestURI,
			TimeStamp:  time.Now(),
			StatusCode: ctx.statusCode,
			ClientIP:   ctx.ClientIP(),
		}

		param.Latency = param.TimeStamp.Sub(timeStart)
		fmt.Fprint(writer, formatter(param))
	}
}
