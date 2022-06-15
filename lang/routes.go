package lang

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skatekrak/scribe/middlewares"
	"gorm.io/gorm"
)

func Route(r *gin.Engine, db *gorm.DB) {
	apiKey := os.Getenv("API_KEY")
	controller := NewController(db)
	auth := middlewares.Authorization(apiKey)

	router := r.Group("/langs")
	{

		router.GET("", controller.FindAll)
		router.Use(auth)
		{
			router.POST("", controller.Create)
			router.PATCH("/:isoCode", controller.Update)
			router.DELETE("/:isoCode", controller.Delete)
		}
	}
}
