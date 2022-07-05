package refresh

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/clients/vimeo"
	"github.com/skatekrak/scribe/clients/youtube"
	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/middlewares"
	"gorm.io/gorm"
)

type RefreshQuery struct {
	Types []string `query:"types" validate:"required,dive,eq=vimeo|eq=youtube|eq=rss"`
}

func Route(app *fiber.App, db *gorm.DB) {
	apiKey := os.Getenv("API_KEY")

	youtubeClient := youtube.New(os.Getenv("YOUTUBE_API_KEY"))
	vimeoClient := vimeo.New(os.Getenv("VIMEO_API_KEY"))
	fetcher := fetchers.New(vimeoClient, youtubeClient, nil)

	controller := New(db, fetcher, os.Getenv("FEEDLY_FETCH_CATEGORY_ID"))
	auth := middlewares.Authorization(apiKey)

	router := app.Group("refresh")

	router.Post("", auth, middlewares.QueryHandler[RefreshQuery](), controller.RefreshByTypes)
	router.Post("/:sourceID", auth, controller.LoaderHandler(), controller.RefreshSource)
}
