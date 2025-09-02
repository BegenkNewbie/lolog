package log

// Output define currently supported target output for logging.
type Output int8

const (
	CONSOLE   Output = iota // CONSOLE target log output to console/terminal with unstructured Level, Message, and Json Log
	FILE                    // FILE target log output to local file with structured data stdOut (json formatted)
	CONTAINER               // CONTAINER target log output to console/terminal with structured data stdOut (json formatted)
)

// A Level is a logging priority. Higher levels are more important.
type Level int8

const (
	// DebugLevel most verbose logs, and are usually disabled in production.
	DebugLevel Level = iota
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
)

// Type indicates how the Logger implementer should treat each
// Log.
type Type uint8

const (
	// StringType use field str string of Log as the value.
	StringType Type = iota
	// NumType use field num int of Log as the value.
	NumType
	// FloatType use field flt float64 of Log as the value.
	FloatType
	// BoolType use field b bool of Log as the value.
	BoolType
	// AnyType use field any interface of Log as the value.
	AnyType
	// ErrorType use field err from error interface of Log as the value.
	ErrorType
)
