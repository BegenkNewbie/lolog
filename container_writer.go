package log

import (
	"io"
	"os"
)

// newContainerWriter return Writer implementer that write logs to os.Stdout
// and set given lvl as the log Level.
func newContainerWriter(lvl Level) Writer {
	return &containerOutput{lvl: lvl}
}

type containerOutput struct {
	lvl Level
}

func (c *containerOutput) GetWriter() io.Writer { return os.Stdout }
func (c *containerOutput) GetOutput() Output    { return CONTAINER }
func (c *containerOutput) Level() Level         { return c.lvl }
func (c *containerOutput) Flush()               {}
