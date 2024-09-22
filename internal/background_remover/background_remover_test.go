package background_remover

import (
	"rory-pearson/environment"
	"rory-pearson/pkg/log"
	"testing"
)

func createLogger() log.Log {
	return log.New(log.Config{
		ID:            "background_remover_test",
		ConsoleOutput: true,
		FileOutput:    false,
		StoragePath:   environment.CreateStorageDirectory("temp_storage"),
	})
}

func TestBackgroundRemover_Trigger(t *testing.T) {
	// Logger
	l := createLogger()

	// BackgroundRemover
	_, err := Initialize(Config{
		Log:         l,
		StoragePath: environment.CreateStorageDirectory("temp_storage"),
	})
	if err != nil {
		t.Fatal(err)
	}
}
