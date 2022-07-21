package loaders

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

const LANG_LOADER_LOCAL = "isoCode"

func LangLoader(s *services.LangService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		isoCode := ctx.Params("isoCode")

		lang, err := s.Get(isoCode)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Lang not found")
		}

		ctx.Locals(LANG_LOADER_LOCAL, lang)
		return ctx.Next()
	}
}

func GetLang(ctx *fiber.Ctx) model.Lang {
	return ctx.Locals(LANG_LOADER_LOCAL).(model.Lang)
}

const CONTENT_LOADER_LOCAL = "contentId"

func ContentLoader(s *services.ContentService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		contentId := ctx.Params(CONTENT_LOADER_LOCAL)

		content, err := s.Get(contentId)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Content not found")
		}

		ctx.Locals(CONTENT_LOADER_LOCAL, content)
		return ctx.Next()
	}
}
