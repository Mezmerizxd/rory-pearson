package spotify_controllers

import (
	"net/http"
	"rory-pearson/internal/spotify_manager"
	"rory-pearson/pkg/server"
	"rory-pearson/pkg/util"

	"github.com/gin-gonic/gin"
)

func Initialize(server *server.Server) {
	server.Cfg.Log.Info().Msg("Initializing Spotify controllers")

	sm := spotify_manager.GetInstance()
	if sm == nil {
		server.Cfg.Log.Error().Msg("Spotify manager not initialized")
		return
	}

	server.Engine.GET("/api/spotify/login", func(c *gin.Context) {
		// Generate a new UUID for the state
		state := util.GenerateUUIDv4()

		// Generate the authorization URL with the state
		url := sm.Auth.AuthURL(state)

		// Store the state in the session
		sm.StoreSession(c.Request.Context(), state, nil) // Initial store for state, token will be added later

		c.JSON(http.StatusOK, gin.H{"url": url, "state": state})
	})

	server.Engine.GET("/api/spotify/callback", func(c *gin.Context) {
		state := c.Query("state")

		if state == "" {
			c.JSON(403, gin.H{"error": "State is required"})
			return
		}

		// Verify state matches
		session := sm.GetSession(state)
		if session == nil {
			c.JSON(403, gin.H{"error": "Session not found"})
			return
		}

		token, err := sm.Auth.Token(c.Request.Context(), state, c.Request)
		if err != nil {
			c.JSON(403, gin.H{"error": "Couldn't get token"})
			return
		}

		sm.StoreSession(c.Request.Context(), state, token)

		c.Redirect(http.StatusFound, "/spotify")
	})

	server.Engine.GET("/api/spotify/profile", func(c *gin.Context) {
		state := c.Query("state") // Use the state parameter to identify the session

		if state == "" {
			c.JSON(403, gin.H{"error": "State is required"})
			return
		}

		session := sm.GetSession(state)
		if session == nil {
			c.JSON(403, gin.H{"error": "Session not found"})
			return
		}

		// Retrieve the current user profile
		user, err := session.Client.CurrentUser(c.Request.Context())
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		playlists, err := session.Client.GetPlaylistsForUser(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"user": user, "playlists": playlists})
	})

	server.Engine.GET("/api/spotify/validate", func(c *gin.Context) {
		state := c.Query("state") // Use the state parameter to identify the session

		if state == "" {
			c.JSON(403, gin.H{"error": "State is required"})
			return
		}

		session := sm.GetSession(state)
		if session == nil {
			c.JSON(403, gin.H{"error": "Session not found"})
			return
		}

		// Retrieve the current user profile
		_, err := session.Client.CurrentUser(c.Request.Context())
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Session is valid"})
	})
}
