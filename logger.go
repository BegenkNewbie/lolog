package log

// Logger unified front-end to log.
type Logger interface {
	// Init do initialization and should be called at very first before
	// other functions.
	init()
	// Flush do clean up task that make sure any leftover logs written properly
	// to the destination by Log therefor should be called at very last
	// after other functions.
	Flush()
	// With add given Log(s) as structured context.
	With(pr ...Log) Logger
	// Group create new group or namespace with given key and Log(s) as the
	// content. In JSON, this is like creating new object using given key and
	// Log(s) as the fields.
	//
	// Ref(s):
	//  - https://pkg.go.dev/go.uber.org/zap#Namespace
	//  - https://pkg.go.dev/golang.org/x/exp/slog#Group
	Group(key string, pr ...Log) Logger
	// Dbg logs a message at DebugLevel.
	Dbg(msg string, pr ...Log)
	// Inf logs a message at InfoLevel.
	Inf(msg string, pr ...Log)
	// Wrn logs a message at WarnLevel.
	Wrn(msg string, pr ...Log)
	// Err logs a message at ErrorLevel.
	Err(msg string, pr ...Log)
}

// NewNop returns a no-op Logger. Do nothing and never writes out any logs.
func NewNop() Logger {
	return &nopLogger{}
}

type nopLogger struct{}

func (n nopLogger) init()                           {}
func (n nopLogger) Flush()                          {}
func (n nopLogger) With(_ ...Log) Logger            { return n }
func (n nopLogger) Group(_ string, _ ...Log) Logger { return n }
func (n nopLogger) Dbg(_ string, _ ...Log)          {}
func (n nopLogger) Inf(_ string, _ ...Log)          {}
func (n nopLogger) Wrn(_ string, _ ...Log)          {}
func (n nopLogger) Err(_ string, _ ...Log)          {}

type Logs []Log
