package environment

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

// Environment holds the environment variables required by the application.
// These are loaded from the .env file during initialization.
type Environment struct {
	ServerHost          string `json:"SERVER_HOST"`           // Host for the server to run on
	ServerPort          string `json:"SERVER_PORT"`           // Port for the server to run on
	UIBuildPath         string `json:"UI_BUILD_PATH"`         // Path to the UI build files
	DbUrl               string `json:"DB_URL"`                // Database connection URL
	SpotifyClientId     string `json:"SPOTIFY_CLIENT_ID"`     // Spotify client ID
	SpotifyClientSecret string `json:"SPOTIFY_CLIENT_SECRET"` // Spotify client secret
}

// Initialize loads the environment variables from the .env file and
// returns a populated Environment struct.
func Initialize() (*Environment, error) {
	environment := Environment{}

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("error loading .env file")
	}

	// Retrieve the expected environment keys (JSON tags from Environment struct)
	keys := getExpectedKeys()

	// Fetch environment values using the keys
	envValues, err := getEnvValues(keys)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment values: %v", err)
	}

	// Map environment values into a map[string]interface{}
	envMap := make(map[string]interface{})
	for _, envValue := range envValues {
		envMap[envValue.key] = envValue.value
	}

	// Create a new Environment struct from the map
	if err := mapToStruct(envMap, &environment); err != nil {
		return nil, fmt.Errorf("failed to map environment values to struct: %v", err)
	}

	env = &environment
	return &environment, nil
}

var env *Environment

// Get returns the current environment instance. It should be called after Initialize.
func Get() *Environment {
	return env
}

// EnvValues represents a key-value pair of environment variables.
type EnvValues struct {
	key   string // The environment variable key
	value string // The corresponding value of the key
}

// getEnvValues retrieves the environment variable values for the given keys.
// It returns an error if any required environment variable is not set.
func getEnvValues(keys []string) ([]EnvValues, error) {
	var envValues []EnvValues

	for _, k := range keys {
		v := os.Getenv(k)

		// If a required environment variable is missing, return an error
		if v == "" {
			return nil, fmt.Errorf("environment variable %s is not set", k)
		}

		envValues = append(envValues, EnvValues{k, v})
	}

	return envValues, nil
}

// getExpectedKeys returns a list of all expected environment variable keys.
// It extracts these keys from the JSON tags in the Environment struct.
func getExpectedKeys() []string {
	var keys []string
	val := reflect.TypeOf(Environment{})
	for i := 0; i < val.NumField(); i++ {
		keys = append(keys, strings.ToUpper(val.Field(i).Tag.Get("json")))
	}
	return keys
}

// mapToStruct converts a map of environment variable key-value pairs
// into an Environment struct using JSON marshaling and unmarshaling.
func mapToStruct(m map[string]interface{}, s *Environment) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}
