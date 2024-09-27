package background_remover_controllers

import (
	"rory-pearson/internal/background_remover"
	"rory-pearson/pkg/server"

	"github.com/gin-gonic/gin"
)

func Initialize(server *server.Server) {
	server.Cfg.Log.Info().Msg("Initializing background remover controllers")

	server.Engine.POST("/api/background-remover", func(c *gin.Context) {
		server.Cfg.Log.Info().Msg("Background remover request")

		formFile, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		bg := background_remover.GetInstance()
		if bg == nil {
			c.JSON(500, gin.H{
				"error": "background remover not initialized",
			})
			return
		}

		storedFile, err := bg.Trigger(formFile)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.File(*&storedFile.FilePath)
		storedFile.RemoveFile()
	})
}
