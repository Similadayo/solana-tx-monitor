package output

import (
	"encoding/csv"
	"os"
	"sync"
	"time"
)

type CSVOutput struct {
	file   *os.File
	writer *csv.Writer
	mu     sync.Mutex
}

func NewCSVOutput(filePath string) (*CSVOutput, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	writer := csv.NewWriter(file)
	// Write header if file is new
	if info, _ := file.Stat(); info.Size() == 0 {
		writer.Write([]string{"Timestamp", "Message"})
	}
	return &CSVOutput{file: file, writer: writer}, nil
}

func (c *CSVOutput) Write(msg string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.writer.Write([]string{time.Now().Format(time.RFC3339), msg})
	if err != nil {
		return err
	}
	c.writer.Flush()
	return c.writer.Error()
}

func (c *CSVOutput) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.writer.Flush()
	c.file.Close()
}
