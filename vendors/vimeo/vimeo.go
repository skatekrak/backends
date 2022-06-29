package vimeo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type ChannelPictureSize struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Link   string `json:"link"`
}

type ChannelPictures struct {
	URI      string               `json:"uri"`
	Active   bool                 `json:"active"`
	BaseLink string               `json:"base_link"`
	Sizes    []ChannelPictureSize `json:"sizes"`
}

type FetchChannelResponse struct {
	URI         string          `json:"uri"`
	Name        string          `json:"name"`
	Link        string          `json:"link"`
	Bio         string          `json:"bio"`
	ShortBio    string          `json:"short_bio"`
	CreatedTime time.Time       `json:"created_time"`
	Pictures    ChannelPictures `json:"pictures"`
}

func FetchChannel(userID, apiKey string) (FetchChannelResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.vimeo.com/users/%s", userID), nil)
	if err != nil {
		return FetchChannelResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
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

	log.Println("data from vimeo", data)

	return data, err
}

func IsVimeoUser(url string) bool {
	match, _ := regexp.MatchString(`(?:https:\/\/)?vimeo.com\/user\d+`, url)
	return match
}

func GetUserID(url string) (string, error) {
	r, _ := regexp.Compile(`user\d+`)

	userID := r.FindString(url)

	if userID == "" {
		return "", errors.New("userID not found")
	}
	return userID, nil
}
