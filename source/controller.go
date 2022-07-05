package source

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/clients/vimeo"
	"github.com/skatekrak/scribe/clients/youtube"
	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/helpers"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/model"
	"github.com/skatekrak/scribe/services"
)

type Controller struct {
	s                *services.SourceService
	ls               *services.LangService
	cs               *services.ContentService
	fetcher          *fetchers.Fetcher
	feedlyCategoryID string
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
	source := middlewares.GetSource(ctx)

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
	source := middlewares.GetSource(ctx)

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
