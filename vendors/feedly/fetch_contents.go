package feedly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FeedlyItemOrigin struct {
	StreamID string `json:"streamId"`
	Title    string `json:"title"`
	HtmlURL  string `json:"htmlUrl"`
}

type FeedlyItemContent struct {
	Content   string `json:"content"`
	Direction string `json:"direction"`
}

type FeedlyItemAlternate struct {
	Type string `json:"type"`
	Href string `json:"href"`
}

type FeedlyItemVisual struct {
	Height         int    `json:"height"`
	Processor      string `json:"string"`
	ExpirationDate int    `json:"expirationDate"`
	EdgeCacheURL   string `json:"edgeCacheUrl"`
	ContentType    string `json:"contentType"`
	Width          int    `json:"width"`
	URL            string `json:"url"`
}

type FeedlyFetchContentItem struct {
	Language     string                `json:"language"`
	Fingerprint  string                `json:"fingerprint"`
	ID           string                `json:"id"`
	Keywords     []string              `json:"keywords"`
	OriginID     string                `json:"originId"`
	Origin       FeedlyItemOrigin      `json:"origin"`
	Title        string                `json:"title"`
	Author       string                `json:"author"`
	Crawled      int                   `json:"crawled"`
	Content      FeedlyItemContent     `json:"content"`
	Published    int                   `json:"published"`
	Summary      FeedlyItemContent     `json:"summary"`
	CanonicalURL string                `json:"canonicalUrl"`
	Alternate    []FeedlyItemAlternate `json:"alternate"`
	Visual       FeedlyItemVisual      `json:"visual"`
}

type FeedlyFetchContentsResponse struct {
	ID           string                   `json:"id"`
	Updated      int                      `json:"updated"`
	Continuation string                   `json:"string"`
	Items        []FeedlyFetchContentItem `json:"items"`
}

func FetchContents(token string, categoryID string) (FeedlyFetchContentsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://cloud.feedly.com/v3/streams/contents?streamId=%s&count=1000", categoryID), nil)
	if err != nil {
		return FeedlyFetchContentsResponse{}, err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return FeedlyFetchContentsResponse{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return FeedlyFetchContentsResponse{}, err
	}

	var data FeedlyFetchContentsResponse
	if err := json.Unmarshal(responseData, &data); err != nil {
		return FeedlyFetchContentsResponse{}, err
	}

	return data, err
}
