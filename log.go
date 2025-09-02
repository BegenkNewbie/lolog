package log

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Log object that holds data for each field inserted to each log message. How
// Logger implementer is treating this object should read the field typ and
// follow the guideline from Type and each of the supported types.
type Log struct {
	typ Type
	key string
	str string
	num int
	flt float64
	b   bool
	any interface{}
	err error
}

// ErrorField is an object to display file and func name on stack trace key
type ErrorField struct {
	File string `json:"file"`
	Func string `json:"func"`
}

// String constructs a Log with the given key and value. This set the type
// to StringType.
func String(k, v string) Log {
	return Log{typ: StringType, key: k, str: v}
}

// Num constructs a Log with the given key and value. This set the type
// to NumType.
func Num(k string, num int) Log {
	return Log{typ: NumType, key: k, num: num}
}

// Float constructs a Log with the given key and value. This set the type
// to FloatType.
func Float(k string, f float64) Log {
	return Log{typ: FloatType, key: k, flt: f}
}

// Bool constructs a Log with the given key and value. This set the type
// to BoolType.
func Bool(k string, b bool) Log {
	return Log{typ: BoolType, key: k, b: b}
}

// Any constructs a Log with the given key and value. This set the type
// to AnyType.
func Any(k string, any interface{}) Log {
	return Log{typ: AnyType, key: k, any: any}
}

// Error constructs a Log with the given err value and 'error' as the key. This
// set the type to ErrorType.
func Error(err error) Log {
	file, line, funcName := trace()
	err = errors.Join(err, errors.New(fmt.Sprintf("\n\n%s:%d => %s", file, line, funcName)))
	return Log{typ: ErrorType, key: "error", err: err}
}

// StackTrace constructs a Log with to capture runtime stack trace and given 'stack_trace' as the key. This
// set the type to ErrorType.
func StackTrace() Log {
	return Log{typ: AnyType, key: "stack_trace", any: traces()}
}

// trace
// this is function to get single row of stack trace on runtime
func trace() (string, int, string) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "?", 0, "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return file, line, "?"
	}

	return file, line, fn.Name()
}

// traces
// this is function to get stack trace based on error variable from runtime
func traces() (fields []ErrorField) {
	var field ErrorField
	stackBuf := make([]uintptr, 6)
	length := runtime.Callers(3, stackBuf[:])
	s := stackBuf[:length]

	frames := runtime.CallersFrames(s)
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			field.File = fmt.Sprintf("%s:%d", frame.File, frame.Line)
			field.Func = frame.Function
			fields = append(fields, field)
		}
		if !more {
			break
		}
	}
	return
}

// Query
// is a wrapper function with 'query' as obj and separate between query and value.
// param: sql (string sql like `select * from blablabla`)
// param: values
// ex : {"query":{"statement":"select * from t where foo = ? and bar = ?", "values":[0,1]}}
func Query(sql string, values ...interface{}) Log {
	return Log{
		typ: AnyType,
		key: "query",
		any: map[string]interface{}{"statement": sql, "values": values},
	}
}

// Request
// is a wrapper function with 'request' as obj.
// param: value (anything)
// ex : {"request":{"hello":"world"}}
func Request(value any) Log {
	return Log{typ: AnyType, key: "request", any: value}
}

// Response
// is a wrapper function with 'response' as obj.
// param: value (anyting)
// ex : {"response":{"hello":"world"}}
func Response(value interface{}) Log {
	return Log{typ: AnyType, key: "response", any: value}
}

// Header
// is a wrapper function with 'header' as obj.
// param: value (anyting)
// ex : {"header":{"hello":"world"}}
func Header(values interface{}) Log {
	return Log{typ: AnyType, key: "header", any: values}
}
