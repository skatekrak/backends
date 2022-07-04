package fetchers

import (
	"errors"

	"github.com/skatekrak/scribe/clients/vimeo"
	"github.com/skatekrak/scribe/clients/youtube"
)

func (fe *Fetcher) GetSourceID(url string) (string, error) {
	sourceType := fe.SourceType(url)
	if sourceType == "rss" {
		return "", errors.New("rss not supported")
	} else if sourceType == "youtube" {
		return youtube.GetChannelID(url)
	} else if sourceType == "vimeo" {
		return vimeo.GetUserID(url)
	}

	return "", errors.New("url invalid or not supported")
}
