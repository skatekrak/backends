package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

func (y *YoutubeClient) FetchChannel(channelID string) (FetchResponse[ChannelItem], error) {
	response, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=snippet,brandingSettings&id=%s&key=%s", channelID, y.apiKey))

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
