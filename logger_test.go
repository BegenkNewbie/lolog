package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNop(t *testing.T) {
	t.Run("Should have nop Logger type", func(t *testing.T) {
		nl := NewNop()
		assert.NotNil(t, nl)
		assert.Equal(t, &nopLogger{}, nl)
		assert.IsType(t, &nopLogger{}, nl)
	})
	t.Run("Should do nothing", func(t *testing.T) {
		nl := NewNop()
		assert.NotNil(t, nl)
		assert.NotNil(t, nl.With())
		assert.NotNil(t, nl.Group("group"))

		// just run it, since its just do literary nothing
		nl.init()
		nl.Flush()
		nl.Dbg("")
		nl.Inf("")
		nl.Wrn("")
		nl.Err("")
	})
}
