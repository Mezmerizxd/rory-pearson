// server_test.go
package server

import (
	"net/http"
	"testing"
	"time"

	"rory-pearson/pkg/log"
)

func getLogger() log.Log {
	return log.New(log.Config{
		ID:            "server_test",
		ConsoleOutput: false,
		FileOutput:    false,
		StoragePath:   "",
	})
}

// TestNewServer tests the creation of a new server instance.
func TestNewServer(t *testing.T) {
	logger := getLogger()
	config := Config{
		Port: "8080",
		Log:  logger,
	}

	srv, err := New(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if srv.Cfg.Port != config.Port {
		t.Fatalf("Expected port %s, got: %s", config.Port, srv.Cfg.Port)
	}
}

// TestServerStart tests starting the server.
func TestServerStart(t *testing.T) {
	logger := getLogger()
	config := Config{
		Port: "8080",
		Log:  logger,
	}

	srv, err := New(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			t.Errorf("Expected no error while starting server, got: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Send a test request to the health check endpoint
	resp, err := http.Get("http://localhost:" + config.Port + "/health")
	if err != nil {
		t.Fatalf("Expected no error while sending request, got: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got: %v", resp.StatusCode)
	}

	// Stop the server gracefully if implemented
	srv.Stop()
}
