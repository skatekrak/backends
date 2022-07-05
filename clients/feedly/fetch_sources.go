package feedly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type FeedlyCollectionFeed struct {
	ID          string   `json:"id"`
	VisualURL   string   `json:"visualUrl"`
	Partial     bool     `json:"partial"`
	ContentType string   `json:"contentType"`
	Language    string   `json:"language"`
	FeedID      string   `json:"feedId"` // Stream ID in Contents
	Title       string   `json:"title"`
	Topics      []string `json:"topics"`
	Updated     int      `json:"updated"`
	Website     string   `json:"website"`
	Subscribers int      `json:"subscribers"`
	IconURL     string   `json:"iconUrl"`
}

type FeedlyCollectionResponse struct {
	Customizable bool   `json:"customizable"`
	Label        string `json:"label"`
	Created      int    `json:"created"`
	Enterprise   bool   `json:"enterprise"`
	NumFeeds     int    `json:"numFeeds"`
	ID           string `json:"id"`
}

func (f *FeedlyClient) FetchSources(categoryID string) ([]FeedlyCollectionResponse, error) {
	url := fmt.Sprintf("https://cloud.feedly.com/v3/collections/%s", url.QueryEscape(categoryID))
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return []FeedlyCollectionResponse{}, err
	}

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return []FeedlyCollectionResponse{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []FeedlyCollectionResponse{}, err
	}

	var data []FeedlyCollectionResponse
	if err := json.Unmarshal(responseData, &data); err != nil {
		return []FeedlyCollectionResponse{}, err
	}

	return data, err
}
