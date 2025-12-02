package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// BenchmarkLogger handles thread-safe CSV writing
type BenchmarkLogger struct {
	file   *os.File
	writer *csv.Writer
	mu     sync.Mutex
}

// NewBenchmarkLogger creates the file and writes headers if empty
func NewBenchmarkLogger(filename string) (*BenchmarkLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)

	// Check if file is empty to write header
	stat, _ := file.Stat()
	if stat.Size() == 0 {
		writer.Write([]string{"timestamp", "endpoint", "duration_ns", "size_bytes"})
		writer.Flush()
	}

	return &BenchmarkLogger{
		file:   file,
		writer: writer,
	}, nil
}

// Log writes a single record to the CSV
func (l *BenchmarkLogger) Log(endpoint string, duration time.Duration, size int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	record := []string{
		strconv.FormatInt(time.Now().Unix(), 10),
		endpoint,
		strconv.FormatInt(duration.Nanoseconds(), 10),
		strconv.Itoa(size),
	}

	l.writer.Write(record)
	l.writer.Flush()
}

// NetworkMonitorMiddleware intercepts the request to measure time and size
func NetworkMonitorMiddleware(logger *BenchmarkLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // Process request

		duration := time.Since(start)
		responseSize := c.Writer.Size()

		// Write to CSV
		logger.Log(c.Request.URL.Path, duration, responseSize)

		// Set headers for debugging
		c.Writer.Header().Set("X-Duration", fmt.Sprintf("%v", duration))
		c.Writer.Header().Set("X-Size-Bytes", fmt.Sprintf("%d", responseSize))
	}
}
