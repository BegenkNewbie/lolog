package log

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// newFileWriter return Writer implementer that write logs to designated file
func newFileWriter(lvl Level, fileName string) (Writer, error) {
	fw, err := setupFileWriter(fileName)
	if err != nil {
		return nil, err
	}
	return &fileOutput{lvl: lvl, wr: fw}, nil
}

type fileOutput struct {
	wr  *fileWriter
	lvl Level
}

type fileWriter struct {
	filename string
	file     *os.File
	mu       sync.Mutex
}

func (f *fileOutput) GetWriter() io.Writer { return f.wr }
func (f *fileOutput) GetOutput() Output    { return FILE }
func (f *fileOutput) Level() Level         { return f.lvl }
func (f *fileOutput) Flush() {
	err := f.wr.Close()
	if err != nil {
		return
	}
}

func (c *fileWriter) Write(p []byte) (n int, err error) {
	// Check if the file still exists
	if _, err := os.Stat(c.filename); os.IsNotExist(err) {
		// File has been deleted; reopen it
		if err := c.openLogFile(); err != nil {
			return n, fmt.Errorf("failed to reopen log file: %w", err)
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Write data to file
	n, err = c.file.Write(p)
	if err != nil {
		return n, fmt.Errorf("write error: %w", err)
	}
	return n, nil
}

func (c *fileWriter) openLogFile() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close existing file if open
	if c.file != nil {
		c.file.Close()
	}

	// Open a new log file
	file, err := os.OpenFile(c.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open log file: %w", err)
	}
	c.file = file
	return nil
}

func (c *fileWriter) Close() interface{} {
	return c.file.Close()
}

// setupFileWriter init and set default value if no value
// provided in given config.
func setupFileWriter(fileName string) (*fileWriter, error) {
	lj := fileWriter{
		filename: fileName,
	}

	// Open initial log file
	if err := lj.openLogFile(); err != nil {
		return nil, err
	}

	return &lj, nil
}
