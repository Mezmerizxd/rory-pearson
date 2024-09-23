package image_convert_controllers

import (
	"os"
	"path/filepath"
	"rory-pearson/environment"
	"rory-pearson/internal/image_convert"
	"rory-pearson/pkg/server"

	"github.com/gin-gonic/gin"
)

func Initialize(server *server.Server) {
	server.Engine.POST("/api/image-convert/upload", func(c *gin.Context) {
		// Get the uploaded file
		formFile, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Save the file to a temp location
		tempFilePath := filepath.Join(environment.GetRootTempDirectory(), formFile.Filename)
		err = c.SaveUploadedFile(formFile, tempFilePath)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Convert the image to icons and compress them into a zip file
		uuid, err := image_convert.Convert(tempFilePath)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Remove the uploaded file after conversion
		if err := os.Remove(tempFilePath); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Return a UUID for the user to download later, no URL is returned now
		c.JSON(200, gin.H{
			"message":     "File uploaded and converted",
			"download_id": uuid, // UUID for download use
		})
	})

	server.Engine.GET("/api/image-convert/download/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Get the zip file path based on the provided ID
		zipPath, err := image_convert.GetConvertedZipPath(id) // File should not be deleted yet
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Serve the file directly to the user
		c.File(zipPath)

		// Delete the file after serving it
		if err := image_convert.DeleteConvertedZip(id); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	})

}
