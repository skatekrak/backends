package fetchers

import "time"

type FetchResponse struct {
	Title       string
	Description string
	IconURL     string
	CoverURL    string
	PublishedAt time.Time
}

type RefreshResponse struct {
	Title          string
	PublishedAt    time.Time
	Description    string // Or Summary
	RawDescription string // or RawSummary
	ThumbnailURL   string
	ContentID      string // or VideoID
	ContentURL     string
}

type SourceFetcher interface {
	Type() string
	ContentType() string
	IsFromSource(url string) bool
	Fetch(url string) (FetchResponse, error)
	GetSourceID(url string) (string, error)
	RefreshSource(sourceID string) ([]RefreshResponse, error)
}
