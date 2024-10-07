package controllers

import (
	"rory-pearson/controllers/background_remover_controllers"
	"rory-pearson/controllers/board_controllers"
	"rory-pearson/controllers/image_convert_controllers"
	"rory-pearson/controllers/spotify"
	"rory-pearson/pkg/server"
	"rory-pearson/pkg/util"
	"rory-pearson/plugins"

	"github.com/gin-gonic/gin"
)

func Initialize(server *server.Server) {
	board_controllers.Initialize(server)
	image_convert_controllers.Initialize(server)
	background_remover_controllers.Initialize(server)
	spotify.Initialize(server)

	server.Engine.GET("/api/ping", func(c *gin.Context) {
		info, err := util.GetSystemInfo()
		c.JSON(200, gin.H{
			"message": info,
			"error":   err,
		})
	})

	server.Engine.GET("/api/v2/spotify/login", func(c *gin.Context) {
		state := util.GenerateUUIDv4()

		// plugins.GetInstance().Spotify.Auth.AuthURL(state)

		// debug plugins, im getting nil pointer error

		spotify := plugins.GetInstance().Spotify

		url := spotify.Auth.AuthURL(state)

		c.JSON(200, gin.H{"url": url})
	})
}
