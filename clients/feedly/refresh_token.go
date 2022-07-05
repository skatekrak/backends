package feedly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FeedlyRefreshTokenResopnse struct {
	Provider     string `json:"providre"`
	IsNewAccount bool   `json:"is_new_account"`
	Plan         string `json:"plan"`
	ID           string `json:"id"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func (f *FeedlyClient) RefreshToken() (string, error) {
	req, err := http.Get(fmt.Sprintf("https://cloud.feedly.com/v3/auth/token?refresh_token=%s&client_id=%s&client_secret=%s&grant_type=%s", f.refreshToken, "feedlydev", "feedlydev", "refresh_token"))
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

	f.accessToken = data.AccessToken

	return data.AccessToken, nil
}

func (f *FeedlyClient) HasRefreshToken() bool {
	return f.accessToken != ""
}
