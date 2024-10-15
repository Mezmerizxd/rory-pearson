package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"rory-pearson/controllers"
	"rory-pearson/environment"
	"rory-pearson/internal/background_remover"
	"rory-pearson/internal/board"
	"rory-pearson/internal/spotify"
	"rory-pearson/pkg/features"
	"rory-pearson/pkg/log"
	"rory-pearson/pkg/python"
	"rory-pearson/pkg/server"
	"rory-pearson/plugins"
	"strings"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure resources are cleaned up

	mainLogger := log.New(log.Config{
		ID:            "main",
		ConsoleOutput: true,
		FileOutput:    false,
		StoragePath:   "",
	})
	defer mainLogger.Close()

	mainLogger.Info().Msg("starting")

	// Environment
	_, err := environment.Initialize()
	if err != nil {
		mainLogger.Error().Err(err).Msg("failed to initialize environment")
		return
	}

	// Python
	_, err = python.Initialize(python.Config{
		Log: mainLogger,
	})
	if err != nil {
		mainLogger.Error().Err(err).Msg("failed to initialize python")
		return
	}

	// Plugins
	pl, err := plugins.Initialize(plugins.Config{
		Log: mainLogger,
	})
	if err != nil {
		panic(err)
	}
	defer pl.Close()

	// Features
	f := features.Initialize(features.Config{
		Log: mainLogger,
	})
	err = f.InitializeAll()
	if err != nil {
		panic(err)
	}

	initInternal(mainLogger)
	initServer(mainLogger)

	// Listen for termination signals (interrupt)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM) // Capture both interrupt and SIGTERM signals

	// ###################################

	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-interrupt:
			mainLogger.Info().Msg("received interrupt signal")
			return
		case <-ctx.Done():
			mainLogger.Info().Msg("context cancelled")
			return
		default:
			// Get input
			fmt.Printf("Enter command: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				mainLogger.Error().Err(err).Msg("error reading input")
				continue
			}

			// Trim the input to remove newline characters
			input = strings.TrimSpace(input)

			if input == "exit" {
				mainLogger.Info().Msg("exiting")
				return
			}

			// Split input into command and arguments
			inputParts := strings.Split(input, " ")
			command := inputParts[0]
			args := inputParts[1:]

			// Convert arguments to `any` type
			var parsedArgs []any
			for _, arg := range args {
				parsedArgs = append(parsedArgs, arg)
			}

			// Execute command with arguments
			err = pl.Commands.ExecuteCommand(command, parsedArgs...)
			if err != nil {
				mainLogger.Error().Err(err).Msgf("error executing command %s, reason: %s", command, err.Error())
			}
		}
	}
}

func initInternal(mainLogger log.Log) {
	// Board
	err := board.Initialize(board.Config{
		Log: mainLogger, // Pass main logger to board for logging purposes.
	})
	if err != nil {
		// Log any error during board initialization and halt the program.
		mainLogger.Error().Err(err).Msg("Failed to initialize board")
		return
	}

	// Background Remover
	_, err = background_remover.Initialize(background_remover.Config{
		Log:         mainLogger, // Pass the logger.
		StoragePath: environment.CreateStorageDirectory("background_remover"),
	})
	if err != nil {
		// Log any initialization errors for the background remover and halt execution.
		mainLogger.Error().Err(err).Msg("Failed to initialize background remover")
		return
	}

	// Spotify
	_ = spotify.Initialize(spotify.Config{
		Log: mainLogger, // Use logger.
	})
}

func initServer(mainLogger log.Log) {
	// Server
	svr, err := server.New(server.Config{
		Port: environment.Get().ServerPort, // Dynamically set port from environment variables.
		Log:  mainLogger,                   // Pass server logger to the server.
	})
	if err != nil {
		// Log any server creation errors and stop execution.
		mainLogger.Error().Err(err).Msg("Failed to create server")
		return
	}

	// Controllers
	controllers.Initialize(svr)

	// Setup ui
	svr.ServeUI(environment.Get().UIBuildPath)

	// Start the server
	err = svr.Start()
	if err != nil {
		mainLogger.Error().Err(err).Msg("Failed to start server")
		return
	}
}
