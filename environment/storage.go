package environment

import (
	"fmt"
	"os"
)

const RootDirectory = "storage/"

/*
CreateStorageDirectory creates a storage directory under the root directory
The path parameter could be "test", "/test", "test/", "/test/test2"
The function should return the full path of the created directory
*/
func CreateStorageDirectory(path string) string {
	// Check if the path is empty
	if path == "" {
		panic("path is empty")
	}

	// Check if the path starts with a slash
	if path[0] == '/' {
		path = path[1:]
	}

	// Check if the path ends with a slash
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// Create the full path
	fullPath := fmt.Sprintf("%s%s", RootDirectory, path)

	// Check if the directory exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Create the directory
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			panic(err)
		}
	}

	return fullPath
}

func GetRootTempDirectory() string {
	return fmt.Sprintf("%stemp/", RootDirectory)
}
