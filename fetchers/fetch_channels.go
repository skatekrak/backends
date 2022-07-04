package fetchers

import (
	"errors"

	"github.com/skatekrak/scribe/vendors/clients/vimeo"
	"github.com/skatekrak/scribe/vendors/clients/youtube"
)

// For a given channel URL, will fetch the vimeo or youtube channel data
func (fe *Fetcher) FetchChannelData(url string) (ChannelFetchData, error) {
	sourceType := fe.SourceType(url)
	if sourceType == "rss" {
		return ChannelFetchData{}, errors.New("can't get data from rss")
	} else if sourceType == "youtube" {
		channelId, err := youtube.GetChannelID(url)
		if err != nil {
			return ChannelFetchData{}, err
		}
		data, err := fe.y.FetchChannel(channelId)
		if err != nil {
			return ChannelFetchData{}, err
		}

		if len(data.Items) <= 0 {
			return ChannelFetchData{}, errors.New("channel not found")
		}

		channel := data.Items[0]

		return ChannelFetchData{
			Title:       channel.Snippet.Title,
			Description: channel.Snippet.Description,
			IconURL:     youtube.GetBestThumbnail(channel.Snippet.Thumbnails),
			CoverURL:    channel.BrandingSettings.Image.BannerExternalURL,
			PublishedAt: channel.Snippet.PublishedAt,
		}, nil
	} else if sourceType == "vimeo" {
		userID, err := vimeo.GetUserID(url)
		if err != nil {
			return ChannelFetchData{}, err
		}

		data, err := fe.v.FetchChannel(userID)
		if err != nil {
			return ChannelFetchData{}, err
		}

		coverURL := vimeo.GetLargerImageLink(data.Pictures.Sizes)

		return ChannelFetchData{
			Title:       data.Name,
			Description: data.Bio,
			PublishedAt: data.CreatedTime,
			IconURL:     coverURL,
			CoverURL:    coverURL,
		}, nil
	}

	return ChannelFetchData{}, errors.New("not supported source")
}
