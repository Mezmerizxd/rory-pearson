package board

import (
	"rory-pearson/internal/board"
	"rory-pearson/pkg/server"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Initialize(server *server.Server) {
	server.Cfg.Log.Info().Msg("Initializing board controllers")

	server.Engine.GET("/api/board/get", func(c *gin.Context) {
		// Get query parameters with default values
		page, err := strconv.Atoi(c.DefaultQuery("page", "1")) // Default to page 1
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid page parameter",
			})
			return
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", strconv.Itoa(board.DefaultPageSize))) // Default to the board's DefaultPageSize
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid pageSize parameter",
			})
			return
		}

		posts, err := board.GetPosts(page, pageSize)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, posts)
	})

	server.Engine.POST("/api/board/create", func(c *gin.Context) {
		var body board.CreateBoardPost
		err := c.BindJSON(&body)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = board.CreatePost(body)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Post created",
		})
	})

}
