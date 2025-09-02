package log

import (
	"io"
)

// Writer unified log writer that's responsible where the log from Logger
// should be written to.
type Writer interface {
	// Writer return where and how the implementer should write the logs.
	GetWriter() io.Writer
	// Output define the Output type.
	GetOutput() Output
	// Level define the logs level.
	Level() Level
	// Flush any necessary clean up task that will be run by Producer at the last order.
	Flush()
}
