// plugins_test.go
package plugins

import (
	"rory-pearson/pkg/log"
	"testing"
)

func getLogger() log.Log {
	return log.New(log.Config{
		ID:            "plugins_test",
		ConsoleOutput: false,
		FileOutput:    false,
		StoragePath:   "",
	})
}

func TestInitialize(t *testing.T) {
	_, err := Initialize(Config{
		Log: getLogger(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

type PluginOneConfig struct{}

type PluginOne struct{}

func (p *PluginOne) Initialize(c PluginOneConfig) (*PluginOne, error) {
	return p, nil
}

func (p *PluginOne) DoSomething() {
	// Do something
}

func TestPluginOneInitialization(t *testing.T) {
	Initialize(Config{
		Log: getLogger(),
	})

	pluginOne := &PluginOne{}
	_, err := pluginOne.Initialize(PluginOneConfig{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if pluginOne == nil {
		t.Fatal("expected PluginOne to be initialized")
	}
}
