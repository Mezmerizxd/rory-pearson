package users

import (
	"context"
	"time"

	"golang.org/x/oauth2"
)

func (f *UsersFeature) SpotifyStoreUserAuth(ctx context.Context, state string, token *oauth2.Token) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Store the session with a new Spotify client using the provided token.
	f.Users[state] = &User{
		ID: state,
		SpotifyAuth: &SpotifyAuth{
			Token: token,
			State: state,
		},
	}

	return nil
}

func (f *UsersFeature) SpotifyGetUserAuth(state string) (*SpotifyAuth, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Retrieve the session from the stored sessions.
	user := f.Users[state]
	if user == nil {
		return nil, nil
	}

	if user.SpotifyAuth.Token == nil {
		return user.SpotifyAuth, nil
	}

	// Check if the token is expired. If expired, refresh it.
	if user.SpotifyAuth.Token.Expiry.Before(time.Now()) {
		// TODO: Refresh token
		return nil, nil
	}

	return user.SpotifyAuth, nil
}

func (f *UsersFeature) SpotifyDeleteUserAuth(state string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Delete the session from the map.
	delete(f.Users, state)

	return nil
}
