package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type SnippetThumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type SnippetThumbnails struct {
	Default  SnippetThumbnail  `json:"default"`
	Medium   *SnippetThumbnail `json:"medium"`
	High     *SnippetThumbnail `json:"high"`
	Standard *SnippetThumbnail `json:"standard"`
}

type ChannelItemSnippet struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	PublishedAt time.Time         `json:"publishedAt"`
	Country     string            `json:"country"`
	Thumbnails  SnippetThumbnails `json:"thumbnails"`
}

type ChannelBrandingSettingsImage struct {
	BannerExternalURL string `json:"bannerExternalUrl"`
}

type ChannelBrandingSettings struct {
	Image ChannelBrandingSettingsImage `json:"image"`
}

type ChannelItem struct {
	Kind             string                  `json:"kind"`
	Etag             string                  `json:"etag"`
	ID               string                  `json:"id"`
	Snippet          ChannelItemSnippet      `json:"snippet"`
	BrandingSettings ChannelBrandingSettings `json:"brandingSettings"`
}

type FetchResponse[T any] struct {
	Kind     string   `json:"kind"`
	Etag     string   `json:"etag"`
	PageInfo PageInfo `json:"pageInfo"`
	Items    []T      `json:"items"`
}

func FetchChannel(channelID, accessToken string) (FetchResponse[ChannelItem], error) {
	response, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=snippet,brandingSettings&id=%s&key=%s", channelID, accessToken))

	if err != nil {
		return FetchResponse[ChannelItem]{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return FetchResponse[ChannelItem]{}, err
	}

	var data FetchResponse[ChannelItem]
	if err := json.Unmarshal(responseData, &data); err != nil {
		return FetchResponse[ChannelItem]{}, err
	}

	return data, nil
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

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

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
