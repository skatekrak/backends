package source

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

	router := r.Group("sources")
	{
		router.GET("", controller.FindAll)
		router.Use(auth)
		{
			router.POST("", controller.Create)
		}
	}
}
