package python

import (
	"os"
	"path/filepath"
	"rory-pearson/pkg/log"
	"testing"
)

func getLogger() log.Log {
	return log.New(log.Config{
		ID:            "python_test",
		ConsoleOutput: false,
		FileOutput:    false,
		StoragePath:   "",
	})
}

func TestInitialize(t *testing.T) {
	tempDir := t.TempDir() // Create a temporary directory for testing
	logger := getLogger()

	config := Config{
		Log:         logger,
		StoragePath: tempDir,
		Librarys:    []string{"requests"}, // Example library to test installation
	}

	// Initialize Python instance
	pythonInstance, err := Initialize(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !pythonInstance.IsInitialized {
		t.Fatal("Expected Python instance to be initialized")
	}

	// Check if the virtual environment was created
	venvPath := filepath.Join(tempDir, PythonVirtualEnvDirectory)
	if _, err := os.Stat(venvPath); os.IsNotExist(err) {
		t.Fatalf("Expected virtual environment directory %s to exist, but it does not", venvPath)
	}
}

func TestGetInstance(t *testing.T) {
	tempDir := t.TempDir() // Create a temporary directory for testing
	logger := getLogger()

	config := Config{
		Log:         logger,
		StoragePath: tempDir,
		Librarys:    []string{"requests"},
	}

	// Initialize Python instance
	_, err := Initialize(config)
	if err != nil {
		t.Fatalf("Expected no error during initialization, got: %v", err)
	}

	// Now retrieve the initialized instance
	instance, err := GetInstance()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !instance.IsInitialized {
		t.Fatal("Expected Python instance to be initialized")
	}
}

func TestDestroy(t *testing.T) {
	tempDir := t.TempDir()
	logger := getLogger()

	config := Config{
		Log:         logger,
		StoragePath: tempDir,
		Librarys:    []string{"requests"},
	}

	// Initialize Python instance
	pythonInstance, err := Initialize(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Destroy the Python instance
	pythonInstance.Destroy()

	// Check if the virtual environment has been removed
	venvPath := filepath.Join(tempDir, PythonVirtualEnvDirectory)
	if _, err := os.Stat(venvPath); !os.IsNotExist(err) {
		t.Fatalf("Expected virtual environment directory %s to be removed, but it still exists", venvPath)
	}
}

func TestLibraryInstallation(t *testing.T) {
	tempDir := t.TempDir()
	logger := getLogger()

	config := Config{
		Log:         logger,
		StoragePath: tempDir,
		Librarys:    []string{"requests"}, // Example library
	}

	// Initialize Python instance
	_, err := Initialize(config)
	if err != nil {
		t.Fatalf("Expected no error during initialization, got: %v", err)
	}

	// Check if the library was installed
	pythonInstance, _ := GetInstance()
	if !pythonInstance.doesLibraryExist("requests") {
		t.Fatal("Expected 'requests' library to be installed, but it is not")
	}
}
