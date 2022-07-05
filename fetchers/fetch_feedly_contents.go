package fetchers

import (
	"time"

	"github.com/k3a/html2text"
)

func (fe *Fetcher) FetchFeedlyContents(categoryID string) ([]ContentFetchData, error) {
	if !fe.f.HasRefreshToken() {
		if _, err := fe.f.RefreshToken(); err != nil {
			return []ContentFetchData{}, err
		}
	}

	data, err := fe.f.FetchContents(categoryID)
	if err != nil {
		return []ContentFetchData{}, err
	}

	items := make([]ContentFetchData, len(data.Items))

	for i, item := range data.Items {
		var url string

		if len(item.Alternate) > 0 && item.Alternate[0].Href != "" {
			url = item.Alternate[0].Href
		} else {
			url = item.CanonicalURL
		}

		items[i] = ContentFetchData{
			Title:          item.Title,
			Description:    html2text.HTML2Text(item.Summary.Content),
			RawDescription: item.Summary.Content,
			PublishedAt:    time.UnixMilli(int64(item.Published)),
			ThumbnailURL:   item.Visual.URL,
			ContentID:      item.ID,
			ContentURL:     url,
			SourceID:       item.Origin.StreamID,
		}
	}

	return items, nil
}
