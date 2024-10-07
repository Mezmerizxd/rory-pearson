package spotify

import (
	"net/http"
	"rory-pearson/internal/users"
	"rory-pearson/pkg/features"
	"rory-pearson/pkg/server"
	"rory-pearson/pkg/util"
	"rory-pearson/plugins"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(server *server.Server) {
	f, _ := features.GetInstance().GetFeature(users.UsersFeatureType)
	usersFeature, _ := f.(*users.UsersFeature)

	server.Engine.GET("/api/spotify/login", func(c *gin.Context) {
		// Generate a new UUID for the state
		state := util.GenerateUUIDv4()

		// Generate the authorization URL with the state
		url := plugins.GetInstance().Spotify.Auth.AuthURL(state)

		// Store the state in the session
		usersFeature.SpotifyStoreUserAuth(c.Request.Context(), state, nil)

		c.JSON(http.StatusOK, gin.H{"url": url, "state": state})
	})

	server.Engine.GET("/api/spotify/callback", func(c *gin.Context) {
		state := c.Query("state")

		if state == "" {
			c.JSON(403, gin.H{"error": "State is required"})
			return
		}

		// Verify state matches
		user, err := usersFeature.SpotifyGetUserAuth(state)
		if user == nil {
			c.JSON(403, gin.H{"error": "Session not found"})
			return
		}

		token, err := plugins.GetInstance().Spotify.Auth.Token(c.Request.Context(), state, c.Request)
		if err != nil {
			c.JSON(403, gin.H{"error": "Couldn't get token"})
			return
		}

		usersFeature.SpotifyStoreUserAuth(c.Request.Context(), state, token)

		c.Redirect(http.StatusFound, "/spotify")
	})

	server.Engine.GET("/api/spotify/validate", func(c *gin.Context) {
		state := c.Query("state") // Use the state parameter to identify the session

		if state == "" {
			c.JSON(403, gin.H{"error": "State is required"})
			return
		}

		user, err := usersFeature.SpotifyGetUserAuth(state)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		if user == nil {
			c.JSON(403, gin.H{"error": "Session not found"})
			return
		}

		// Retrieve the current user profile
		_, err = user.Client.CurrentUser(c.Request.Context())
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Session is valid"})
	})

	server.Engine.GET("/api/spotify/disconnect", func(c *gin.Context) {
		state := c.Query("state") // Use the state parameter to identify the session

		if state == "" {
			c.JSON(403, gin.H{"error": "State is required"})
			return
		}

		user, err := usersFeature.SpotifyGetUserAuth(state)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		if user == nil {
			c.JSON(403, gin.H{"error": "Session not found"})
			return
		}

		usersFeature.SpotifyDeleteUserAuth(state)

		c.JSON(200, gin.H{"message": "Session disconnected"})
	})
}
