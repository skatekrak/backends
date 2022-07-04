package vimeo

import (
	"errors"
	"regexp"
	"sort"
)

type VimeoClient struct {
	apiKey string
}

func New(apiKey string) *VimeoClient {
	return &VimeoClient{apiKey}
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

func GetLargerImageLink(sizes []ChannelPictureSize) string {
	if len(sizes) <= 0 {
		return ""
	}

	sort.SliceStable(sizes, func(i, j int) bool {
		return sizes[i].Height >= sizes[j].Height
	})

	return sizes[0].Link
}
