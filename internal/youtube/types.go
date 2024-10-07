package youtube

import (
	youtubeLib "google.golang.org/api/youtube/v3"
)

type SimpleVideoData struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Channel string `json:"channel"`
	Views   uint64 `json:"views"`
	URL     string `json:"url"`

	Snippet    *youtubeLib.VideoSnippet    `json:"snippet"`
	Statistics *youtubeLib.VideoStatistics `json:"statistics"`
}
