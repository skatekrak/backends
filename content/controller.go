package content

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/services"
)

type Controller struct {
	s *services.ContentService
}

func (c *Controller) Find(ctx *fiber.Ctx) error {
	query := ctx.Locals(middlewares.QUERY).(FindQuery)

	pagination, err := c.s.Find(query.SourceTypes, query.Sources, query.Page)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(pagination)
}
