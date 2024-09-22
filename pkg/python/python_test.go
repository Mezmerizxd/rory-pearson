package python

import (
	"rory-pearson/pkg/log"
	"testing"
)

func TestInitialize(t *testing.T) {
	l := log.New(log.Config{
		ID:            "python_test",
		ConsoleOutput: true,
		FileOutput:    false,
	})

	_, err := Initialize(Config{
		Log: l,
		Librarys: []string{
			"backgroundremover",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDestroy(t *testing.T) {
	l := log.New(log.Config{
		ID:            "python_test",
		ConsoleOutput: true,
		FileOutput:    false,
	})

	py, err := Initialize(Config{
		Log: l,
		Librarys: []string{
			"backgroundremover",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	py.Destroy()
}

// go test -v -timeout 0 -run TestBackgroundRemover
func TestBackgroundRemover(t *testing.T) {
	l := log.New(log.Config{
		ID:            "python_test",
		ConsoleOutput: true,
		FileOutput:    false,
	})

	py, err := Initialize(Config{
		Log: l,
		Librarys: []string{
			"backgroundremover",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[TESTING] Background Initialized")

	t.Log("[TESTING] Triggering backgroundremover")
	cmd, err := py.Command("backgroundremover", "-i", "input.png", "-a", "-ae", "15", "-o", "output.png")
	if err != nil {
		t.Fatal(err)
	}
	cmd.Run()

	t.Log("[TESTING] Background Successfully Removed")
}
