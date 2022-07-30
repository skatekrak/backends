package user

import "github.com/gofiber/fiber/v2"

const USER_LOADER = "userID"

// Will fetch the user, if it doesn't exists will return a 404,
// or go to the next handler if it does
func (c *Controller) Loader() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID := ctx.Params(USER_LOADER)

		user, err := c.usersService.Get(userID)
		if err != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}

		ctx.Locals(USER_LOADER, user)
		return ctx.Next()
	}
}
