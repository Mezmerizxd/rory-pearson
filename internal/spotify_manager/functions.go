package spotify_manager

import (
	"fmt"

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
		// Format: Title - Artist, Artist, Artist...

		var artists string
		for i, artist := range track.Track.Track.Artists {
			if i == 0 {
				artists = artist.Name
			} else {
				artists += ", " + artist.Name
			}
		}

		names[i] = fmt.Sprintf("%s - %s", track.Track.Track.Name, artists)
	}

	return names, nil
}
