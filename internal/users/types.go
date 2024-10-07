package users

import (
	zSpotify "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

type User struct {
	ID          string
	SpotifyAuth *SpotifyAuth
}

type SpotifyAuth struct {
	Client *zSpotify.Client
	Token  *oauth2.Token
	State  string
}
