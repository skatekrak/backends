package refresh

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/clients/feedly"
	"github.com/skatekrak/scribe/clients/vimeo"
	"github.com/skatekrak/scribe/clients/youtube"
	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/services"
	"gorm.io/gorm"
)

type RefreshQuery struct {
	Types []string `query:"types" validate:"required,dive,eq=vimeo|eq=youtube|eq=rss"`
}

func Route(app *fiber.App, db *gorm.DB) {
	apiKey := os.Getenv("API_KEY")

	youtubeClient := youtube.New(os.Getenv("YOUTUBE_API_KEY"))
	vimeoClient := vimeo.New(os.Getenv("VIMEO_API_KEY"))
	feedlyClient := feedly.New(os.Getenv("FEEDLY_API_KEY"))
	fetcher := fetchers.New(vimeoClient, youtubeClient, feedlyClient)

	sourceService := services.NewSourceService(db)
	contentService := services.NewContentService(db)

	controller := &Controller{
		ss:               sourceService,
		cs:               contentService,
		fetcher:          fetcher,
		feedlyCategoryID: os.Getenv("FEEDLY_FETCH_CATEGORY_ID"),
	}
	auth := middlewares.Authorization(apiKey)
	sourceLoader := middlewares.SourceLoader(sourceService)

	router := app.Group("refresh")

	router.Post("", auth, middlewares.QueryHandler[RefreshQuery](), controller.RefreshByTypes)
	router.Post("/:sourceID", auth, sourceLoader, controller.RefreshSource)
}
