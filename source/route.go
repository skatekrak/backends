package source

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/scribe/fetchers"
	"github.com/skatekrak/scribe/fetchers/vimeo"
	"github.com/skatekrak/scribe/fetchers/youtube"
	"github.com/skatekrak/scribe/middlewares"
	"gorm.io/gorm"
)

type FindAllQuery struct {
	Types []string `form:"types[]" validate:"dive,eq=vimeo|eq=youtube|eq=rss"`
}

type CreateBody struct {
	URL           string `json:"url" validated:"required"`
	LangIsoCode   string `json:"lang" validate:"required"`
	IsSkateSource bool   `json:"isSkateSource"`
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

	youtubeFetcher := youtube.NewYoutubeFetcher(os.Getenv("YOUTUBE_API_KEY"))
	vimeoFetcher := vimeo.NewVimeoSourceFetcher(os.Getenv("VIMEO_API_KEY"))

	controller := NewController(db, []fetchers.SourceFetcher{youtubeFetcher, vimeoFetcher})
	auth := middlewares.Authorization(apiKey)

	router := app.Group("sources")

	router.Get("", middlewares.QueryHandler[FindAllQuery](), controller.FindAll)
	router.Post("", auth, middlewares.JSONHandler[CreateBody](), controller.Create)
	router.Patch("/:sourceID", auth, controller.LoaderHandler(), middlewares.JSONHandler[UpdateBody](), controller.Update)
	router.Delete("/:sourceID", auth, controller.LoaderHandler(), controller.Delete)

	router.Post("/:sourceID/refresh", auth, controller.LoaderHandler(), controller.RefreshSource)
}
