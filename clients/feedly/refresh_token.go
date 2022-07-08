package feedly

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type FeedlyRefreshTokenResopnse struct {
	Provider     string `json:"provider"`
	IsNewAccount bool   `json:"is_new_account"`
	Plan         string `json:"plan"`
	ID           string `json:"id"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func (f *FeedlyClient) RefreshToken() (string, error) {
	url := fmt.Sprintf("https://cloud.feedly.com/v3/auth/token?refresh_token=%s&client_id=feedlydev&client_secret=feedlydev&grant_type=refresh_token", f.refreshToken)
	req, err := http.Post(url, "application/json", nil) //#nosec G107 -- False positive
	if err != nil {
		return "", err
	}

	responseData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	var data FeedlyRefreshTokenResopnse
	if err := json.Unmarshal(responseData, &data); err != nil {
		return "", err
	}

	if data.AccessToken == "" {
		return "", errors.New("empty access token")
	}

	f.accessToken = data.AccessToken
	log.Println("Feedly token refreshed")

	return data.AccessToken, nil
}

func (f *FeedlyClient) HasRefreshToken() bool {
	return f.accessToken != ""
}
