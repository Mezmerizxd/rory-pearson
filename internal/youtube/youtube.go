package youtube

import (
	"context"
	"rory-pearson/environment"
	"sync"

	"google.golang.org/api/option"
	youtubeLib "google.golang.org/api/youtube/v3"
)

type Config struct{}

type Youtube struct {
	Ctx     context.Context
	Service youtubeLib.Service

	mu sync.Mutex
}

var instance *Youtube

func Initialize(c Config) (*Youtube, error) {
	if instance != nil {
		return instance, nil
	}

	ctx, _ := context.WithCancel(context.Background())

	service, err := youtubeLib.NewService(ctx, option.WithAPIKey(environment.Get().YoutubeApiKey))
	if err != nil {
		return nil, err
	}

	youtube := &Youtube{
		Ctx:     ctx,
		Service: *service,
	}

	instance = youtube

	return youtube, nil
}

func GetInstance() *Youtube {
	if instance == nil {
		panic("Youtube not initialized")
	}
	return instance
}

func (y *Youtube) SimpleSearch(query string, scope int64) (*[]SimpleVideoData, error) {
	var videos []SimpleVideoData

	searchResponse, err := y.Service.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(scope).Do()
	if err != nil {
		return nil, err
	}

	var videoIds []string
	for _, item := range searchResponse.Items {
		if item.Id.Kind == "youtube#video" {
			videoIds = append(videoIds, item.Id.VideoId)
		}
	}

	if len(videoIds) > 0 {
		videoStatsResponse, err := y.Service.Videos.List([]string{"snippet", "statistics"}).Id(videoIds...).Do()
		if err != nil {
			return nil, err
		}

		for _, item := range videoStatsResponse.Items {
			videos = append(videos, SimpleVideoData{
				ID:         item.Id,
				Title:      item.Snippet.Title,
				Channel:    item.Snippet.ChannelTitle,
				Views:      item.Statistics.ViewCount,
				URL:        "https://www.youtube.com/watch?v=" + item.Id,
				Snippet:    item.Snippet,
				Statistics: item.Statistics,
			})
		}
	}

	return &videos, nil
}

func (y *Youtube) SearchVideoWithHighestViews(query string) (*SimpleVideoData, error) {
	videos, err := y.SimpleSearch(query, 3)
	if err != nil {
		return nil, err
	}

	if len(*videos) == 0 {
		return nil, nil
	}

	var highestViewedVideo SimpleVideoData

	for _, video := range *videos {
		if video.Views > highestViewedVideo.Views {
			highestViewedVideo = video
		}
	}

	return &highestViewedVideo, nil
}

func (y *Youtube) BulkSearchWithHighestViews(query []string) (*[]SimpleVideoData, error) {
	var videos []SimpleVideoData

	for _, q := range query {
		video, err := y.SearchVideoWithHighestViews(q)
		if err != nil {
			return nil, err
		}

		if video != nil {
			videos = append(videos, *video)
		}
	}

	return &videos, nil
}
