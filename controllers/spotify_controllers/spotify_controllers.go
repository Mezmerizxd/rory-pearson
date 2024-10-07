package spotify_controllers

import (
	"net/http"
	"rory-pearson/internal/spotify_manager"
	"rory-pearson/internal/youtube"
	"rory-pearson/pkg/server"
	"rory-pearson/pkg/util"

	"github.com/gin-gonic/gin"
	zSpotify "github.com/zmb3/spotify/v2"
)

func Initialize(server *server.Server) {
	server.Cfg.Log.Info().Msg("Initializing Spotify controllers")

	sm := spotify_manager.GetInstance()
	yt := youtube.GetInstance()

	server.Engine.GET("/api/spotify/login", func(c *gin.Context) {
		// Generate a new UUID for the state
		state := util.GenerateUUIDv4()

		// Generate the authorization URL with the state
		url := sm.Auth.AuthURL(state)

		// Store the state in the session
		sm.StoreSession(c.Request.Context(), state, nil) // Initial store for state, token will be added later

		c.JSON(http.StatusOK, gin.H{"url": url, "state": state})
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

	server.Engine.GET("/api/spotify/disconnect", func(c *gin.Context) {
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

		sm.DestroySession(state)

		c.JSON(200, gin.H{"message": "Session disconnected"})
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

	type PlaylistToYoutubeRequest struct {
		PlaylistId string `json:"playlistId"`
	}
	server.Engine.POST("/api/spotify/playlist-to-youtube", func(c *gin.Context) {
		state := c.Query("state") // Use the state parameter to identify the session

		var req PlaylistToYoutubeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		if state == "" {
			c.JSON(403, gin.H{"error": "State is required"})
			return
		}

		if req.PlaylistId == "" {
			c.JSON(403, gin.H{"error": "Playlist ID is required"})
			return
		}

		session := sm.GetSession(state)
		if session == nil {
			c.JSON(403, gin.H{"error": "Session not found"})
			return
		}

		playlistNames, err := sm.GetNamesFromPlaylistTracks(*session, req.PlaylistId)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		// Convert the playlist to a YouTube playlist
		youtubeData, err := yt.BulkSearchWithHighestViews(playlistNames)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"youtubeData": youtubeData})
	})
}
