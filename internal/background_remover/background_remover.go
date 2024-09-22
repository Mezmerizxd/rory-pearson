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

type BackgroundRemover struct {
	Log         log.Log
	StoragePath string
	Python      *python.Python

	JobsRunning   int
	JobsCompleted int
}

const (
	MaxConcurrentJobs = 5
)

var instance *BackgroundRemover

func Initialize(c Config) (*BackgroundRemover, error) {
	if instance != nil {
		return instance, nil
	}

	instance = &BackgroundRemover{
		Log:         c.Log,
		StoragePath: c.StoragePath,
	}

	p, err := python.GetInstance()
	if err != nil {
		return nil, err
	}
	instance.Python = p

	instance.Log.Info().Msg("Background remover initialized")

	return instance, nil
}

func GetInstance() *BackgroundRemover {
	return instance
}

type StoredFile struct {
	FileName string
	FilePath string
}

func (b *BackgroundRemover) Trigger(file *multipart.FileHeader) (*StoredFile, error) {
	b.Log.Info().Msg("Background remover request")

	if b.JobsRunning >= MaxConcurrentJobs {
		return nil, errors.New("max concurrent jobs reached")
	}

	b.JobsRunning++

	storedFile, err := b.TemporarilySaveFile(file)
	if err != nil {
		b.JobsRunning--
		return nil, err
	}

	modifiedFileName := "output_" + *&storedFile.FileName
	modifiedFilePath := filepath.Join(b.StoragePath, "temp", modifiedFileName)

	cmd, err := b.Python.Command("backgroundremover", "-i", storedFile.FilePath, "-a", "-ae", "15", "-o", modifiedFilePath)
	if err != nil {
		b.JobsRunning--
		return nil, err
	}
	cmd.Run()

	// Check if the file exists
	if _, err := os.Stat(modifiedFilePath); os.IsNotExist(err) {
		b.JobsRunning--
		return nil, errors.New("file not found")
	}

	err = storedFile.RemoveFile()
	if err != nil {
		return nil, err
	}

	b.JobsRunning--

	b.Log.Info().Msg("Background remover request completed")

	return &StoredFile{
		FileName: modifiedFileName,
		FilePath: modifiedFilePath,
	}, nil
}

func (b *BackgroundRemover) TemporarilySaveFile(file *multipart.FileHeader) (*StoredFile, error) {
	b.Log.Info().Msg("Temporarily saving file")

	// Check if `temp/` directory exists
	if _, err := os.Stat(filepath.Join(b.StoragePath, "temp")); os.IsNotExist(err) {
		b.Log.Info().Msg("Creating temp directory")
		err := os.MkdirAll(filepath.Join(b.StoragePath, "temp"), 0755)
		if err != nil {
			return nil, err
		}
	}

	// Generate a unique filename
	tempFileName := util.GenerateUUIDv4() + util.GetFileExtension(file.Filename)

	b.Log.Info().Msg("Saving file to temp directory")
	// Read the file
	fileData, err := file.Open()
	defer fileData.Close()
	if err != nil {
		return nil, err
	}

	b.Log.Info().Msg("Creating temp file")
	// Create the file
	tempFile, err := os.Create(filepath.Join(b.StoragePath, "temp", tempFileName))
	defer tempFile.Close()
	if err != nil {
		return nil, err
	}

	b.Log.Info().Msg("Writing to temp file")
	// Write the file
	_, err = tempFile.ReadFrom(fileData)
	if err != nil {
		return nil, err
	}

	b.Log.Info().Msg("File saved to temp directory")
	return &StoredFile{
		FileName: tempFileName,
		FilePath: filepath.Join(b.StoragePath, "temp", tempFileName),
	}, nil
}

func (s *StoredFile) RemoveFile() error {
	return os.Remove(s.FilePath)
}
