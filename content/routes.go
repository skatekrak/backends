package content

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/middlewares"
	"gorm.io/gorm"
)

type FindQuery struct {
	SourceTypes []string `json:"sourceTypes" validate:"dive,eq=vimeo|eq=youtube|eq=rss"`
	Page        int      `json:"page"`
}

func Route(app *fiber.App, db *gorm.DB) {
	controller := NewController(db)

	router := app.Group("contents")

	router.Get("", middlewares.QueryHandler[FindQuery](), controller.Find)
}
