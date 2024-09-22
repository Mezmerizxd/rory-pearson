package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CompressDirectoryAndDelete(sourceDir, outputDir, zipName string) error {
	// Create the output ZIP file
	zipFilePath := filepath.Join(outputDir, zipName)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("could not create output ZIP file: %v", err)
	}
	defer zipFile.Close()

	// Create a new ZIP archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk the source directory and add files to the ZIP archive
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("could not walk directory: %v", err)
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Create a new file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("could not create file header: %v", err)
		}

		// Set the file name to be relative to the source directory
		header.Name, err = filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("could not set file name: %v", err)
		}

		// Write the file header
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("could not create file header: %v", err)
		}

		// Open the source file
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("could not open file: %v", err)
		}
		defer file.Close()

		// Copy the file to the ZIP archive
		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("could not copy file: %v", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking the directory: %v", err)
	}

	// Remove the source directory
	if err := os.RemoveAll(sourceDir); err != nil {
		return fmt.Errorf("could not remove source directory: %v", err)
	}

	return nil
}
