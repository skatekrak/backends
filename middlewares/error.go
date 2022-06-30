package middlewares

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/formatter"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		if verr, ok := err.(validator.ValidationErrors); ok {
			f := formatter.NewJSONFormatter()
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": f.Simple(verr),
			})
		}

		return nil
	}
}
