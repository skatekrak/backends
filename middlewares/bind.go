package middlewares

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/formatter"
)

const BODY = "body"
const QUERY = "query"
const URI = "uri"

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func JSONHandler[T interface{}]() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var body T

		if err := ctx.BodyParser(&body); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := validate.Struct(body); err != nil {
			if verr, ok := err.(validator.ValidationErrors); ok {
				f := formatter.NewJSONFormatter()
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": f.Simple(verr),
				})
			}
		}

		ctx.Locals("body", body)
		return ctx.Next()
	}
}

func QueryHandler[T interface{}]() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var query T

		if err := ctx.QueryParser(&query); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := validate.Struct(query); err != nil {
			if verr, ok := err.(validator.ValidationErrors); ok {
				f := formatter.NewJSONFormatter()
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": f.Simple(verr),
				})
			}
		}

		ctx.Locals("query", query)
		return ctx.Next()
	}
}
