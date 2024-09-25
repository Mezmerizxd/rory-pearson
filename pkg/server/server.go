package server

import (
	"net/http"
	"rory-pearson/pkg/log"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Config holds the configuration options for the server.
type Config struct {
	Port string  // The port on which the server will listen.
	Log  log.Log // Logger instance for logging server activities.
}

// Server represents the HTTP server and its configuration.
type Server struct {
	Mutex  sync.Mutex  // Mutex for safe concurrent access.
	Cfg    Config      // Server configuration.
	Engine *gin.Engine // The Gin engine for handling HTTP requests.
}

// New initializes a new server instance with the provided configuration.
func New(cfg Config) (*Server, error) {
	e := gin.Default()

	// Middleware for the server
	e.Use(gin.Recovery())    // Recover from panics and log the error.
	e.Use(gin.Logger())      // Logger middleware to log HTTP requests.
	e.Use(gin.ErrorLogger()) // Logger middleware for errors.
	e.Use(cors.Default())    // Enable CORS with default settings.

	return &Server{
		Cfg:    cfg,
		Engine: e,
	}, nil
}

// Start begins listening on the configured port for incoming HTTP requests.
func (s *Server) Start() error {
	s.HealthCheck() // Register the health check route.

	err := s.Engine.Run(":" + s.Cfg.Port)
	if err != nil {
		return err // Return error if the server fails to start.
	}

	return nil
}

// Stop stops the server gracefully.
func (s *Server) Stop() {
	// Implementation for stopping the server gracefully (optional)
}

// ServeUI serves static files for the UI.
func (s *Server) ServeUI(pathToBuild string) {
	s.Engine.Use(static.Serve("/", static.LocalFile(pathToBuild, true)))

	s.Engine.NoRoute(func(c *gin.Context) {
		c.File(pathToBuild + "/index.html") // Serve the index.html for unmatched routes.
	})
}

// HealthCheck is a simple route for health checking.
func (s *Server) HealthCheck() {
	s.Engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})
}
