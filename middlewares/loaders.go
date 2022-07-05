package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/model"
	"github.com/skatekrak/scribe/services"
)

// Key used to pass the source interface between middlewares
const SOURCE_LOADER_LOCAL = "sourceID"

func SourceLoader(s *services.SourceService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sourceID := ctx.Params(SOURCE_LOADER_LOCAL)

		source, err := s.Get(sourceID)
		if err != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Source not found",
			})
		}

		ctx.Locals(SOURCE_LOADER_LOCAL, source)
		return ctx.Next()
	}
}

func GetSource(ctx *fiber.Ctx) model.Source {
	return ctx.Locals(SOURCE_LOADER_LOCAL).(model.Source)
}
