package controllers

import (
	"rory-pearson/internal/controllers/background_remover_controllers"
	"rory-pearson/internal/controllers/board_controllers"
	"rory-pearson/internal/controllers/image_convert_controllers"
	"rory-pearson/pkg/server"
	"rory-pearson/pkg/util"

	"github.com/gin-gonic/gin"
)

func Initialize(server *server.Server) {
	board_controllers.Initialize(server)
	image_convert_controllers.Initialize(server)
	background_remover_controllers.Initialize(server)

	server.Engine.GET("/ping", func(c *gin.Context) {
		info, err := util.GetSystemInfo()
		c.JSON(200, gin.H{
			"message": info,
			"error":   err,
		})
	})
}
