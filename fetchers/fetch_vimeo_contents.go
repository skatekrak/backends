package fetchers

import (
	"fmt"

	"github.com/skatekrak/scribe/vendors/clients/vimeo"
)

func (fe *Fetcher) fetchVimeoChannelContent(userID string) ([]ContentFetchData, error) {
	data, err := fe.v.FetchVideos(userID)
	if err != nil {
		return []ContentFetchData{}, err
	}

	items := make([]ContentFetchData, len(data.Data))

	for i, item := range data.Data {
		items[i] = ContentFetchData{
			Title:          item.Name,
			Description:    item.Description,
			PublishedAt:    item.ReleaseTime,
			RawDescription: item.Description,
			ThumbnailURL:   vimeo.GetLargerImageLink(item.Pictures.Sizes),
			ContentID:      item.URI,
			ContentURL:     fmt.Sprintf("https://vimeo.com%s", item.URI),
		}
	}

	return items, nil
}

func (fe *Fetcher) FetcherVimeoContent(userIDs []string, contents []ContentFetchData) map[string]error {
	errors := make(map[string]error)

	for _, userID := range userIDs {
		c, err := fe.fetchVimeoChannelContent(userID)
		if err != nil {
			errors[userID] = err
		} else {
			contents = append(contents, c...)
		}
	}

	return errors
}
