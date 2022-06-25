package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func Authorization(apiKey string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorization := ctx.Get("Authorization")

		if authorization == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing Authorization")
		}

		if authorization != apiKey {
			return fiber.NewError(fiber.StatusForbidden, "Wait, that's illegal")
		}

		return ctx.Next()
	}
}
