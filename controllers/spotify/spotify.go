package spotify

import (
	"rory-pearson/internal/users"
	"rory-pearson/pkg/features"
	"rory-pearson/pkg/server"

	"github.com/gin-gonic/gin"
	zSpotify "github.com/zmb3/spotify/v2"
)

func Initialize(server *server.Server) {
	server.Cfg.Log.Info().Msg("Initializing Spotify controllers")

	AuthRoutes(server)

	f, _ := features.GetInstance().GetFeature(users.UsersFeatureType)
	usersFeature, _ := f.(*users.UsersFeature)

	server.Engine.GET("/api/spotify/profile", func(c *gin.Context) {
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
		profile, err := user.Client.CurrentUser(c.Request.Context())
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, profile)
	})

	server.Engine.GET("/api/spotify/playlists", func(c *gin.Context) {
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

		// Retrieve the current user playlists
		playlists, err := user.Client.CurrentUsersPlaylists(c.Request.Context())
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

		user, err := usersFeature.SpotifyGetUserAuth(state)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		if user == nil {
			c.JSON(403, gin.H{"error": "Session not found"})
			return
		}

		tracks, err := user.Client.GetPlaylistItems(c.Request.Context(), zSpotify.ID(playlistId))
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

		// session := sm.GetSession(state)
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
		np, err := user.Client.PlayerCurrentlyPlaying(c.Request.Context())
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, np)
	})

	// type PlaylistToYoutubeRequest struct {
	// 	PlaylistId string `json:"playlistId"`
	// }
	// server.Engine.POST("/api/spotify/playlist-to-youtube", func(c *gin.Context) {
	// 	state := c.Query("state") // Use the state parameter to identify the session

	// 	var req PlaylistToYoutubeRequest
	// 	if err := c.ShouldBindJSON(&req); err != nil {
	// 		c.JSON(403, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	if state == "" {
	// 		c.JSON(403, gin.H{"error": "State is required"})
	// 		return
	// 	}

	// 	if req.PlaylistId == "" {
	// 		c.JSON(403, gin.H{"error": "Playlist ID is required"})
	// 		return
	// 	}

	// 	session := sm.GetSession(state)
	// 	if session == nil {
	// 		c.JSON(403, gin.H{"error": "Session not found"})
	// 		return
	// 	}

	// 	playlistNames, err := sm.GetNamesFromPlaylistTracks(*session, req.PlaylistId)
	// 	if err != nil {
	// 		c.JSON(403, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	// Convert the playlist to a YouTube playlist
	// 	youtubeData, err := yt.BulkSearchWithHighestViews(playlistNames)
	// 	if err != nil {
	// 		c.JSON(403, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	c.JSON(200, gin.H{"youtubeData": youtubeData})
	// })
}
