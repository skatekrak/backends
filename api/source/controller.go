package source

import (
	"errors"
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
	"gorm.io/gorm"
)

type Controller struct {
	s                *services.SourceService
	ls               *services.LangService
	cs               *services.ContentService
	fetcher          *fetchers.Fetcher
	feedlyCategoryID string
}

// Fetch all sources
// @Summary  Fetch all sources
// @Tags      sources
// @Success  200    {array}   []model.Source
// @Failure  500    {object}  api.JSONError
// @Param    types  query     []string  false  "Filter by source types"  Enums(rss,vimeo,youtube)
// @Router   /sources [get]
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

// Add a new source
// @Summary   Add a new source
// @Tags      sources
// @Security ApiKeyAuth
// @Success   200   {object}  model.Source
// @Failure   404   {object}  api.JSONError
// @Failure   500   {object}  api.JSONError
// @Param     body  body      source.CreateBody  true  "body"
// @Router    /sources [post]
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
		PublishedAt: data.PublishedAt,
	}

	if err := c.s.Create(&source); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Couldn't create the source",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(source)
}

// Update a source
// @Summary   Update a source
// @Security ApiKeyAuth
// @Tags     sources
// @Success   200       {object}  model.Source
// @Failure   500       {object}  api.JSONError
// @Param     body      body      source.UpdateBody  true  "Update body"
// @Param     sourceID  path      integer     true  "ID of the source"
// @Router    /sources/{sourceID} [patch]
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

// Delete a source
// @Summary   Delete a source
// @Security ApiKeyAuth
// @Tags      sources
// @Success   200       {object}  api.JSONMessage
// @Failure   404       {object}  api.JSONError
// @Failure   500       {object}  api.JSONError
// @Param     sourceID  path      integer  true  "ID of the source"
// @Router    /sources/{sourceID} [delete]
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

// Refresh feedly sources
// @Summary   Query sources used in feedly and add missing ones in Scribe
// @Security ApiKeyAuth
// @Tags      sources
// @Success   200  {array}   []model.Source
// @Failure   500  {object}  api.JSONError
// @Router    /sources/sync-feedly [patch]
func (c *Controller) RefreshFeedly(ctx *fiber.Ctx) error {
	data, err := c.fetcher.FetchFeedlySources(c.feedlyCategoryID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	nextOrder, err := c.s.GetNextOrder()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Couldn't get the next order",
			"error":   err.Error(),
		})
	}

	sources := []*model.Source{}
	index := 0

	for _, s := range data {
		if _, err := c.s.GetBySourceID(s.SourceID); err != nil {
			// Only attempt to create source that are not already here
			if errors.Is(err, gorm.ErrRecordNotFound) {
				sources = append(sources, &model.Source{
					Order:       nextOrder + index,
					SourceType:  "rss",
					SourceID:    s.SourceID,
					Title:       s.Title,
					Description: s.Description,
					ShortTitle:  s.Title,
					CoverURL:    s.CoverURL,
					IconURL:     s.IconURL,
					WebsiteURL:  s.WebsiteURL,
					SkateSource: s.SkateSource,
					PublishedAt: s.PublishedAt,
					Lang: model.Lang{
						IsoCode: s.Lang,
					},
				})
				index++
			}
		}
	}

	if err := c.s.AddMany(sources); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(sources)
}

// Update orders of the sources
// @Summary   Update orders of the sources
// @Security ApiKeyAuth
// @Tags      sources
// @Success   200   {array}   []model.Source
// @Failure   400   {object}  api.JSONError
// @Failure   500   {object}  api.JSONError
// @Param     body  body      UpdateBody  true  "Update body"
// @Router    /sources/order [patch]
func (c *Controller) UpdateOrder(ctx *fiber.Ctx) error {
	body := ctx.Locals(middlewares.BODY).(UpdateOrderBody)

	orders := make([]int, len(body))
	updates := map[int]map[string]interface{}{}

	for key, order := range body {
		if helpers.Has(orders, order) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "You have an duplicate order",
			})
		}

		orders = append(orders, order)

		updates[key] = map[string]interface{}{
			"order": order,
		}
	}

	sources, err := c.s.UpdateOrder(updates)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(sources)
}
