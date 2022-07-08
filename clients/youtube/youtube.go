package youtube

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type YoutubeClient struct {
	apiKey string
}

func New(apiKey string) *YoutubeClient {
	return &YoutubeClient{apiKey}
}

func IsYoutubeChannel(url string) bool {
	// Check if URL is of a youtube channel
	match, _ := regexp.MatchString(`(?:https?:\/\/)?(?:(?:(?:www\.?)?youtube\.com\/c\/\w+))`, url)
	return match
}

// For a given youtube channel URL
func GetChannelID(url string) (string, error) {
	if !IsYoutubeChannel(url) {
		return "", errors.New("URL is not a youtube channel")
	}

	res, err := http.Get(url) //#nosec G107 -- Ok because it's only admin and urls are check for domain before
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", errors.New("error fetching url")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	s := doc.Find(`[itemprop="channelId"]`)
	if s.Length() <= 0 {
		return "", errors.New("channelID not found")
	}

	if channelId, ok := s.First().Attr("content"); ok {
		return channelId, nil
	}

	return "", errors.New("channelID not found")
}

func GetBestThumbnail(thumbnails SnippetThumbnails) string {
	if thumbnails.Standard != nil {
		return thumbnails.Standard.URL
	} else if thumbnails.High != nil {
		return thumbnails.High.URL
	} else if thumbnails.Medium != nil {
		return thumbnails.Medium.URL
	}
	return thumbnails.Default.URL
}
