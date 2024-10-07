package spotify

import (
	"context"
	"fmt"
	"rory-pearson/environment"
	"sync"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type SpotifyPlugin struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	RedirectUrl string
	Auth        *spotifyauth.Authenticator
}

// Initialize sets up the Spotify plugin with necessary configurations.
func (p *SpotifyPlugin) Initialize() (*SpotifyPlugin, error) {
	// Lock to avoid race conditions if Initialize is called concurrently.
	p.mu.Lock()
	defer p.mu.Unlock()

	// Ensure the plugin isn't already initialized.
	if p.ctx != nil {
		return p, nil
	}

	// Retrieve environment variables.
	env := environment.Get()
	if env.ServerHost == "" || env.SpotifyClientId == "" || env.SpotifyClientSecret == "" {
		return nil, fmt.Errorf("missing environment variables for Spotify setup")
	}

	// Create a context with cancellation for managing session cleanup and other tasks.
	ctx, cancel := context.WithCancel(context.Background())

	// Generate the Spotify redirect URL using the server host environment variable.
	redirectUrl := fmt.Sprintf("%s/api/spotify/callback", env.ServerHost)

	// Set up the Spotify authenticator.
	p.ctx = ctx
	p.cancel = cancel
	p.RedirectUrl = redirectUrl
	p.Auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectUrl),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadEmail,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopeUserReadCurrentlyPlaying,
		),
		spotifyauth.WithClientID(env.SpotifyClientId),
		spotifyauth.WithClientSecret(env.SpotifyClientSecret),
	)

	return p, nil
}

// Close cancels the plugin's context, releasing any resources.
func (p *SpotifyPlugin) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cancel != nil {
		p.cancel() // Cancel any ongoing cleanup or long-running operations
		p.cancel = nil
	}
}
