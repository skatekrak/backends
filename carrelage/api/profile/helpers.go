package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/carrelage/models"
)

func (c *Controller) getProfile(ctx *fiber.Ctx) *models.Profile {
	return ctx.Locals(PROFILE_LOADER).(*models.Profile)
}
