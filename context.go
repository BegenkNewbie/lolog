package log

import (
	"context"
	"sync"
)

// loggerKeyType custom type for log type inside context.
type loggerKeyType int

// loggerKey identifier for logger inside context.
const loggerKey loggerKeyType = iota

// singletonLogger holder of Logger and intended to be used as singleton.
var singletonLogger Logger

// mutex to protect singletonLogger.
var mutex sync.Mutex

// WithCtx return a copy of ctx with given logger attached.
func WithCtx(ctx context.Context, w Logger) context.Context {
	if w == nil {
		return ctx
	}
	if ww, ok := ctx.Value(loggerKey).(Logger); ok {
		// do not store same Logger
		if ww == singletonLogger {
			return ctx
		}
	}
	return context.WithValue(ctx, loggerKey, w)
}

// FromCtx return the singleton Logger associated with given ctx. If no
// logger is associated, the default logger is returned, unless it is nil in
// which case a no-op logger is returned.
func FromCtx(ctx context.Context) Logger {
	if ww, ok := ctx.Value(loggerKey).(Logger); ok {
		// singleton already set
		if singletonLogger != nil {
			// and the value is same with the Logger inside context
			if ww == singletonLogger {
				return singletonLogger
			}
		}
		return ww
	}
	return NewNop()
}
