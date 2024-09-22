package python

import (
	"errors"
	"os/exec"
	"rory-pearson/pkg/log"
	"strings"
)

const (
	PythonVirtualEnvDirectory      = "venv"
	PythonVirtualEnvActivateScript = "venv/bin/activate"
)

var (
	ErrorInstanceNotInitialized     = errors.New("python instance not initialized")
	ErrorCreatingVirtualEnvironment = errors.New("error creating virtual environment")
)

var instance *Python

type Config struct {
	Log      log.Log
	Librarys []string
}

type Python struct {
	Log           log.Log
	Librarys      []string
	IsInitialized bool
}

func Initialize(c Config) (*Python, error) {
	var python = &Python{
		Log:      c.Log,
		Librarys: c.Librarys,
	}

	// Get librarys
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

func GetInstance() (*Python, error) {
	if instance == nil {
		return nil, ErrorInstanceNotInitialized
	}
	return instance, nil
}

func (p *Python) Destroy() {
	// Delete the virtual environment
	cmd := exec.Command("rm", "-rf", PythonVirtualEnvDirectory)
	cmd.Run()

	p.Log.Info().Msg("Python instance destroyed")
}

func (p *Python) Command(arg ...string) (*exec.Cmd, error) {
	p.Log.Info().Msg("Creating command")

	// Check if environment is ready
	cmd := exec.Command("ls", PythonVirtualEnvDirectory)
	err := cmd.Run()
	if err != nil {
		p.Log.Info().Msg("Creating virtual environment")

		// Create the virtual environment if it doesn't exist
		cmd := exec.Command("python3", "-m", "venv", PythonVirtualEnvDirectory)
		err := cmd.Run()
		if err != nil {
			return nil, ErrorCreatingVirtualEnvironment
		}
	}

	// Join the arguments into a single string, properly escaping them
	commandStr := strings.Join(arg, " ")

	// Use bash to source the virtual environment and run the command
	mainCmd := exec.Command("bash", "-c", "source "+PythonVirtualEnvActivateScript+" && "+commandStr)

	mainCmd.Stdout = NewLoggerWriter(LoggerWriter{
		Log:  p.Log,
		Type: "info",
	})
	mainCmd.Stderr = NewLoggerWriter(LoggerWriter{
		Log:  p.Log,
		Type: "error",
	})

	// Activate the virtual environment and run any other command
	return mainCmd, nil
}

func (p *Python) doesLibraryExist(name string) bool {
	cmd, err := p.Command("pip show " + name)
	if err != nil {
		return false
	}

	err = cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func (p *Python) installLibrary(name string) error {
	cmd, err := p.Command("pip install " + name)
	if err != nil {
		return err
	}

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
