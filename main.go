package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/skatekrak/scribe/database"
	"github.com/skatekrak/scribe/lang"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	db, err := database.Open("./local.db")
	if err != nil {
		log.Fatalf("unable to open database: %s", err)
	}

	if err = db.AutoMigrate(&lang.Lang{}); err != nil {
		log.Fatalf("unable to migrate database: %s", err)
	}

	r := gin.Default()

	apikey := os.Getenv("API_KEY")

	langRouter := r.Group("/langs")
	{
		controller := lang.NewController(db)
		auth := Authorization(apikey)
		langRouter.GET("", controller.FindAll)
		langRouter.Use(auth).POST("", controller.Create)
		langRouter.Use(auth).PATCH("/:isoCode", controller.Update)
		langRouter.Use(auth).DELETE("/:isoCode", controller.Delete)
	}

	if err := r.Run(); err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}

type authHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}

func Authorization(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header authHeader
		if err := c.BindHeader(&header); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Missing Authorization header",
			})
			return
		}

		if header.Authorization != apiKey {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Wait, that's illegal",
			})
			return
		}

		c.Next()
	}
}
