package source

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

type FindAllQuery struct {
	Types []string `query:"types" validate:"dive,eq=vimeo|eq=youtube|eq=rss"`
}

type CreateBody struct {
	URL           string `json:"url" validated:"required"`
	LangIsoCode   string `json:"lang" validate:"required"`
	IsSkateSource bool   `json:"isSkateSource"`
	Type          string `json:"type" validate:"required,oneof=vimeo youtube"`
}

type SourceURI struct {
	SourceID string `uri:"sourceID" validate:"required"`
}

type UpdateBody struct {
	LangIsoCode   *string `json:"lang"`
	IsSkateSource *bool   `json:"isSkateSource"`
	Title         *string `json:"title"`
	ShortTitle    *string `json:"shortTitle"`
	Description   *string `json:"description"`
	IconURL       *string `json:"iconURL"`
	CoverURL      *string `json:"coverURL"`
	WebsiteURL    *string `json:"websiteURL"`
}

func Route(app *fiber.App, db *gorm.DB) {
	apiKey := os.Getenv("API_KEY")

	youtubeClient := youtube.New(os.Getenv("YOUTUBE_API_KEY"))
	vimeoClient := vimeo.New(os.Getenv("VIMEO_API_KEY"))
	feedlyClient := feedly.New(os.Getenv("FEEDLY_API_KEY"))
	fetcher := fetchers.New(vimeoClient, youtubeClient, feedlyClient)

	sourceService := services.NewSourceService(db)
	contentService := services.NewContentService(db)
	langService := services.NewLangService(db)

	controller := &Controller{
		s:                sourceService,
		cs:               contentService,
		ls:               langService,
		fetcher:          fetcher,
		feedlyCategoryID: os.Getenv("FEEDLY_FETCH_CATEGORY_ID"),
	}
	auth := middlewares.Authorization(apiKey)
	sourceLoader := middlewares.SourceLoader(sourceService)

	router := app.Group("sources")

	router.Get("", middlewares.QueryHandler[FindAllQuery](), controller.FindAll)
	router.Post("", auth, middlewares.JSONHandler[CreateBody](), controller.Create)
	router.Post("/sync-feedly", auth, controller.RefreshFeedly)
	router.Patch("/:sourceID", auth, sourceLoader, middlewares.JSONHandler[UpdateBody](), controller.Update)
	router.Delete("/:sourceID", auth, sourceLoader, controller.Delete)

}
