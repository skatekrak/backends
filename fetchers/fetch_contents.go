package fetchers

import (
	"errors"
)

func (fe *Fetcher) FetchChannelContents(sourceID string, sourceType string) ([]ContentFetchData, error) {
	switch sourceType {
	case "youtube":
		return fe.FetchYoutubeChannelContents(sourceID)
	case "vimeo":
		return fe.FetchVimeoChannelContents(sourceID)
	default:
		return []ContentFetchData{}, errors.New("sourceType not supported")
	}
}
