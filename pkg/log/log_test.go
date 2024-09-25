package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewLogger ensures that a new logger can be created and log messages to both console and file.
func TestNewLogger(t *testing.T) {
	// Define the configuration for the logger.
	cfg := Config{
		ID:            "test_logger",
		ConsoleOutput: true,
		FileOutput:    true,
		StoragePath:   "./test_logs/", // Use a test-specific log directory.
	}

	// Create the logger.
	logger := New(cfg)
	if logger == nil {
		t.Fatal("Expected non-nil logger, got nil")
	}

	// Log some messages.
	logger.Info().Msg("This is an info log")
	logger.Warn().Msg("This is a warning log")
	logger.Error().Msg("This is an error log")
	logger.Debug().Msg("This is a debug log")

	// Ensure log file is created.
	logFilePath := filepath.Join(cfg.StoragePath, "logs", fmt.Sprintf("%s.log", cfg.ID))
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		t.Fatalf("Expected log file to exist, but it does not: %s", logFilePath)
	}

	// Read the log file to verify content.
	logContent, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Error reading log file: %v", err)
	}
	if !strings.Contains(string(logContent), "This is an info log") {
		t.Errorf("Expected info log in file, got: %s", string(logContent))
	}
	if !strings.Contains(string(logContent), "This is a warning log") {
		t.Errorf("Expected warning log in file, got: %s", string(logContent))
	}

	// Close the logger.
	if err := logger.Close(); err != nil {
		t.Errorf("Error closing logger: %v", err)
	}

	// Clean up after test.
	cleanup(t, cfg.StoragePath)
}

// cleanup removes any test files or directories created during testing.
func cleanup(t *testing.T, path string) {
	if err := os.RemoveAll(path); err != nil {
		t.Errorf("Failed to clean up test files: %v", err)
	}
}

// TestGetLogger tests the Get function to ensure the same logger instance is retrieved.
func TestGetLogger(t *testing.T) {
	// Define the configuration for the logger.
	cfg := Config{
		ID:            "shared_logger",
		ConsoleOutput: true,
		FileOutput:    false,
		StoragePath:   "./test_logs/",
	}

	// Create the logger.
	logger := New(cfg)
	if logger == nil {
		t.Fatal("Expected non-nil logger, got nil")
	}

	// Fetch the same logger by ID.
	sameLogger := Get("shared_logger")
	if sameLogger == nil {
		t.Fatal("Expected to retrieve existing logger, got nil")
	}
	if sameLogger != logger {
		t.Errorf("Expected the same logger instance, got different instances")
	}

	// Log a message and verify the output.
	logger.Info().Msg("Log message from shared_logger")

	// Close the logger.
	if err := logger.Close(); err != nil {
		t.Errorf("Error closing logger: %v", err)
	}

	// Clean up after test.
	cleanup(t, cfg.StoragePath)
}

// TestLogClose ensures that the logger closes file handles properly.
func TestLogClose(t *testing.T) {
	cfg := Config{
		ID:            "close_test_logger",
		ConsoleOutput: false,
		FileOutput:    true,
		StoragePath:   "./test_logs/",
	}

	logger := New(cfg)
	if logger == nil {
		t.Fatal("Expected non-nil logger, got nil")
	}

	// Log a message.
	logger.Info().Msg("Testing close")

	// Close the logger.
	if err := logger.Close(); err != nil {
		t.Errorf("Error closing logger: %v", err)
	}

	// Attempt to log again (should not panic, but should not log either).
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logger should not panic after close, but did: %v", r)
		}
	}()
	logger.Info().Msg("Log after close")

	// Clean up after test.
	cleanup(t, cfg.StoragePath)
}
