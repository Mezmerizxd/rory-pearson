package environment

import (
	"os"
	"testing"
)

// TestInitializeEnvironment checks that environment variables are loaded correctly
// and that the Initialize function creates an Environment struct.
func TestInitializeEnvironment(t *testing.T) {
	// Set temporary environment variables for testing
	os.Setenv("SERVER_HOST", "hostname")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("UI_BUILD_PATH", "/ui/build")
	os.Setenv("DB_URL", "postgres://user:pass@localhost/db")
	os.Setenv("SPOTIFY_CLIENT_ID", "client_id")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "client_secret")

	// Ensure environment variables are cleaned up after the test
	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("UI_BUILD_PATH")
		os.Unsetenv("DB_URL")
		os.Unsetenv("SPOTIFY_CLIENT_ID")
		os.Unsetenv("SPOTIFY_CLIENT_SECRET")
	}()

	// Initialize the environment
	env, err := Initialize(nil)
	if err != nil {
		t.Fatalf("failed to initialize environment: %v", err)
	}

	// Validate the environment values
	if env.ServerHost != "hostname" {
		t.Errorf("expected ServerHost 'hostname', got '%s'", env.ServerHost)
	}
	if env.ServerPort != "8080" {
		t.Errorf("expected ServerPort '8080', got '%s'", env.ServerPort)
	}
	if env.UIBuildPath != "/ui/build" {
		t.Errorf("expected UIBuildPath '/ui/build', got '%s'", env.UIBuildPath)
	}
	if env.DbUrl != "postgres://user:pass@localhost/db" {
		t.Errorf("expected DbUrl 'postgres://user:pass@localhost/db', got '%s'", env.DbUrl)
	}
	if env.SpotifyClientId != "client_id" {
		t.Errorf("expected SpotifyClientId 'client_id', got '%s'", env.SpotifyClientId)
	}
	if env.SpotifyClientSecret != "client_secret" {
		t.Errorf("expected SpotifyClientSecret 'client_secret', got '%s'", env.SpotifyClientSecret)
	}
}

// TestCreateStorageDirectory verifies that directories are created correctly.
func TestCreateStorageDirectory(t *testing.T) {
	// Test with a simple directory path
	path := CreateStorageDirectory("testdir")
	expectedPath := "storage/testdir"
	if path != expectedPath {
		t.Errorf("expected path '%s', got '%s'", expectedPath, path)
	}

	// Test with leading and trailing slashes
	path = CreateStorageDirectory("/testdir/test/")
	expectedPath = "storage/testdir/test"
	if path != expectedPath {
		t.Errorf("expected path '%s', got '%s'", expectedPath, path)
	}

	// Ensure that the directories created during the test are cleaned up
	defer func() {
		err := os.RemoveAll("storage/testdir")
		if err != nil {
			t.Fatalf("failed to clean up test directory: %v", err)
		}
	}()
}

// TestGetRootTempDirectory verifies that the root temp directory path is returned correctly.
func TestGetRootTempDirectory(t *testing.T) {
	expectedPath := "storage/temp/"
	if GetRootTempDirectory() != expectedPath {
		t.Errorf("expected '%s', got '%s'", expectedPath, GetRootTempDirectory())
	}
}

// TestDestroyStorage checks that the DestroyStorage function properly removes the root directory.
func TestDestroyStorage(t *testing.T) {
	// Create a directory
	CreateStorageDirectory("testdir")

	// Ensure the directory exists
	if _, err := os.Stat("storage/testdir"); os.IsNotExist(err) {
		t.Fatalf("expected directory 'storage/testdir' to exist")
	}

	// Destroy the storage
	DestroyStorage()

	// Ensure the directory is removed
	if _, err := os.Stat("storage/testdir"); err == nil || !os.IsNotExist(err) {
		t.Fatalf("expected directory 'storage/testdir' to be removed")
	}
}
