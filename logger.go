package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// BenchmarkLogger handles thread-safe CSV writing
type BenchmarkLogger struct {
	filename string
	file     *os.File
	writer   *csv.Writer
	mutex    sync.Mutex
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
		WriteHeaders(writer)
	}

	return &BenchmarkLogger{
		filename: filename,
		file:     file,
		writer:   writer,
	}, nil
}

func WriteHeaders(writer *csv.Writer) {
	writer.Write([]string{"timestamp", "endpoint", "duration_ns", "size_bytes"})
	writer.Flush()
}

// Log writes a single record to the CSV
func (logger *BenchmarkLogger) Log(endpoint string, duration time.Duration, size int) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	// Check if file exists and is the same
	reopen := false

	if info, err := os.Stat(logger.filename); err != nil {
		if os.IsNotExist(err) {
			reopen = true
		}
	} else {
		if fdInfo, err := logger.file.Stat(); err == nil && !os.SameFile(info, fdInfo) {
			reopen = true
		}
	}

	if reopen {
		logger.file.Close()
		file, err := os.OpenFile(logger.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			logger.file = file
			logger.writer = csv.NewWriter(file)
			if stat, _ := file.Stat(); stat.Size() == 0 {
				WriteHeaders(logger.writer)
			}
		} else {
			fmt.Println("Error reopening log file:", err)
			return
		}
	}

	record := []string{
		strconv.FormatInt(time.Now().Unix(), 10),
		endpoint,
		strconv.FormatInt(duration.Nanoseconds(), 10),
		strconv.Itoa(size),
	}

	logger.writer.Write(record)
	logger.writer.Flush()
}

// Clear wipes the CSV content and rewrites headers
func (logger *BenchmarkLogger) Clear() error {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	if err := logger.file.Truncate(0); err != nil {
		return err
	}
	if _, err := logger.file.Seek(0, 0); err != nil {
		return err
	}

	logger.writer = csv.NewWriter(logger.file)
	logger.writer.Write([]string{"timestamp", "endpoint", "duration_ns", "size_bytes"})
	logger.writer.Flush()
	return nil
}

// DownloadHandler serves the CSV file, creates an snapshot that sends so we don't enter any race
// condition when calling from multiple devices
func (logger *BenchmarkLogger) DownloadHandler(c *gin.Context) {
	logger.mutex.Lock()
	logger.writer.Flush()

	tempFile, err := os.CreateTemp("", "benchmark_*.csv")
	if err != nil {
		logger.mutex.Unlock()
		c.String(500, "Error creating temp file")
		return
	}
	defer os.Remove(tempFile.Name())

	srcFile, err := os.Open(logger.filename)

	if err != nil {
		logger.mutex.Unlock()
		c.String(500, "Error creating temp file")
		return
	}

	io.Copy(tempFile, srcFile)

	srcFile.Close()
	tempFile.Close()

	logger.mutex.Unlock()

	c.FileAttachment(tempFile.Name(), "benchmark_data.csv")
}

// NetworkMonitorMiddleware intercepts the request to measure time and size
func NetworkMonitorMiddleware(logger *BenchmarkLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		responseSize := c.Writer.Size()

		path := c.Request.URL.Path

		// Trying to make a set out of a map
		isMonitoredPath := map[string]bool{"/optimal": true, "/slower": true}

		// Write to CSV
		if isMonitoredPath[path] {
			logger.Log(path, duration, responseSize)
		}

		// Set headers for debugging
		c.Writer.Header().Set("X-Duration", fmt.Sprintf("%v", duration))
		c.Writer.Header().Set("X-Size-Bytes", fmt.Sprintf("%d", responseSize))
	}
}
