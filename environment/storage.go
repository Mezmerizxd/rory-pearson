package environment

import (
	"fmt"
	"os"
)

const RootDirectory = "storage/" // Base directory for storage operations

// CreateStorageDirectory creates a directory under the root directory.
// It normalizes the input path by removing leading/trailing slashes.
// Returns the full path of the created directory.
func CreateStorageDirectory(path string) string {
	// Check if the path is empty
	if path == "" {
		panic("path is empty")
	}

	// Remove leading slash, if any
	if path[0] == '/' {
		path = path[1:]
	}

	// Remove trailing slash, if any
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// Construct the full directory path
	fullPath := fmt.Sprintf("%s%s", RootDirectory, path)

	// Create the directory if it does not exist
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			panic(err)
		}
	}

	return fullPath
}

// GetRootTempDirectory returns the path to the temporary directory under the root storage directory.
func GetRootTempDirectory() string {
	return fmt.Sprintf("%stemp/", RootDirectory)
}

// DestroyStorage removes the root storage directory and all its contents.
func DestroyStorage() {
	if err := os.RemoveAll(RootDirectory); err != nil {
		panic(err)
	}
}
