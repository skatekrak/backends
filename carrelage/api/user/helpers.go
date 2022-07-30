package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/carrelage/models"
)

func (c *Controller) getUser(ctx *fiber.Ctx) *models.User {
	return ctx.Locals(USER_LOADER).(*models.User)
}
