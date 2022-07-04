package vimeo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ChannelPictureSize struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Link   string `json:"link"`
}

type VimeoPictures struct {
	URI      string               `json:"uri"`
	Active   bool                 `json:"active"`
	BaseLink string               `json:"base_link"`
	Sizes    []ChannelPictureSize `json:"sizes"`
}

type FetchChannelResponse struct {
	URI         string        `json:"uri"`
	Name        string        `json:"name"`
	Link        string        `json:"link"`
	Bio         string        `json:"bio"`
	ShortBio    string        `json:"short_bio"`
	CreatedTime time.Time     `json:"created_time"`
	Pictures    VimeoPictures `json:"pictures"`
}

func (v *VimeoClient) FetchChannel(userID string) (FetchChannelResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.vimeo.com/users/%s", userID), nil)
	if err != nil {
		return FetchChannelResponse{}, err
	}

	defer req.Body.Close()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", v.apiKey))
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return FetchChannelResponse{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return FetchChannelResponse{}, err
	}

	var data FetchChannelResponse
	if err := json.Unmarshal(responseData, &data); err != nil {
		return FetchChannelResponse{}, err
	}

	return data, err
}
