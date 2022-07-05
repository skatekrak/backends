package services

import (
	"errors"
	"time"

	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/helpers"
	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

type RefreshErrors struct {
	Message string           `json:"message,omitempty"`
	Error   error            `json:"error,omitempty"`
	Errors  map[string]error `json:"errors,omitempty"`
}

type RefreshService struct {
	fetcher          *fetchers.Fetcher
	feedlyCategoryID string
	cs               *ContentService
	ss               *SourceService
}

func NewRefreshService(db *gorm.DB, fetcher *fetchers.Fetcher, feedlyCategoryID string) *RefreshService {
	return &RefreshService{
		fetcher:          fetcher,
		feedlyCategoryID: feedlyCategoryID,
		cs:               NewContentService(db),
		ss:               NewSourceService(db),
	}
}

func (rs *RefreshService) RefreshByTypes(types []string) ([]*model.Content, error) {
	sources, err := rs.ss.FindAll(types)
	if err != nil {
		return []*model.Content{}, err
	}

	if len(sources) <= 0 {
		return []*model.Content{}, errors.New("no sources to update")
	}

	formattedContents := []*model.Content{}
	errs := make(map[uint]error)
	now := time.Now()

	for _, source := range sources {
		if source.SourceType == "youtube" || source.SourceType == "vimeo" {
			contents, err := rs.fetcher.FetchChannelContents(source.SourceID, source.SourceType)

			if err != nil {
				errs[source.ID] = err
				continue
			}

			for _, content := range contents {
				// Only add content not already here
				if _, err := rs.cs.FindOneByContentID(content.ContentID); err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						formattedContents = append(formattedContents, formatContent(content, source))
					}
				}
			}

			source.RefreshedAt = &now
		}
	}

	if helpers.Has(types, "rss") {
		contents, err := rs.fetcher.FetchFeedlyContents(rs.feedlyCategoryID)
		if err != nil {
			return []*model.Content{}, err
		}

		for _, content := range contents {
			// Only add content not already here
			if _, err := rs.cs.FindOneByContentID(content.ContentID); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Look into sources we've aleady fetched
					source, ok := helpers.Find(sources, func(s *model.Source) bool {
						return s.SourceID == content.SourceID
					})

					if ok {
						formattedContents = append(formattedContents, formatContent(content, source))

						source.RefreshedAt = &now
					}
				}
			}
		}
	}

	if len(errs) > 0 {
		return []*model.Content{}, errors.New("errors fetching contents for some video channels")
	}

	if err := rs.cs.AddMany(formattedContents, sources); err != nil {
		return []*model.Content{}, err
	}

	return formattedContents, nil
}

func (rs *RefreshService) RefreshBySource(source model.Source) ([]*model.Content, *RefreshErrors) {

	if source.SourceType == "rss" {
		return []*model.Content{}, &RefreshErrors{Error: errors.New("rss sources cannot be individually refreshed")}
	}

	contentsMap := map[string][]fetchers.ContentFetchData{}
	if source.SourceType == "youtube" {
		if errors := rs.fetcher.FetchYoutubeContent([]string{source.SourceID}, contentsMap); len(errors) > 0 {
			return []*model.Content{}, &RefreshErrors{Errors: errors}
		}
	} else if source.SourceType == "vimeo" {
		if errors := rs.fetcher.FetcherVimeoContent([]string{source.SourceID}, contentsMap); len(errors) > 0 {
			return []*model.Content{}, &RefreshErrors{Errors: errors}
		}
	} else {
		return []*model.Content{}, &RefreshErrors{Error: errors.New("Oops")}
	}

	formattedContents := []*model.Content{}

	for _, content := range contentsMap[source.SourceID] {
		// Only add content not already here
		if _, err := rs.cs.FindOneByContentID(content.ContentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				formattedContents = append(formattedContents, formatContent(content, &source))
			}
		}
	}

	now := time.Now()
	source.RefreshedAt = &now

	if err := rs.cs.AddMany(formattedContents, []*model.Source{&source}); err != nil {
		return []*model.Content{}, &RefreshErrors{Error: err}
	}

	return formattedContents, nil
}

func formatContent(content fetchers.ContentFetchData, source *model.Source) *model.Content {
	contentType := "video"
	if source.SourceType == "rss" {
		contentType = "article"
	}

	return &model.Content{
		SourceID:     source.ID,
		ContentID:    content.ContentID,
		PublishedAt:  content.PublishedAt,
		Title:        content.Title,
		ThumbnailURL: content.ThumbnailURL,
		ContentURL:   content.ContentURL,
		RawSummary:   content.RawDescription,
		Summary:      content.Description,
		Type:         contentType,
	}
}
