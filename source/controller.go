package source

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/lang"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

// Key used to pass the source interface between middlewares
const context_source = "source"

type Controller struct {
	s  *Service
	ls *lang.Service
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{
		s:  NewService(db),
		ls: lang.NewService(db),
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

	if _, err := c.s.GetBySourceID(body.ChannelID); err == nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "This youtube channel is already added",
		})
	}

	nextOrder, err := c.s.GetNextOrder()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Couldn't process the next order",
			"error":   err.Error(),
		})
	}

	// TODO: fetch youtube info

	source := model.Source{
		Order:       nextOrder,
		SourceType:  body.Type,
		SkateSource: body.IsSkateSource,
		LangIsoCode: body.LangIsoCode,
		SourceID:    body.ChannelID,
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

	// TODO
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"body":   body,
		"source": source,
	})
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
