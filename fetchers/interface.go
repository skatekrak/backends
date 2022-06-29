package fetchers

type FetchResponse struct {
	Title       string
	Description string
	IconURL     string
	CoverURL    string
}

type SourceFetcher interface {
	Type() string
	IsFromSource(url string) bool
	Fetch(url string) (FetchResponse, error)
	GetSourceID(url string) (string, error)
}
