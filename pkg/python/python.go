package python

import (
	"errors"
	"fmt"
	"os/exec"
	"rory-pearson/pkg/log"
	"strings"
)

const (
	// Virtual environment directory and activation script paths
	PythonVirtualEnvDirectory      = "venv"
	PythonVirtualEnvActivateScript = "venv/bin/activate"
)

var (
	// Common error messages
	ErrorInstanceNotInitialized     = errors.New("python instance not initialized")
	ErrorCreatingVirtualEnvironment = errors.New("error creating virtual environment")
)

var instance *Python

// Config holds the configuration for setting up a Python environment
type Config struct {
	Log         log.Log
	StoragePath string
	Librarys    []string
}

// Python struct manages the Python environment
type Python struct {
	Log           log.Log
	StoragePath   string
	Librarys      []string
	IsInitialized bool
}

// Initialize sets up a Python virtual environment and installs required libraries
func Initialize(c Config) (*Python, error) {
	var python = &Python{
		Log:         c.Log,
		StoragePath: c.StoragePath,
		Librarys:    c.Librarys,
	}

	// Ensure all required libraries are installed
	for _, library := range python.Librarys {
		if !python.doesLibraryExist(library) {
			c.Log.Info().Msg("Installing " + library + " library")

			err := python.installLibrary(library)
			if err != nil {
				return nil, err
			}
		}
	}

	c.Log.Info().Msg("Python instance initialized")
	python.IsInitialized = true
	instance = python
	return python, nil
}

// GetInstance returns the current Python instance, or an error if it's not initialized
func GetInstance() (*Python, error) {
	if instance == nil {
		return nil, ErrorInstanceNotInitialized
	}
	return instance, nil
}

// Destroy deletes the virtual environment and logs the action
func (p *Python) Destroy() {
	cmd := exec.Command("rm", "-rf", fmt.Sprintf("%s/%s", p.StoragePath, PythonVirtualEnvDirectory))
	cmd.Run()
	p.Log.Info().Msg("Python instance destroyed")
}

// Command creates a command to run within the Python virtual environment
func (p *Python) Command(arg ...string) (*exec.Cmd, error) {
	p.Log.Info().Msg("Creating command")

	// Check if virtual environment exists, and create it if necessary
	cmd := exec.Command("ls", fmt.Sprintf("%s/%s", p.StoragePath, PythonVirtualEnvDirectory))
	err := cmd.Run()
	if err != nil {
		p.Log.Info().Msg("Creating virtual environment")
		cmd := exec.Command("python3", "-m", "venv", fmt.Sprintf("%s/%s", p.StoragePath, PythonVirtualEnvDirectory))
		if err := cmd.Run(); err != nil {
			return nil, ErrorCreatingVirtualEnvironment
		}
	}

	// Build and execute the command within the virtual environment
	commandStr := strings.Join(arg, " ")
	mainCmd := exec.Command("bash", "-c", "source "+fmt.Sprintf("%s/%s", p.StoragePath, PythonVirtualEnvActivateScript)+" && "+commandStr)

	mainCmd.Stdout = NewLoggerWriter(LoggerWriter{Log: p.Log, Type: "info"})
	mainCmd.Stderr = NewLoggerWriter(LoggerWriter{Log: p.Log, Type: "error"})

	return mainCmd, nil
}

// doesLibraryExist checks if a Python library is installed in the virtual environment
func (p *Python) doesLibraryExist(name string) bool {
	cmd, err := p.Command("pip show " + name)
	if err != nil {
		return false
	}

	return cmd.Run() == nil
}

// installLibrary installs a Python library using pip
func (p *Python) installLibrary(name string) error {
	cmd, err := p.Command("pip install " + name)
	if err != nil {
		return err
	}

	return cmd.Run()
}
