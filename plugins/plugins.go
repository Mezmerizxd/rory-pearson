package plugins

import (
	"rory-pearson/plugins/spotify"
)

type Config struct{}

type Plugins struct {
	Spotify *spotify.SpotifyPlugin
}

var instance *Plugins

func Initialize() (*Plugins, error) {
	if instance != nil {
		return instance, nil
	}

	instance = &Plugins{
		Spotify: &spotify.SpotifyPlugin{},
	}

	// Initialize plugins here
	_, err := instance.Spotify.Initialize()
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (p *Plugins) Close() {
	p.Spotify.Close()
}

func GetInstance() Plugins {
	if instance == nil {
		panic("Plugins not initialized")
	}
	return *instance
}
