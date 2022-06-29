package youtube

import (
	"errors"

	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/vendors/youtube"
)

type YoutubeSourceFetcher struct {
	apiKey string
}

func NewYoutubeFetcher(apiKey string) YoutubeSourceFetcher {
	return YoutubeSourceFetcher{apiKey}
}

func (y YoutubeSourceFetcher) Type() string {
	return "youtube"
}

func (y YoutubeSourceFetcher) IsFromSource(url string) bool {
	// Check if URL is of a youtube channel
	return youtube.IsYoutubeChannel(url)
}

func (y YoutubeSourceFetcher) GetSourceID(url string) (string, error) {
	return youtube.GetChannelID(url)
}

func (y YoutubeSourceFetcher) Fetch(url string) (fetchers.FetchResponse, error) {
	channelId, err := youtube.GetChannelID(url)
	if err != nil {
		return fetchers.FetchResponse{}, err
	}

	channelData, err := youtube.FetchChannel(channelId, y.apiKey)
	if err != nil {
		return fetchers.FetchResponse{}, err
	}

	if len(channelData.Items) > 0 {
		item := channelData.Items[0]

		return fetchers.FetchResponse{
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			IconURL:     item.Snippet.Thumbnails.Default.URL,
			CoverURL:    item.BrandingSettings.Image.BannerExternalURL,
			PublishedAt: item.Snippet.PublishedAt,
		}, nil
	}
	return fetchers.FetchResponse{}, errors.New("weird, no items")
}
