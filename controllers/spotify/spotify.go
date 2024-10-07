package spotify

import (
	"rory-pearson/internal/spotify"
	"rory-pearson/pkg/server"

	"github.com/gin-gonic/gin"
	zSpotify "github.com/zmb3/spotify/v2"
)

func Initialize(server *server.Server) {
	server.Cfg.Log.Info().Msg("Initializing Spotify controllers")

	AuthRoutes(server)

	sm := spotify.GetInstance()
	if sm == nil {
		server.Cfg.Log.Error().Msg("Spotify manager not initialized")
		return
	}

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

		c.JSON(200, user)
	})

	server.Engine.GET("/api/spotify/playlists", func(c *gin.Context) {
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

		// Retrieve the current user playlists
		playlists, err := session.Client.CurrentUsersPlaylists(c.Request.Context())
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, playlists)
	})

	server.Engine.GET("/api/spotify/tracks", func(c *gin.Context) {
		state := c.Query("state") // Use the state parameter to identify the session
		playlistId := c.Query("playlistId")

		if state == "" {
			c.JSON(403, gin.H{"error": "State is required"})
			return
		}

		if playlistId == "" {
			c.JSON(403, gin.H{"error": "Playlist ID is required"})
			return
		}

		session := sm.GetSession(state)
		if session == nil {
			c.JSON(403, gin.H{"error": "Session not found"})
			return
		}

		tracks, err := session.Client.GetPlaylistItems(c.Request.Context(), zSpotify.ID(playlistId))
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, tracks)
	})

	server.Engine.GET("/api/spotify/now-playing", func(c *gin.Context) {
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
		np, err := session.Client.PlayerCurrentlyPlaying(c.Request.Context())
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, np)
	})
}
