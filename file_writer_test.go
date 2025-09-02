package log

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSetupFileWriter(t *testing.T) {
	t.Run("Valid Filename", func(t *testing.T) {
		// Arrange
		fileName := "./test.log"
		defer os.Remove(fileName) // Cleanup after test

		// Act
		writer, err := setupFileWriter(fileName)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, writer)
		assert.Equal(t, fileName, writer.filename)
		assert.NotNil(t, writer.file)

		// Check if the file was created
		_, statErr := os.Stat(fileName)
		assert.NoError(t, statErr)
	})

	t.Run("File Opening Error", func(t *testing.T) {
		// Arrange
		fileName := "/invalid/path/to/log/file.log"

		// Act
		writer, err := setupFileWriter(fileName)

		// Assert
		assert.Nil(t, writer)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not open log file")
	})
}
