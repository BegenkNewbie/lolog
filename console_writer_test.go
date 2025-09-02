package log

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConsoleWriter(t *testing.T) {
	cns := newConsoleWriter(DebugLevel)

	assert.Equal(t, DebugLevel, cns.Level())
	assert.Equal(t, os.Stdout, cns.GetWriter())
	assert.Equal(t, CONSOLE, cns.GetOutput())

	cns.Flush() // do nothing
}
