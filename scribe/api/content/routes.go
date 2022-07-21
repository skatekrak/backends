package content

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/loaders"
	"github.com/skatekrak/scribe/services"
	"github.com/skatekrak/utils/middlewares"
	"gorm.io/gorm"
)

type FindQuery struct {
	SourceTypes []string `json:"sourceTypes" validate:"dive,eq=vimeo|eq=youtube|eq=rss"`
	Sources     []int    `json:"sources"`
	Page        int      `json:"page"`
}

func Route(app *fiber.App, db *gorm.DB) {
	contentService := services.NewContentService(db)
	controller := &Controller{
		s: contentService,
	}

	router := app.Group("contents")

	contentLoader := loaders.ContentLoader(contentService)

	router.Get("", middlewares.QueryHandler[FindQuery](), controller.Find)
	router.Get("/:contentId", contentLoader, controller.Get)
}
