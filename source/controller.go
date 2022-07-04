package source

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/content"
	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/helpers"
	"github.com/skatekrak/scribe/lang"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/model"
	"github.com/skatekrak/scribe/vendors/clients/vimeo"
	"github.com/skatekrak/scribe/vendors/clients/youtube"
	"gorm.io/gorm"
)

// Key used to pass the source interface between middlewares
const context_source = "source"

type Controller struct {
	s       *Service
	ls      *lang.Service
	cs      *content.Service
	fetcher *fetchers.Fetcher
}

func NewController(db *gorm.DB, fetcher *fetchers.Fetcher) *Controller {
	return &Controller{
		s:       NewService(db),
		ls:      lang.NewService(db),
		cs:      content.NewService(db),
		fetcher: fetcher,
	}
}

func (c *Controller) LoaderHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sourceID := ctx.Params("sourceID")

		source, err := c.s.Get(sourceID)
		if err != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Lang not founc",
			})
		}

		ctx.Locals(context_source, source)
		return ctx.Next()
	}
}

func (c *Controller) FindAll(ctx *fiber.Ctx) error {
	query := ctx.Locals(middlewares.QUERY).(FindAllQuery)

	sources, err := c.s.FindAll(query.Types)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(sources)
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
	body := ctx.Locals(middlewares.BODY).(CreateBody)

	var sourceID string

	if body.Type == "youtube" && !youtube.IsYoutubeChannel(body.URL) {
		return ctx.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
			"message": "This isn't a youtube url",
		})
	}
	if body.Type == "vimeo" && !vimeo.IsVimeoUser(body.URL) {
		return ctx.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
			"message": "This isn't a vimeo url",
		})
	}

	sourceID, err := c.fetcher.GetSourceID(body.URL)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "This url seems invalid or not supported",
		})
	}

	if _, err := c.s.GetBySourceID(sourceID); err == nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": fmt.Sprintf("This %s source is already added", body.Type),
		})
	}

	nextOrder, err := c.s.GetNextOrder()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Couldn't process the next order",
			"error":   err.Error(),
		})
	}

	data, err := c.fetcher.FetchChannelData(body.URL)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	source := model.Source{
		Order:       nextOrder,
		SourceType:  body.Type,
		SkateSource: body.IsSkateSource,
		LangIsoCode: body.LangIsoCode,
		SourceID:    sourceID,
		Title:       data.Title,
		ShortTitle:  data.Title,
		Description: data.Description,
		CoverURL:    data.CoverURL,
		IconURL:     data.IconURL,
		WebsiteURL:  body.URL,
		PublishedAt: &data.PublishedAt,
	}

	if err := c.s.Create(&source); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Couldn't create the source",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(source)
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
	body := ctx.Locals(middlewares.BODY).(UpdateBody)
	source := ctx.Locals(context_source).(model.Source)

	source.LangIsoCode = helpers.SetIfNotNil(body.LangIsoCode, source.LangIsoCode)
	source.SkateSource = helpers.SetIfNotNil(body.IsSkateSource, source.SkateSource)
	source.Title = helpers.SetIfNotNil(body.Title, source.Title)
	source.ShortTitle = helpers.SetIfNotNil(body.ShortTitle, source.ShortTitle)
	source.Description = helpers.SetIfNotNil(body.Description, source.Description)
	source.IconURL = helpers.SetIfNotNil(body.IconURL, source.IconURL)
	source.CoverURL = helpers.SetIfNotNil(body.CoverURL, source.CoverURL)
	source.WebsiteURL = helpers.SetIfNotNil(body.WebsiteURL, source.WebsiteURL)

	if err := c.s.Update(&source); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(source)
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	source := ctx.Locals(context_source).(model.Source)

	if err := c.s.Delete(&source); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Couldn't delete this source",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Source deleted",
	})
}

func (c *Controller) RefreshSource(ctx *fiber.Ctx) error {
	source := ctx.Locals(context_source).(model.Source)

	if source.SourceType == "rss" {
		return ctx.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"message": "rss sources cannot be individually refreshed",
		})
	}

	contents := []fetchers.ContentFetchData{}
	if source.SourceType == "youtube" {
		if errors := c.fetcher.FetchYoutubeContent([]string{source.SourceID}, &contents); len(errors) > 0 {
			return ctx.Status(fiber.StatusInternalServerError).JSON(errors)
		}
	} else if source.SourceType == "vimeo" {
		if errors := c.fetcher.FetcherVimeoContent([]string{source.SourceID}, contents); len(errors) > 0 {
			return ctx.Status(fiber.StatusInternalServerError).JSON(errors)
		}
	} else {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Oops",
		})
	}

	formattedContents := []*model.Content{}

	for _, content := range contents {
		// Only add content not already here
		if _, err := c.cs.FindOneByContentID(content.ContentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				formattedContents = append(formattedContents, &model.Content{
					SourceID:     source.ID,
					ContentID:    content.ContentID,
					PublishedAt:  content.PublishedAt,
					Title:        content.Title,
					ThumbnailURL: content.ThumbnailURL,
					ContentURL:   content.ContentURL,
					RawSummary:   content.RawDescription,
					Summary:      content.Description,
					Type:         "video",
				})
			}
		}
	}

	if err := c.cs.AddMany(formattedContents); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(formattedContents)
}

func (c *Controller) RefreshTypes(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "",
	})
}
