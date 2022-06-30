package vimeo

import (
	"errors"

	"github.com/skatekrak/scribe/fetchers"
)

func (v VimeoSourceFetcher) RefreshSource(sourceID string) ([]fetchers.RefreshResponse, error) {
	return []fetchers.RefreshResponse{}, errors.New("not implemented")
}
