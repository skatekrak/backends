package vimeo

import (
	"sort"

	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/vendors/vimeo"
)

type VimeoSourceFetcher struct {
	apiKey string
}

func NewVimeoSourceFetcher(apiKey string) VimeoSourceFetcher {
	return VimeoSourceFetcher{apiKey}
}

func (v VimeoSourceFetcher) Type() string {
	return "vimeo"
}

func (v VimeoSourceFetcher) IsFromSource(url string) bool {
	return vimeo.IsVimeoUser(url)
}

func (v VimeoSourceFetcher) GetSourceID(url string) (string, error) {
	return vimeo.GetUserID(url)
}

func (v VimeoSourceFetcher) Fetch(url string) (fetchers.FetchResponse, error) {
	userID, err := vimeo.GetUserID(url)
	if err != nil {
		return fetchers.FetchResponse{}, err
	}

	data, err := vimeo.FetchChannel(userID, v.apiKey)
	if err != nil {
		return fetchers.FetchResponse{}, err
	}

	coverURL := getLargerImageLink(data.Pictures.Sizes)

	return fetchers.FetchResponse{
		Title:       data.Name,
		Description: data.Bio,
		PublishedAt: data.CreatedTime,
		IconURL:     coverURL,
		CoverURL:    coverURL,
	}, nil
}

func getLargerImageLink(sizes []vimeo.ChannelPictureSize) string {
	if len(sizes) <= 0 {
		return ""
	}

	sort.SliceStable(sizes, func(i, j int) bool {
		return sizes[i].Height >= sizes[j].Height
	})

	return sizes[0].Link
}
