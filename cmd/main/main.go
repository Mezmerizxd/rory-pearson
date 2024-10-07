package main

import (
	"rory-pearson/controllers"
	"rory-pearson/environment"
	"rory-pearson/internal/background_remover"
	"rory-pearson/internal/board"
	"rory-pearson/internal/spotify"
	"rory-pearson/pkg/log"
	"rory-pearson/pkg/python"
	"rory-pearson/pkg/server"
)

func main() {
	// ========================================
	// 1. Logging Setup
	// ========================================
	// Initialize main logger with both console and file outputs.
	// Logs will be stored at the root directory defined in environment.
	l := log.New(log.Config{
		ID:            "main",
		ConsoleOutput: true,
		FileOutput:    true,
		StoragePath:   environment.RootDirectory,
	})
	defer l.Close() // Ensure the logger is closed properly when the function ends.

	// ========================================
	// 2. Environment Initialization
	// ========================================
	// Load environment variables and configuration settings.
	env, err := environment.Initialize()
	if err != nil {
		// If environment initialization fails, log the error and stop execution.
		l.Error().Err(err).Msgf("Failed to initialize environment: %s", err.Error())
		return
	}

	// ========================================
	// 3. Board System Initialization
	// ========================================
	// Initialize the board system with logging.
	err = board.Initialize(board.Config{
		Log: l, // Pass main logger to board for logging purposes.
	})
	if err != nil {
		// Log any error during board initialization and halt the program.
		l.Error().Err(err).Msg("Failed to initialize board")
		return
	}

	// ========================================
	// 4. Python Interpreter Setup
	// ========================================
	// Set up logging for Python-related operations.
	pyLog := log.New(log.Config{
		ID:            "python",
		ConsoleOutput: true,
		FileOutput:    true,
		StoragePath:   environment.RootDirectory,
	})
	defer pyLog.Close() // Ensure Python logger is closed after the function completes.

	// Initialize Python with required libraries (e.g., backgroundremover).
	_, err = python.Initialize(python.Config{
		Log:         pyLog, // Python-specific logging.
		StoragePath: environment.CreateStorageDirectory("python"),
		Librarys: []string{
			"backgroundremover", // Ensure necessary Python libraries are included.
		},
	})
	if err != nil {
		// Log any errors during Python initialization and stop the program.
		l.Error().Err(err).Msg("Failed to initialize python")
		return
	}

	// ========================================
	// 5. Background Remover Initialization
	// ========================================
	// Initialize the background remover service.
	_, err = background_remover.Initialize(background_remover.Config{
		Log:         l, // Pass the logger.
		StoragePath: environment.CreateStorageDirectory("background_remover"),
	})
	if err != nil {
		// Log any initialization errors for the background remover and halt execution.
		l.Error().Err(err).Msg("Failed to initialize background remover")
		return
	}

	// ========================================
	// 6. Spotify Manager Initialization
	// ========================================
	// Initialize Spotify Manager to handle session management.
	sm := spotify.Initialize(spotify.Config{
		Log: l, // Use logger.
	})
	defer sm.Close() // Ensure Spotify Manager cleans up resources on shutdown.

	// ========================================
	// 7. Server Initialization
	// ========================================
	// Set up logging for Server-related operations.
	srvLog := log.New(log.Config{
		ID:            "server",
		ConsoleOutput: true,
		FileOutput:    true,
		StoragePath:   environment.RootDirectory,
	})
	defer pyLog.Close() // Ensure Python logger is closed after the function completes.

	// Create a new server instance with the port and logger from environment settings.
	svr, err := server.New(server.Config{
		Port: env.ServerPort, // Dynamically set port from environment variables.
		Log:  srvLog,         // Pass server logger to the server.
	})
	if err != nil {
		// Log any server creation errors and stop execution.
		l.Error().Err(err).Msg("Failed to create server")
		return
	}

	// ========================================
	// 8. Controllers Setup
	// ========================================
	// Initialize application controllers to handle routing and API logic.
	controllers.Initialize(svr)

	// ========================================
	// 9. Serve Static UI Files
	// ========================================
	// Set up the server to serve the UI build directory, defined in the environment.
	svr.ServeUI(env.UIBuildPath)

	// ========================================
	// 10. Start Server
	// ========================================
	// Start the server and listen for incoming requests.
	// If the server fails to start, log the error.
	err = svr.Start()
	if err != nil {
		l.Error().Err(err).Msg("Failed to start server")
		return
	}
}
