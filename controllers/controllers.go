package controllers

import (
	"rory-pearson/controllers/background_remover"
	"rory-pearson/controllers/board"
	"rory-pearson/controllers/image_convert"
	"rory-pearson/controllers/spotify"
	"rory-pearson/pkg/server"
	"rory-pearson/pkg/util"

	"github.com/gin-gonic/gin"
)

func Initialize(server *server.Server) {
	board.Initialize(server)
	image_convert.Initialize(server)
	background_remover.Initialize(server)
	spotify.Initialize(server)

	server.Engine.GET("/api/ping", func(c *gin.Context) {
		info, err := util.GetSystemInfo()
		c.JSON(200, gin.H{
			"message": info,
			"error":   err,
		})
	})
}
