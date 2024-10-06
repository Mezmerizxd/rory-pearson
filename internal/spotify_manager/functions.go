package spotify_manager

import (
	zSpotify "github.com/zmb3/spotify/v2"
)

func (s *SpotifyManager) GetNamesFromPlaylistTracks(session Session, playlistId string) ([]string, error) {
	// Get the playlist tracks
	tracks, err := session.Client.GetPlaylistItems(s.ctx, zSpotify.ID(playlistId))
	if err != nil {
		return nil, err
	}

	names := make([]string, len(tracks.Items))
	for i, track := range tracks.Items {
		names[i] = track.Track.Track.Name
	}

	return names, nil
}
