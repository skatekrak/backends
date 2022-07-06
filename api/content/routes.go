package content

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/services"
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

	router.Get("", middlewares.QueryHandler[FindQuery](), controller.Find)
}
