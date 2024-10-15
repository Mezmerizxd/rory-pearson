package controllers

import (
	"rory-pearson/controllers/background_remover"
	"rory-pearson/controllers/board"
	"rory-pearson/controllers/image_convert"
	"rory-pearson/controllers/spotify"
	"rory-pearson/pkg/server"
	"rory-pearson/pkg/util"
	"rory-pearson/plugins"

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

	// Handle Commands
	type CommandRequest struct {
		Command string `json:"command"`
		Args    []any  `json:"args"`
	}
	server.Engine.POST("/api/command", func(c *gin.Context) {
		/* Local IP's Only */
		if !server.IsLocalRequest(c) {
			c.JSON(403, gin.H{"error": "Forbidden"})
			return
		}

		var req CommandRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err := plugins.GetInstance().Commands.ExecuteCommand(req.Command, req.Args...)
		if err != nil {
			if err == plugins.ErrCommandNotFound {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}

			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Command executed successfully"})
	})
}
