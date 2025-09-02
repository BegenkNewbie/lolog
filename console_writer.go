package log

import (
	"io"
	"os"
)

// newConsoleWriter return Writer implementer that write logs to os.Stdout and
// set given lvl as the log Level.
func newConsoleWriter(lvl Level) Writer {
	return &consoleOutput{lvl: lvl}
}

type consoleOutput struct {
	lvl Level
}

func (c *consoleOutput) GetWriter() io.Writer { return os.Stdout }
func (c *consoleOutput) GetOutput() Output    { return CONSOLE }
func (c *consoleOutput) Level() Level         { return c.lvl }
func (c *consoleOutput) Flush()               {}
