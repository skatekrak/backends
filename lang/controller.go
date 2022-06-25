package lang

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

const context_lang = "lang"

type Controller struct {
	s *Service
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{s: NewService(db)}
}

func (c *Controller) LoaderHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		isoCode := ctx.Params("isoCode")

		lang, err := c.s.Get(isoCode)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Lang not found")
		}

		ctx.Locals(context_lang, lang)
		return ctx.Next()
	}

}
func (c *Controller) FindAll(ctx *fiber.Ctx) error {
	langs, err := c.s.FindAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(langs)
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
	body := ctx.Locals(middlewares.BODY).(CreateBody)

	if _, err := c.s.Get(body.IsoCode); err == nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "isoCode already used",
		})
	}

	lang := model.Lang{
		IsoCode:  body.IsoCode,
		ImageURL: body.ImageURL,
	}

	if err := c.s.Create(&lang); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(lang)
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
	lang := ctx.Locals(context_lang).(model.Lang)
	body := ctx.Locals(middlewares.BODY).(UpdateBody)

	lang.ImageURL = body.ImageURL
	if err := c.s.Update(&lang); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(lang)
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	lang := ctx.Locals(context_lang).(model.Lang)

	if err := c.s.Delete(&lang); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Lang deleted",
	})
}
