package background_remover_controllers

import (
	"rory-pearson/internal/background_remover"
	"rory-pearson/pkg/server"

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

		// uuid := util.GenerateUUIDv4() + util.GetFileExtension(formFile.Filename)

		// // Save file to temp
		// filePath := environment.GetRootTempDirectory() + uuid
		// err = c.SaveUploadedFile(formFile, filePath)
		// if err != nil {
		// 	c.JSON(500, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		// p, err := python.GetInstance()
		// if err != nil {
		// 	c.JSON(500, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		// filePathOutput := environment.GetRootTempDirectory() + uuid + "_output.png"
		// cmd, err := p.Command("backgroundremover", "-i", filePath, "-a", "-ae", "15", "-o", filePathOutput)
		// if err != nil {
		// 	c.JSON(500, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }
		// cmd.Run()

		// c.File(filePathOutput)

		// // Remove the uploaded file
		// if err := os.Remove(filePath); err != nil {
		// 	c.JSON(500, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }
		// if err := os.Remove(filePathOutput); err != nil {
		// 	c.JSON(500, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }
	})
}
