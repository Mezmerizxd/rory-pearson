package main

import (
	"rory-pearson/environment"
	"rory-pearson/internal/board"
	"rory-pearson/internal/controllers"
	"rory-pearson/pkg/log"
	"rory-pearson/pkg/server"
)

func main() {
	// Log
	log := log.New(log.Config{
		ID:            "main",
		ConsoleOutput: true,
		FileOutput:    false,
	})

	// Environment variables
	env, err := environment.Initialize()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to initialize environment", err.Error())
		return
	}

	// Board
	err = board.Initialize(board.Config{
		Log: log,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize board")
		return
	}

	// Server
	svr, err := server.New(server.Config{
		Port: env.ServerPort,
		Log:  log,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create server")
		return
	}

	// Controllers
	controllers.Initialize(svr)

	// Serve UI
	svr.ServeUI(env.UIBuildPath)

	// Start server
	err = svr.Start()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start server")
		return
	}
}
