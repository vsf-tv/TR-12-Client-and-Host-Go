package cddlogger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	dir := t.TempDir()
	l, err := New(dir, "test-device", nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer l.Close()

	logFile := filepath.Join(dir, "cdd_sdk.log")
	if _, err := os.Stat(logFile); err != nil {
		t.Fatalf("expected log file to exist: %v", err)
	}
}

func TestInfoWritesJSON(t *testing.T) {
	dir := t.TempDir()
	l, _ := New(dir, "dev-1", nil)
	defer l.Close()

	l.Info("hello world")
	l.logFile.Sync()

	data, _ := os.ReadFile(filepath.Join(dir, "cdd_sdk.log"))
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 0 {
		t.Fatal("expected at least one log line")
	}

	var rec LogRecord
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if rec.Level != "INFO" {
		t.Fatalf("expected INFO, got %q", rec.Level)
	}
	if rec.Message != "hello world" {
		t.Fatalf("expected 'hello world', got %q", rec.Message)
	}
	if rec.DeviceID != "dev-1" {
		t.Fatalf("expected dev-1, got %q", rec.DeviceID)
	}
}

func TestErrorAndErrorf(t *testing.T) {
	dir := t.TempDir()
	l, _ := New(dir, "dev-1", nil)
	defer l.Close()

	l.Error("something broke")
	l.Errorf("code %d", 42)
	l.logFile.Sync()

	data, _ := os.ReadFile(filepath.Join(dir, "cdd_sdk.log"))
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	var rec LogRecord
	json.Unmarshal([]byte(lines[0]), &rec)
	if rec.Level != "ERROR" || rec.Message != "something broke" {
		t.Fatalf("unexpected first record: %+v", rec)
	}
	json.Unmarshal([]byte(lines[1]), &rec)
	if rec.Level != "ERROR" || rec.Message != "code 42" {
		t.Fatalf("unexpected second record: %+v", rec)
	}
}

func TestInfof(t *testing.T) {
	dir := t.TempDir()
	l, _ := New(dir, "dev-1", nil)
	defer l.Close()

	l.Infof("count=%d name=%s", 5, "test")
	l.logFile.Sync()

	data, _ := os.ReadFile(filepath.Join(dir, "cdd_sdk.log"))
	var rec LogRecord
	json.Unmarshal([]byte(strings.TrimSpace(string(data))), &rec)
	if rec.Message != "count=5 name=test" {
		t.Fatalf("unexpected message: %q", rec.Message)
	}
}

func TestException(t *testing.T) {
	dir := t.TempDir()
	l, _ := New(dir, "dev-1", nil)
	defer l.Close()

	l.Exception("oops", os.ErrPermission)
	l.logFile.Sync()

	data, _ := os.ReadFile(filepath.Join(dir, "cdd_sdk.log"))
	var rec LogRecord
	json.Unmarshal([]byte(strings.TrimSpace(string(data))), &rec)
	if rec.Exception != "permission denied" {
		t.Fatalf("expected 'permission denied', got %q", rec.Exception)
	}
}

func TestUpdateDeviceID(t *testing.T) {
	dir := t.TempDir()
	l, _ := New(dir, "old-id", nil)
	defer l.Close()

	l.UpdateDeviceID("new-id")
	l.Info("after update")
	l.logFile.Sync()

	data, _ := os.ReadFile(filepath.Join(dir, "cdd_sdk.log"))
	var rec LogRecord
	json.Unmarshal([]byte(strings.TrimSpace(string(data))), &rec)
	if rec.DeviceID != "new-id" {
		t.Fatalf("expected new-id, got %q", rec.DeviceID)
	}
}

func TestLogRotation(t *testing.T) {
	dir := t.TempDir()
	var mu sync.Mutex
	var rotatedCalled bool
	callback := func(path string) {
		mu.Lock()
		rotatedCalled = true
		mu.Unlock()
	}
	_ = rotatedCalled
	l, _ := New(dir, "dev-1", callback)
	defer l.Close()

	// Write enough data to trigger rotation (LogFileMaxBytes = 500KB)
	bigMsg := strings.Repeat("x", 1000)
	for i := 0; i < 600; i++ {
		l.Info(bigMsg)
	}
	l.logFile.Sync()

	// Check that .1 file was created
	rotated := filepath.Join(dir, "cdd_sdk.log.1")
	if _, err := os.Stat(rotated); err != nil {
		t.Fatalf("expected rotated file %s to exist", rotated)
	}
}

func TestDump(t *testing.T) {
	dir := t.TempDir()
	l, _ := New(dir, "dev-1", nil)
	defer l.Close()

	l.Info("before dump")
	l.Dump()

	rotated := filepath.Join(dir, "cdd_sdk.log.1")
	if _, err := os.Stat(rotated); err != nil {
		t.Fatalf("expected rotated file after Dump: %v", err)
	}
}

func TestClose(t *testing.T) {
	dir := t.TempDir()
	l, _ := New(dir, "dev-1", nil)
	l.Close()
	// Should not panic on double close
	l.Close()
}
