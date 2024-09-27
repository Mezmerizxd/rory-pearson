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

type Session struct {
	Client *zSpotify.Client
	Token  *oauth2.Token
	State  string
}

type SpotifyManager struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	Log         log.Log
	RedirectUrl string
	Auth        *spotifyauth.Authenticator
	Sessions    map[string]*Session // Use a map to store sessions by state
}

var instance *SpotifyManager

func Initialize(c Config) *SpotifyManager {
	if instance != nil {
		return instance
	}

	ctx, cancel := context.WithCancel(context.Background())

	redirectUrl := fmt.Sprintf("%s/api/spotify/callback", environment.Get().ServerHost)

	authentication := &SpotifyManager{
		ctx:    ctx,
		cancel: cancel,
		Log:    c.Log,
		Auth: spotifyauth.New(
			spotifyauth.WithRedirectURL(redirectUrl),
			spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
			spotifyauth.WithClientID(environment.Get().SpotifyClientId),         // Update with your actual Client ID
			spotifyauth.WithClientSecret(environment.Get().SpotifyClientSecret), // Update with your actual Client Secret
		),
		Sessions: make(map[string]*Session),
	}

	authentication.Log.Info().Msg("SpotifyManager initialized")

	instance = authentication
	return instance
}

func GetInstance() *SpotifyManager {
	if instance == nil {
		panic("SpotifyManager not initialized")
	}
	return instance
}

func (s *SpotifyManager) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel() // Stop the cleanup goroutine
	}

	s.Sessions = make(map[string]*Session)

	s.Log.Info().Msg("SpotifyManager closed")
	s.Log.Close()
}

func (s *SpotifyManager) StoreSession(ctx context.Context, state string, token *oauth2.Token) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Sessions[state] = &Session{
		Client: zSpotify.New(s.Auth.Client(ctx, token)),
		Token:  token,
		State:  state,
	}
}

func (s *SpotifyManager) GetSession(state string) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Sessions[state]
}

func (s *SpotifyManager) DestroySession(state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Sessions, state)
}

func (s *SpotifyManager) StartSessionCleanup(interval time.Duration) {
	// Use the context and cancel created during initialization
	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.mu.Lock()
				for state, session := range s.Sessions {
					if isTokenExpired(session.Token) {
						delete(s.Sessions, state)
					}
				}
				s.mu.Unlock()
			case <-s.ctx.Done(): // Stop the ticker when context is canceled
				ticker.Stop()
				return
			}
		}
	}()
}

func isTokenExpired(token *oauth2.Token) bool {
	return token.Expiry.Before(time.Now())
}
