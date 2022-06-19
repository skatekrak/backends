package lang

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skatekrak/scribe/middlewares"
	"gorm.io/gorm"
)

type CreateBody struct {
	IsoCode  string `json:"isoCode" binding:"required,len=2"`
	ImageURL string `json:"imageURL" binding:"required"`
}

type UpdateBody struct {
	ImageURL string `json:"imageURL" binding:"required"`
}

type LangUri struct {
	IsoCode string `uri:"isoCode" binding:"required"`
}

func Route(r *gin.Engine, db *gorm.DB) {
	apiKey := os.Getenv("API_KEY")
	controller := NewController(db)
	auth := middlewares.Authorization(apiKey)

	router := r.Group("/langs")
	{

		router.GET("", controller.FindAll)
		router.Use(auth)
		{
			router.POST("", middlewares.JSONHandler[CreateBody](), controller.Create)

			router.Use(middlewares.URIHandler[LangUri]())
			{
				router.PATCH("/:isoCode", middlewares.JSONHandler[UpdateBody](), controller.Update)
				router.DELETE("/:isoCode", controller.Delete)
			}
		}
	}
}
