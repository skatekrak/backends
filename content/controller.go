package content

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/middlewares"
	"gorm.io/gorm"
)

type Controller struct {
	s *Service
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{
		s: NewService(db),
	}
}

func (c *Controller) Find(ctx *fiber.Ctx) error {
	query := ctx.Locals(middlewares.QUERY).(FindQuery)

	contents, err := c.s.Find(query.SourceTypes, query.Page)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(contents)
}
