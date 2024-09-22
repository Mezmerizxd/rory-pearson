package image_convert_controllers

import (
	"os"
	"rory-pearson/internal/image_convert"
	"rory-pearson/pkg/server"

	"github.com/gin-gonic/gin"
)

func Initialize(server *server.Server) {
	server.Engine.POST("/image-convert/upload", func(c *gin.Context) {
		formFile, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Save file to temp
		filePath := "temp/" + formFile.Filename
		err = c.SaveUploadedFile(formFile, filePath)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Convert image to icon
		uuid, err := image_convert.Convert(filePath)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Remove the uploaded file
		if err := os.Remove(filePath); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message":     "File uploaded and converted",
			"download_id": uuid,
		})
	})

	server.Engine.GET("/image-convert/download/:id", func(c *gin.Context) {
		id := c.Param("id")

		zipPath, err := image_convert.GetConvertedZipPath(id, true)
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.File(zipPath)
	})

}
