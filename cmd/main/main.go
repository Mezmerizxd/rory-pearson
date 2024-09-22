package main

import (
	"os"
	"rory-pearson/environment"
	"rory-pearson/internal/board"
	"rory-pearson/internal/controllers"
	"rory-pearson/pkg/log"
	"rory-pearson/pkg/python"
	"rory-pearson/pkg/server"
)

func main() {
	err := initDirectoryStructure()
	if err != nil {
		panic(err)
	}

	// Log
	l := log.New(log.Config{
		ID:            "main",
		ConsoleOutput: true,
		FileOutput:    false,
	})
	defer l.Close()

	// Environment variables
	env, err := environment.Initialize()
	if err != nil {
		l.Error().Err(err).Msgf("Failed to initialize environment", err.Error())
		return
	}

	// Board
	err = board.Initialize(board.Config{
		Log: l,
	})
	if err != nil {
		l.Error().Err(err).Msg("Failed to initialize board")
		return
	}

	// Python
	pyLog := log.New(log.Config{
		ID:            "python",
		ConsoleOutput: true,
		FileOutput:    true,
	})
	defer pyLog.Close()
	_, err = python.Initialize(python.Config{
		Log: pyLog,
		Librarys: []string{
			"backgroundremover",
		},
	})
	if err != nil {
		l.Error().Err(err).Msg("Failed to initialize python")
		return
	}

	// Server
	svr, err := server.New(server.Config{
		Port: env.ServerPort,
		Log:  l,
	})
	if err != nil {
		l.Error().Err(err).Msg("Failed to create server")
		return
	}

	// Controllers
	controllers.Initialize(svr)

	// Serve UI
	svr.ServeUI(env.UIBuildPath)

	// Start server
	err = svr.Start()
	if err != nil {
		l.Error().Err(err).Msg("Failed to start server")
		return
	}
}

var directorys = []string{"temp"}

func initDirectoryStructure() error {
	for _, directory := range directorys {
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			err := os.Mkdir(directory, 0755)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
