// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
//
// CDDLogger implements JSON-formatted rotating log files with optional upload callback,
// mirroring the Python SDK's CDDLogHandler and CustomRotatingFileHandler.
package cddlogger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	LogFileMaxBytes    = 500 * 1000 // Host API spec: 500 KB per log file
	LogFileRotateCount = 3          // Number of rotated files to keep
)

// LogRecord matches the TR-12 host API log format.
type LogRecord struct {
	Timestamp string `json:"timestamp"`
	DeviceID  string `json:"device_id"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Pathname  string `json:"pathname"`
	Lineno    int    `json:"lineno"`
	Exception string `json:"exception,omitempty"`
}

// UploadCallback is called after log rotation with the path to the rotated file.
type UploadCallback func(logFilePath string)

// CDDLogger manages structured JSON logging with rotation and upload.
type CDDLogger struct {
	mu             sync.Mutex
	deviceID       string
	logPath        string
	logFile        *os.File
	currentSize    int64
	uploadCallback UploadCallback
	stdLogger      *log.Logger
}

// New creates a new CDDLogger.
func New(logPath string, deviceID string, callback UploadCallback) (*CDDLogger, error) {
	l := &CDDLogger{
		deviceID:       deviceID,
		logPath:        logPath,
		uploadCallback: callback,
		stdLogger:      log.New(os.Stdout, "", 0),
	}
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return nil, fmt.Errorf("cannot create log directory %s: %w", logPath, err)
	}
	if err := l.openLogFile(); err != nil {
		return nil, err
	}
	return l, nil
}

func (l *CDDLogger) openLogFile() error {
	logFile := filepath.Join(l.logPath, "cdd_sdk.log")
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("cannot open log file %s: %w", logFile, err)
	}
	l.logFile = f
	info, _ := f.Stat()
	if info != nil {
		l.currentSize = info.Size()
	}
	return nil
}

// UpdateDeviceID changes the device ID used in log records.
func (l *CDDLogger) UpdateDeviceID(deviceID string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.deviceID = deviceID
}

// Dump forces a log rotation and triggers the upload callback.
func (l *CDDLogger) Dump() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.doRotation()
}

func (l *CDDLogger) writeRecord(level, message string, exception string) {
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	record := LogRecord{
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000000Z"),
		DeviceID:  l.deviceID,
		Level:     level,
		Message:   message,
		Pathname:  file,
		Lineno:    line,
		Exception: exception,
	}
	data, _ := json.Marshal(record)
	line_bytes := append(data, '\n')

	l.mu.Lock()
	defer l.mu.Unlock()

	// Print to stdout
	l.stdLogger.Print(string(data))

	// Write to file
	if l.logFile != nil {
		n, _ := l.logFile.Write(line_bytes)
		l.currentSize += int64(n)
		if l.currentSize >= LogFileMaxBytes {
			l.doRotation()
		}
	}
}

func (l *CDDLogger) doRotation() {
	if l.logFile == nil {
		return
	}
	l.logFile.Close()
	logFile := filepath.Join(l.logPath, "cdd_sdk.log")

	// Rotate files: .3 -> .4, .2 -> .3, .1 -> .2, current -> .1
	for i := LogFileRotateCount; i > 0; i-- {
		src := fmt.Sprintf("%s.%d", logFile, i)
		dst := fmt.Sprintf("%s.%d", logFile, i+1)
		os.Rename(src, dst)
	}
	os.Rename(logFile, logFile+".1")

	// Open new log file
	l.openLogFile()
	l.currentSize = 0

	// Trigger upload callback for the .1 file
	rotatedFile := logFile + ".1"
	if l.uploadCallback != nil {
		if _, err := os.Stat(rotatedFile); err == nil {
			go l.uploadCallback(rotatedFile)
		}
	}
}

// Info logs an informational message.
func (l *CDDLogger) Info(msg string) {
	l.writeRecord("INFO", msg, "")
}

// Infof logs a formatted informational message.
func (l *CDDLogger) Infof(format string, args ...interface{}) {
	l.writeRecord("INFO", fmt.Sprintf(format, args...), "")
}

// Error logs an error message.
func (l *CDDLogger) Error(msg string) {
	l.writeRecord("ERROR", msg, "")
}

// Errorf logs a formatted error message.
func (l *CDDLogger) Errorf(format string, args ...interface{}) {
	l.writeRecord("ERROR", fmt.Sprintf(format, args...), "")
}

// Warn logs a warning message.
func (l *CDDLogger) Warn(msg string) {
	l.writeRecord("WARNING", msg, "")
}

// Exception logs an error with exception details.
func (l *CDDLogger) Exception(msg string, err error) {
	exc := ""
	if err != nil {
		exc = err.Error()
	}
	l.writeRecord("ERROR", msg, exc)
}

// Close closes the log file.
func (l *CDDLogger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logFile != nil {
		l.logFile.Close()
	}
}

// Discard returns a no-op writer for suppressing other loggers.
func Discard() io.Writer {
	return io.Discard
}
