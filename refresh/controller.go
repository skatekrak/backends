package refresh

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/services"
)

type Controller struct {
	rs               *services.RefreshService
	ss               *services.SourceService
	cs               *services.ContentService
	fetcher          *fetchers.Fetcher
	feedlyCategoryID string
}

func (c *Controller) RefreshByTypes(ctx *fiber.Ctx) error {
	query := ctx.Locals(middlewares.QUERY).(RefreshQuery)

	contents, err := c.rs.RefreshByTypes(query.Types)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(contents)
}

func (c *Controller) RefreshSource(ctx *fiber.Ctx) error {
	source := middlewares.GetSource(ctx)

	contents, errs := c.rs.RefreshBySource(source)
	if errs != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errs)
	}

	return ctx.Status(fiber.StatusOK).JSON(contents)
}

// Refresh feedly sources
// @Summary   Query sources used in feedly and add missing ones in Scribe
// @Security ApiKeyAuth
// @Tags      sources
// @Success   200  {array}   []model.Source
// @Failure   500  {object}  api.JSONError
// @Router    /sources/sync-feedly [patch]
func (c *Controller) RefreshFeedly(ctx *fiber.Ctx) error {
	sources, err := c.rs.RefreshFeedlySource()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(sources)
}
