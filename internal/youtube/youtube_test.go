package youtube

import (
	"rory-pearson/environment"
	"testing"
)

func TestInitialize(t *testing.T) {

	// Initialize environment to load API keys from .env file
	environment.Initialize(&environment.Config{
		Filepath: "../../.env",
	})

	// Create a new instance of the Youtube struct
	youtube, err := Initialize(Config{})
	if err != nil {
		t.Fatalf("Error initializing YouTube service: %v", err)
	}

	video, err := youtube.SearchVideoWithHighestViews("On My Mind - Vibe Chemistry")
	if err != nil {
		t.Fatalf("Error searching for video: %v", err)
	}

	t.Logf("Video: %+v", video)
}
