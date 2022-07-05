package fetchers

import "github.com/skatekrak/scribe/helpers"

func (fe *Fetcher) FetchFeedlySources(categoryID string) ([]ChannelFetchData, error) {
	if !fe.f.HasRefreshToken() {
		if _, err := fe.f.RefreshToken(); err != nil {
			return []ChannelFetchData{}, err
		}
	}

	data, err := fe.f.FetchSources(categoryID)
	if err != nil {
		return []ChannelFetchData{}, err
	}

	feeds := data[0].Feeds
	sources := make([]ChannelFetchData, len(feeds))

	for i, feed := range feeds {
		sources[i] = ChannelFetchData{
			Title:       feed.Title,
			IconURL:     feed.IconURL,
			CoverURL:    feed.VisualURL,
			Description: feed.Description,
			WebsiteURL:  feed.Website,
			SourceID:    feed.FeedID,
			SkateSource: helpers.Has(feed.Topics, "skate"),
			Lang:        feed.Language,
		}
	}

	return sources, nil
}
