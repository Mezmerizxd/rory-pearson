package spotify_manager

import (
	"context"
	"fmt"
	"rory-pearson/environment"
	"rory-pearson/pkg/log"
	"sync"
	"time"

	zSpotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type Config struct {
	Log log.Log
}

// Session represents a user's Spotify session, including the client for making API calls,
// the OAuth2 token for authentication, and the state for session identification.
type Session struct {
	Client *zSpotify.Client
	Token  *oauth2.Token
	State  string
}

// SpotifyManager manages Spotify sessions and handles token lifecycle management.
type SpotifyManager struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	Log         log.Log
	RedirectUrl string
	Auth        *spotifyauth.Authenticator
	Sessions    map[string]*Session // Map of Spotify sessions keyed by the state
}

var instance *SpotifyManager

// Initialize sets up a singleton SpotifyManager instance.
// It initializes the environment, sets up Spotify authentication, and prepares session management.
func Initialize(c Config) *SpotifyManager {
	if instance != nil {
		return instance
	}

	env, err := environment.Initialize()
	if err != nil {
		c.Log.Error().Err(err).Msg("failed to initialize environment")
		return nil
	}

	// Create a context with cancellation for managing session cleanup and other tasks.
	ctx, cancel := context.WithCancel(context.Background())

	// Generate the Spotify redirect URL using the server host environment variable.
	redirectUrl := fmt.Sprintf("%s/api/spotify/callback", env.ServerHost)

	// Create and configure the SpotifyManager with authentication and session management.
	authentication := &SpotifyManager{
		ctx:    ctx,
		cancel: cancel,
		Log:    c.Log,
		Auth: spotifyauth.New(
			spotifyauth.WithRedirectURL(redirectUrl),
			spotifyauth.WithScopes(
				spotifyauth.ScopeUserReadPrivate,
				spotifyauth.ScopeUserReadEmail,
				spotifyauth.ScopePlaylistReadPrivate,
				spotifyauth.ScopeUserReadCurrentlyPlaying,
			),
			spotifyauth.WithClientID(env.SpotifyClientId),         // Replace with actual Client ID
			spotifyauth.WithClientSecret(env.SpotifyClientSecret), // Replace with actual Client Secret
		),
		Sessions: make(map[string]*Session),
	}

	authentication.Log.Info().Msg("SpotifyManager initialized")

	instance = authentication
	return instance
}

// GetInstance returns the singleton instance of SpotifyManager. It panics if the manager has not been initialized.
func GetInstance() *SpotifyManager {
	if instance == nil {
		panic("SpotifyManager not initialized")
	}
	return instance
}

// Close gracefully shuts down the SpotifyManager, cancels any ongoing context tasks, and clears sessions.
func (s *SpotifyManager) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel() // Cancel any ongoing cleanup or long-running operations
	}

	// Clear all stored sessions and log the closure.
	s.Sessions = make(map[string]*Session)
	s.Log.Info().Msg("SpotifyManager closed")
	s.Log.Close()
}

// StoreSession stores a Spotify session identified by the state, including the OAuth token and Spotify client.
func (s *SpotifyManager) StoreSession(ctx context.Context, state string, token *oauth2.Token) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the session with a new Spotify client using the provided token.
	s.Sessions[state] = &Session{
		Client: zSpotify.New(s.Auth.Client(ctx, token)),
		Token:  token,
		State:  state,
	}
}

// GetSession retrieves a session by state and refreshes its token if it has expired.
func (s *SpotifyManager) GetSession(state string) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Retrieve the session from the stored sessions.
	session := s.Sessions[state]
	if session == nil {
		return nil
	}

	if session.Token == nil {
		return session
	}

	// Check if the token is expired. If expired, refresh it.
	if isTokenExpired(session.Token) {
		token, err := s.Auth.RefreshToken(s.ctx, session.Token)
		if err != nil {
			s.Log.Error().Err(err).Msg("failed to refresh token")
			return nil
		}

		// Update the session with the new token.
		session.Token = token
		s.Sessions[state] = session
	}

	return session
}

// DestroySession removes a session by state, effectively logging the user out.
func (s *SpotifyManager) DestroySession(state string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Delete the session from the map.
	delete(s.Sessions, state)
}

// StartSessionCleanup periodically removes expired sessions at the given interval.
// It uses the SpotifyManager's context for canceling the cleanup routine when necessary.
func (s *SpotifyManager) StartSessionCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.mu.Lock()
				// Iterate through sessions and remove expired ones.
				for state, session := range s.Sessions {
					if isTokenExpired(session.Token) {
						delete(s.Sessions, state)
					}
				}
				s.mu.Unlock()
			case <-s.ctx.Done(): // When context is canceled, stop the ticker and exit.
				ticker.Stop()
				return
			}
		}
	}()
}

// isTokenExpired checks whether the OAuth token is expired or nearing expiration.
// If the token is nil or has already expired, the function returns true.
func isTokenExpired(token *oauth2.Token) bool {
	if token == nil {
		return true
	}
	return token.Expiry.Before(time.Now())
}
