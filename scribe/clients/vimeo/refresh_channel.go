package vimeo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type VimeoVideoItem struct {
	URI            string        `json:"uri"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Type           string        `json:"type"`
	Link           string        `json:"link"`
	PlayerEmbedURL string        `json:"player_embed_url"`
	ReleaseTime    time.Time     `json:"release_time"`
	Pictures       VimeoPictures `json:"pictures"`
}

type VimeoPaging struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	First    string `json:"first"`
	Last     string `json:"last"`
}

type VimeoVideosResponse struct {
	Total   int              `json:"total"`
	Page    int              `json:"page"`
	PerPage int              `json:"per_page"`
	Paging  VimeoPaging      `json:"paging"`
	Data    []VimeoVideoItem `json:"data"`
}

func (v *VimeoClient) FetchVideos(userID string) (VimeoVideosResponse, error) {
	url := fmt.Sprintf("https://api.vimeo.com/users/%s/videos", userID)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return VimeoVideosResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", v.apiKey))
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return VimeoVideosResponse{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return VimeoVideosResponse{}, err
	}

	var data VimeoVideosResponse
	if err := json.Unmarshal(responseData, &data); err != nil {
		return VimeoVideosResponse{}, err
	}

	return data, err
}
