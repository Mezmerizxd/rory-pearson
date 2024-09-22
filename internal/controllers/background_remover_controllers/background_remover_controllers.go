package background_remover_controllers

import (
	"os"
	"rory-pearson/pkg/python"
	"rory-pearson/pkg/server"
	"rory-pearson/pkg/util"

	"github.com/gin-gonic/gin"
)

func Initialize(server *server.Server) {
	server.Engine.POST("/background-remover", func(c *gin.Context) {
		server.Cfg.Log.Info().Msg("Background remover request")

		formFile, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		uuid := util.GenerateUUIDv4() + util.GetFileExtension(formFile.Filename)

		// Save file to temp
		filePath := "temp/" + uuid
		err = c.SaveUploadedFile(formFile, filePath)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		p, err := python.GetInstance()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		filePathOutput := "temp/" + uuid + "_output.png"
		cmd, err := p.Command("backgroundremover", "-i", filePath, "-a", "-ae", "15", "-o", filePathOutput)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		cmd.Run()

		c.File(filePathOutput)

		// Remove the uploaded file
		if err := os.Remove(filePath); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := os.Remove(filePathOutput); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	})
}
