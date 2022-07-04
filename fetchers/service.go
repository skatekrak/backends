package fetchers

import (
	"time"

	"github.com/skatekrak/scribe/clients/feedly"
	"github.com/skatekrak/scribe/clients/vimeo"
	"github.com/skatekrak/scribe/clients/youtube"
)

// Abstract representation of a Source from a channel
type ChannelFetchData struct {
	Title       string
	Description string
	IconURL     string
	CoverURL    string
	PublishedAt time.Time
}

type ContentFetchData struct {
	Title          string
	PublishedAt    time.Time
	Description    string // Or Summary
	RawDescription string // or RawSummary
	ThumbnailURL   string
	ContentID      string // or VideoID
	ContentURL     string
}

type Fetcher struct {
	v *vimeo.VimeoClient
	y *youtube.YoutubeClient
	f *feedly.FeedlyClient
}

func New(v *vimeo.VimeoClient, y *youtube.YoutubeClient, f *feedly.FeedlyClient) *Fetcher {
	return &Fetcher{
		v: v,
		y: y,
		f: f,
	}
}

func (fe *Fetcher) SourceType(url string) string {
	if youtube.IsYoutubeChannel(url) {
		return "youtube"
	} else if vimeo.IsVimeoUser(url) {
		return "vimeo"
	}
	return "rss"
}
