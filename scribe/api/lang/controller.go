package lang

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/loaders"
	"github.com/skatekrak/scribe/model"
	"github.com/skatekrak/scribe/services"
	"github.com/skatekrak/utils/middlewares"
)

type Controller struct {
	s *services.LangService
}

// Fetch all langs
// @Tags      langs
// @Success  200  {array}   []model.Lang
// @Failure  500  {object}  api.JSONError
// @Router   /langs [get]
func (c *Controller) FindAll(ctx *fiber.Ctx) error {
	langs, err := c.s.FindAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(langs)
}

// Create a lang
// @Tags      langs
// @Security  ApiKeyAuth
// @Success   200   {object}  model.Lang
// @Failure   409   {object}  api.JSONError
// @Failure   500   {object}  api.JSONError
// @Param     body  body      lang.CreateBody  true  "Create body"
// @Router    /langs [post]
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

// Update a lang
// @Tags      langs
// @Security  ApiKeyAuth
// @Success   200      {object}  model.Lang
// @Failure   404      {object}  api.JSONError
// @Failure   500      {object}  api.JSONError
// @Param     body     body      lang.UpdateBody  true  "Update body"
// @Param     isoCode  path      string           true  "Lang ISO Code"
// @Router    /langs/{isoCode} [patch]
func (c *Controller) Update(ctx *fiber.Ctx) error {
	lang := loaders.GetLang(ctx)
	body := ctx.Locals(middlewares.BODY).(UpdateBody)

	lang.ImageURL = body.ImageURL
	if err := c.s.Update(&lang); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(lang)
}

// Delete a lang
// @Tags     langs
// @Security  ApiKeyAuth
// @Success   200      {object}  api.JSONMessage
// @Failure   404      {object}  api.JSONError
// @Failure   500      {object}  api.JSONError
// @Param     isoCode  path      string  true  "Lang ISO Code"
// @Router    /langs/{isoCode} [delete]
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	lang := loaders.GetLang(ctx)

	if err := c.s.Delete(&lang); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Lang deleted",
	})
}
