package refresh

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/content"
	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/helpers"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/model"
	"github.com/skatekrak/scribe/source"
	"gorm.io/gorm"
)

type Controller struct {
	ss               *source.Service
	cs               *content.Service
	fetcher          *fetchers.Fetcher
	feedlyCategoryID string
}

func New(db *gorm.DB, fetcher *fetchers.Fetcher, feedlyCategoryID string) *Controller {
	return &Controller{
		ss:               source.NewService(db),
		cs:               content.NewService(db),
		fetcher:          fetcher,
		feedlyCategoryID: feedlyCategoryID,
	}
}

func (c *Controller) RefreshByTypes(ctx *fiber.Ctx) error {
	query := ctx.Locals(middlewares.QUERY).(RefreshQuery)

	sources, err := c.ss.FindAll(query.Types)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if len(sources) <= 0 {
		return ctx.Status(fiber.StatusNotModified).JSON(fiber.Map{
			"message": "No sources to update",
		})
	}

	formattedContents := []*model.Content{}
	errs := make(map[uint]error)
	now := time.Now()

	for _, source := range sources {
		if source.SourceType == "youtube" || source.SourceType == "vimeo" {
			contents, err := c.fetcher.FetchChannelContents(source.SourceID, source.SourceType)

			if err != nil {
				errs[source.ID] = err
				continue
			}

			for _, content := range contents {
				// Only add content not already here
				if _, err := c.cs.FindOneByContentID(content.ContentID); err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						formattedContents = append(formattedContents, formatContent(content, source))
					}
				}
			}

			source.RefreshedAt = &now
		}
	}

	if helpers.Has(query.Types, "rss") {
		contents, err := c.fetcher.FetchFeedlyContents(c.feedlyCategoryID)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Coudln't refresh feedly",
				"error":   err.Error(),
			})
		}

		for _, content := range contents {
			// Only add content not already here
			if _, err := c.cs.FindOneByContentID(content.ContentID); err != nil {
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
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error refreshing sources",
			"errors":  errs,
		})
	}

	if err := c.cs.AddMany(formattedContents, sources); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(formattedContents)
}

func (c *Controller) RefreshSource(ctx *fiber.Ctx) error {
	source := ctx.Locals("sourceID").(model.Source)

	if source.SourceType == "rss" {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"message": "rss sources cannot be individually refreshed",
		})
	}

	contentsMap := map[string][]fetchers.ContentFetchData{}
	if source.SourceType == "youtube" {
		if errors := c.fetcher.FetchYoutubeContent([]string{source.SourceID}, contentsMap); len(errors) > 0 {
			return ctx.Status(fiber.StatusInternalServerError).JSON(errors)
		}
	} else if source.SourceType == "vimeo" {
		if errors := c.fetcher.FetcherVimeoContent([]string{source.SourceID}, contentsMap); len(errors) > 0 {
			return ctx.Status(fiber.StatusInternalServerError).JSON(errors)
		}
	} else {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Oops",
		})
	}

	formattedContents := []*model.Content{}

	for _, content := range contentsMap[source.SourceID] {
		// Only add content not already here
		if _, err := c.cs.FindOneByContentID(content.ContentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				formattedContents = append(formattedContents, formatContent(content, &source))
			}
		}
	}

	now := time.Now()
	source.RefreshedAt = &now

	if err := c.cs.AddMany(formattedContents, []*model.Source{&source}); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(formattedContents)
}

func formatContent(content fetchers.ContentFetchData, source *model.Source) *model.Content {
	return &model.Content{
		SourceID:     source.ID,
		ContentID:    content.ContentID,
		PublishedAt:  content.PublishedAt,
		Title:        content.Title,
		ThumbnailURL: content.ThumbnailURL,
		ContentURL:   content.ContentURL,
		RawSummary:   content.RawDescription,
		Summary:      content.Description,
		Type:         "video",
	}
}
