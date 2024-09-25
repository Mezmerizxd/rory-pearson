package background_remover

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"rory-pearson/pkg/log"
	"rory-pearson/pkg/python"
	"rory-pearson/pkg/util"
)

type Config struct {
	Log         log.Log
	StoragePath string
}

// BackgroundRemover manages background removal jobs and interacts with Python for processing.
type BackgroundRemover struct {
	Log         log.Log
	StoragePath string
	Python      *python.Python

	JobsRunning   int
	JobsCompleted int
}

const (
	// MaxConcurrentJobs defines the maximum number of background removal jobs that can run concurrently.
	MaxConcurrentJobs = 5
)

var instance *BackgroundRemover

// Initialize creates and returns a singleton instance of BackgroundRemover.
// It sets up the Python environment and logs the initialization.
func Initialize(c Config) (*BackgroundRemover, error) {
	// Check if the instance is already initialized
	if instance != nil {
		return instance, nil
	}

	// Initialize the BackgroundRemover instance
	instance = &BackgroundRemover{
		Log:         c.Log,
		StoragePath: c.StoragePath,
	}

	// Get the Python instance
	p, err := python.GetInstance()
	if err != nil {
		return nil, err
	}
	instance.Python = p

	// Log the initialization
	instance.Log.Info().Msg("Background remover initialized")

	return instance, nil
}

// GetInstance returns the singleton instance of BackgroundRemover.
func GetInstance() *BackgroundRemover {
	return instance
}

// StoredFile represents a file stored in the temporary directory for background removal.
type StoredFile struct {
	FileName string
	FilePath string
}

// Trigger handles a background removal request. It checks for max concurrent jobs,
// temporarily saves the uploaded file, triggers the Python background remover command,
// and deletes the original file after processing.
func (b *BackgroundRemover) Trigger(file *multipart.FileHeader) (*StoredFile, error) {
	b.Log.Info().Msg("Background remover request")

	// Ensure the maximum number of concurrent jobs is not exceeded
	if b.JobsRunning >= MaxConcurrentJobs {
		return nil, errors.New("max concurrent jobs reached")
	}

	b.JobsRunning++

	// Temporarily save the file
	storedFile, err := b.TemporarilySaveFile(file)
	if err != nil {
		b.JobsRunning--
		return nil, err
	}

	// Prepare the output file path
	modifiedFileName := "output_" + storedFile.FileName
	modifiedFilePath := filepath.Join(b.StoragePath, "temp", modifiedFileName)

	// Run the background remover command using Python
	cmd, err := b.Python.Command("backgroundremover", "-i", storedFile.FilePath, "-a", "-ae", "15", "-o", modifiedFilePath)
	if err != nil {
		b.JobsRunning--
		return nil, err
	}
	cmd.Run()

	// Check if the output file was successfully created
	if _, err := os.Stat(modifiedFilePath); os.IsNotExist(err) {
		b.JobsRunning--
		return nil, errors.New("file not found")
	}

	// Remove the original file
	err = storedFile.RemoveFile()
	if err != nil {
		return nil, err
	}

	b.JobsRunning--

	b.Log.Info().Msg("Background remover request completed")

	// Return the stored output file
	return &StoredFile{
		FileName: modifiedFileName,
		FilePath: modifiedFilePath,
	}, nil
}

// TemporarilySaveFile saves the uploaded file to a temporary directory.
// It returns a StoredFile with the file name and path.
func (b *BackgroundRemover) TemporarilySaveFile(file *multipart.FileHeader) (*StoredFile, error) {
	b.Log.Info().Msg("Temporarily saving file")

	// Ensure the temp directory exists
	tempDir := filepath.Join(b.StoragePath, "temp")
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		b.Log.Info().Msg("Creating temp directory")
		err := os.MkdirAll(tempDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	// Generate a unique file name
	tempFileName := util.GenerateUUIDv4() + util.GetFileExtension(file.Filename)

	// Open the uploaded file
	fileData, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileData.Close()

	// Create a file in the temp directory
	tempFilePath := filepath.Join(tempDir, tempFileName)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	// Write the file data to the temp file
	_, err = tempFile.ReadFrom(fileData)
	if err != nil {
		return nil, err
	}

	b.Log.Info().Msg("File saved to temp directory")

	return &StoredFile{
		FileName: tempFileName,
		FilePath: tempFilePath,
	}, nil
}

// RemoveFile deletes the file from the filesystem.
func (s *StoredFile) RemoveFile() error {
	return os.Remove(s.FilePath)
}
