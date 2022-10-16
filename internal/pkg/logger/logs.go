package logger

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

type Logger interface {
	WithContext(ctx context.Context) *log.Helper
	Log(level log.Level, keyvals ...interface{})
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
	Debugw(keyvals ...interface{})
	Info(a ...interface{})
	Infof(format string, a ...interface{})
	Infow(keyvals ...interface{})
	Warn(a ...interface{})
	Warnf(format string, a ...interface{})
	Warnw(keyvals ...interface{})
	Error(a ...interface{})
	Errorf(format string, a ...interface{})
	Errorw(keyvals ...interface{})
	Fatal(a ...interface{})
	Fatalf(format string, a ...interface{})
	Fatalw(keyvals ...interface{})
}

func New(id, name, version, level string) *log.Filter {
	loggerInstance := log.With(
		log.DefaultLogger,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", name,
		"service.version", version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	logLevel := log.ParseLevel(level)
	logger := log.NewFilter(loggerInstance, log.FilterLevel(logLevel))
	log.SetLogger(logger) // TODO Is it ok?
	return logger
}

func NewHelper(logger log.Logger, kv ...interface{}) *log.Helper {
	return log.NewHelper(log.With(logger, kv...))
}
