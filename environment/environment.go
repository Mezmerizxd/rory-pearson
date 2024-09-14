package environment

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

/*
########################################
########################################
#				 ENVIRENMENT VARIABLES				 #
########################################
# EDIT THE ENVIRONMENT VARIABLES BELOW #
########################################
*/
type Environment struct {
	ServerPort  string `json:"SERVER_PORT"`
	UIBuildPath string `json:"UI_BUILD_PATH"`
	DbUrl       string `json:"DB_URL"`
}

func Initialize() (*Environment, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("error loading .env file")
	}

	// Map Environment keys to a string array
	keys := getExpectedKeys()

	// Get Environment values
	envValues, err := getEnvValues(keys)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment values: %v", err)
	}

	// Create a map of keys and values
	envMap := make(map[string]interface{})
	for _, envValue := range envValues {
		envMap[envValue.key] = envValue.value
	}

	// Create a new Config struct
	env := Environment{}

	// Convert the map to a Config struct
	if err := mapToStruct(envMap, &env); err != nil {
		return nil, fmt.Errorf("failed to map environment values to struct: %v", err)
	}

	return &env, nil
}

var env *Environment

func Get() *Environment {
	return env
}

type EnvValues struct {
	key   string
	value string
}

func getEnvValues(key []string) ([]EnvValues, error) {
	var envValues []EnvValues

	for _, k := range key {
		v := os.Getenv(k)

		if v == "" {
			return nil, fmt.Errorf("environment variable %s is not set", k)
		}

		envValues = append(envValues, EnvValues{k, v})
	}

	return envValues, nil
}

// getExpectedKeys returns a list of all expected keys (JSON tags) in the Config struct.
func getExpectedKeys() []string {
	var keys []string
	val := reflect.TypeOf(Environment{})
	for i := 0; i < val.NumField(); i++ {
		keys = append(keys, strings.ToUpper(val.Field(i).Tag.Get("json")))
	}
	return keys
}

// mapToStruct converts a map to a Config struct using JSON marshaling.
func mapToStruct(m map[string]interface{}, s *Environment) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}
