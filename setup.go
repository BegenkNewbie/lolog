package log

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

type loggerImpl struct {
	loggerConstructorFunc func(w ...Writer) Logger
	writers               []Writer
	level                 Level // default 0 (DebugLevel)
}

type LoggerOption func(*loggerImpl) error

// SetupLogWriter is the constructor using the options pattern
//
// If no options are provided, it defaults to DebugLevel and outputs to the console with ZapLogger
func SetupLogWriter(opts ...LoggerOption) (Logger, error) {
	l := &loggerImpl{
		level:                 DebugLevel,
		writers:               []Writer{newConsoleWriter(DebugLevel)},
		loggerConstructorFunc: newZapLogger,
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(l); err != nil {
			return nil, err
		}
	}

	wr := l.loggerConstructorFunc(l.writers...)
	wr.init()

	return wr, nil
}

// setupLogWriters configures the log writers based on the output destinations
func setupLogWriters(outputs []string, fileSetup *FileSetup, level Level) (out []Writer, err error) {
	for _, output := range outputs {
		switch output {
		case "console":
			out = append(out, newConsoleWriter(level))
		case "container":
			out = append(out, newContainerWriter(level))
		case "file":
			hostName, err := os.Hostname()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("get os.Hostname() failed: %v", err))
			}

			fileName := fmt.Sprintf("%s%s", fileSetup.PrefixFile, hostName)
			nfw, err := newFileWriter(level, fmt.Sprintf("%s%s.log", fileSetup.Path, fileName))
			if err != nil {
				return nil, err
			}
			out = append(out, nfw)
		}
	}
	return out, nil
}

// WithOutput configures the output destinations (console, file, container)
func WithOutput(output []string, fileSetup *FileSetup) LoggerOption {
	return func(l *loggerImpl) (err error) {
		outputAllowed := []string{"console", "file", "container"}
		for _, o := range output {
			if !slices.Contains(outputAllowed, o) {
				return fmt.Errorf("invalid output format '%s': only 'console', 'container' and 'file' allowed", o)
			}
		}
		l.writers, err = setupLogWriters(output, fileSetup, l.level)
		if err != nil {
			return err
		}
		return nil
	}
}

// WithLevel sets the logging level
func WithLevel(level string) LoggerOption {
	return func(l *loggerImpl) (err error) {
		l.level, err = adaptLogLevelFromString(level)
		if err != nil {
			return err
		}
		return nil
	}
}

// WithSlogLogger setup logger with Slog
func WithSlogLogger() func(*loggerImpl) error {
	return func(l *loggerImpl) error {
		l.loggerConstructorFunc = newSlogLogger
		return nil
	}
}

func adaptLogLevelFromString(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warning":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	}

	return DebugLevel, fmt.Errorf("invalid log level '%s'", lvl)
}
