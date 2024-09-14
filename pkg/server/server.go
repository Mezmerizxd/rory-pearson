package server

import (
	"rory-pearson/pkg/log"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port string
	Log  log.Log
}

type Server struct {
	Mutex  sync.Mutex
	Cfg    Config
	Engine *gin.Engine
}

func New(cfg Config) (*Server, error) {
	e := gin.Default()

	// Middleware
	e.Use(gin.Recovery())
	e.Use(gin.Logger())
	e.Use(gin.ErrorLogger())
	// e.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"*"},
	// 	AllowMethods: []string{"*"},
	// 	AllowHeaders: []string{"*"},
	// 	ExposeHeaders: []string{
	// 		"Content-Length",
	// 		"Access-Control-Allow-Origin",
	// 		"Access-Control-Allow-Headers",
	// 		"Access-Control-Allow-Methods",
	// 	},
	// 	AllowCredentials: true,
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return true
	// 	},
	// 	MaxAge: 12 * time.Hour,
	// }))
	e.Use(cors.Default())

	return &Server{
		Cfg:    cfg,
		Engine: e,
	}, nil
}

func (s *Server) Start() error {
	err := s.Engine.Run(":" + s.Cfg.Port)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() {}

// Serve the UI after all other routes
func (s *Server) ServeUI(pathToBuild string) {
	s.Engine.Use(static.Serve("/", static.LocalFile(pathToBuild, true)))

	s.Engine.NoRoute(func(c *gin.Context) {
		c.File(pathToBuild + "/index.html")
	})
}
