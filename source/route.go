package source

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skatekrak/scribe/middlewares"
	"gorm.io/gorm"
)

type FindAllQuery struct {
	Types []string `form:"types[]" binding:"dive,eq=vimeo|eq=youtube|eq=rss"`
}

type CreateBody struct {
	ChannelID     string `json:"channelID" binding:"required"`
	LangIsoCode   string `json:"lang" binding:"required"`
	IsSkateSource bool   `json:"isSkateSource"`
	Type          string `json:"type" binding:"oneof=youtube vimeo feedly"`
}

type SourceURI struct {
	SourceID string `uri:"sourceID" binding:"required"`
}

type UpdateBody struct {
	LangIsoCode   string `json:"lang" binding:"required"`
	IsSkateSource bool   `json:"isSkateSource"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	ShortTitle    string `json:"shortTitle"`
	IconURL       string `json:"iconURL"`
	CoverURL      string `json:"coverURL"`
	WebsiteURL    string `json:"websiteURL"`
}

func Route(r *gin.Engine, db *gorm.DB) {
	apiKey := os.Getenv("API_KEY")
	controller := NewController(db)
	auth := middlewares.Authorization(apiKey)

	router := r.Group("sources")
	{
		router.GET("", middlewares.QueryHandler[FindAllQuery](), controller.FindAll)
		router.Use(auth)
		{
			router.POST("", middlewares.JSONHandler[CreateBody](), controller.Create)

			router.Use(middlewares.URIHandler[SourceURI](), controller.LoaderHandler())
			{

				router.PATCH("/:sourceID", middlewares.JSONHandler[UpdateBody](), controller.Update)
				router.DELETE("/:sourceID", controller.Delete)
			}
		}
	}
}
