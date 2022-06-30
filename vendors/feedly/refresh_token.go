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

func RefreshToken(refreshToken string) (string, error) {
	req, err := http.Get(fmt.Sprintf("https://cloud.feedly.com/v3/auth/token?refresh_token=%s&client_id=%s&client_secret=%s&grant_type=%s", refreshToken, "feedlydev", "feedlydev", "refresh_token"))
	if err != nil {
		return "", err
	}

	defer req.Body.Close()
	responseData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	var data FeedlyRefreshTokenResopnse
	if err := json.Unmarshal(responseData, &data); err != nil {
		return "", err
	}

	return data.AccessToken, nil
}
