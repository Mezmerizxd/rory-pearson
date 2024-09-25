package main

import (
	"rory-pearson/controllers"
	"rory-pearson/environment"
	"rory-pearson/internal/background_remover"
	"rory-pearson/internal/board"
	"rory-pearson/pkg/log"
	"rory-pearson/pkg/python"
	"rory-pearson/pkg/server"
)

func main() {
	// Log
	l := log.New(log.Config{
		ID:            "main",
		ConsoleOutput: true,
		FileOutput:    false,
		StoragePath:   environment.RootDirectory,
	})
	defer l.Close()

	// Environment variables
	env, err := environment.Initialize()
	if err != nil {
		l.Error().Err(err).Msgf("Failed to initialize environment: %s", err.Error())
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
		ConsoleOutput: false,
		FileOutput:    true,
		StoragePath:   environment.RootDirectory,
	})
	defer pyLog.Close()

	_, err = python.Initialize(python.Config{
		Log:         pyLog,
		StoragePath: environment.CreateStorageDirectory("python"),
		Librarys: []string{
			"backgroundremover",
		},
	})
	if err != nil {
		l.Error().Err(err).Msg("Failed to initialize python")
		return
	}

	// Background remover
	bgLog := log.New(log.Config{
		ID:            "background_remover",
		ConsoleOutput: true,
		FileOutput:    true,
		StoragePath:   environment.RootDirectory,
	})
	_, err = background_remover.Initialize(background_remover.Config{
		Log:         bgLog,
		StoragePath: environment.CreateStorageDirectory("background_remover"),
	})

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
