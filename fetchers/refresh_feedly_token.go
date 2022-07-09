package fetchers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

func (fe *Fetcher) RefreshFeedlyToken() (string, time.Time, error) {
	url := fmt.Sprintf("https://cloud.feedly.com/v3/auth/token?refresh_token=%s&client_id=feedlydev&client_secret=feedlydev&grant_type=refresh_token", fe.f.RefreshToken)
	req, err := http.Post(url, "application/json", nil) //#nosec G107 -- False positive
	if err != nil {
		return "", time.Time{}, err
	}

	responseData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", time.Time{}, err
	}

	var data FeedlyRefreshTokenResopnse
	if err := json.Unmarshal(responseData, &data); err != nil {
		return "", time.Time{}, err
	}

	if data.AccessToken == "" {
		return "", time.Time{}, errors.New("empty access token")
	}

	expiresAt := time.Now()
	expiresAt = expiresAt.Add(time.Second * time.Duration(data.ExpiresIn))

	log.Printf("Feedly token refreshed, expires at %s", expiresAt)

	return data.AccessToken, expiresAt, err
}
