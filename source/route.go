package source

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/middlewares"
	"gorm.io/gorm"
)

type FindAllQuery struct {
	Types []string `form:"types[]" binding:"dive,eq=vimeo|eq=youtube|eq=rss"`
}

type CreateBody struct {
	ChannelID     string `json:"channelID" validate:"required"`
	LangIsoCode   string `json:"lang" validate:"required"`
	IsSkateSource bool   `json:"isSkateSource"`
	Type          string `json:"type" validate:"oneof=youtube vimeo feedly"`
}

type SourceURI struct {
	SourceID string `uri:"sourceID" validate:"required"`
}

type UpdateBody struct {
	LangIsoCode   string `json:"lang" validate:"required"`
	IsSkateSource bool   `json:"isSkateSource"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	ShortTitle    string `json:"shortTitle"`
	IconURL       string `json:"iconURL"`
	CoverURL      string `json:"coverURL"`
	WebsiteURL    string `json:"websiteURL"`
}

func Route(app *fiber.App, db *gorm.DB) {
	apiKey := os.Getenv("API_KEY")
	controller := NewController(db)
	auth := middlewares.Authorization(apiKey)

	router := app.Group("sources")

	router.Get("", middlewares.QueryHandler[FindAllQuery](), controller.FindAll)
	router.Post("", auth, middlewares.JSONHandler[CreateBody](), controller.Create)
	router.Patch("/:sourceID", auth, controller.LoaderHandler(), middlewares.JSONHandler[UpdateBody](), controller.Update)
	router.Delete("/:sourceID", auth, controller.LoaderHandler(), controller.Delete)
}
