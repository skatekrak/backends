package content

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/services"
)

type Controller struct {
	s *services.ContentService
}

// Find contents
// @Tags     contents
// @Param    sourceTypes  query     []string  false  "filter contents by source types"  Enums(rss,vimeo,youtube)
// @Param    sources      query     []int     false  "filter contents by source id"
// @Param    page         query     int       false  "Fetch page"  minimum(1)
// @Success  200          {object}  database.Pagination{Items=[]model.Content}
// @Failure  500          {object}  api.JSONError
// @Router   /contents [get]
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