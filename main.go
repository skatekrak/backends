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
		log.Fatal("Error loading .env file")
	}

	db, err := database.Open("./local.db")
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.AutoMigrate(&lang.Lang{}); err != nil {
		log.Println(err)
	}

	r := gin.Default()

	langRouter := r.Group("/langs")
	{
		langRouter.GET("", lang.FindAll)
		langRouter.Use(Authorization()).POST("", lang.Create)
		langRouter.Use(Authorization()).PATCH("/:isoCode", lang.Update)
		langRouter.Use(Authorization()).DELETE("/:isoCode", lang.Delete)
	}

	r.Run()
}

type authHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var header authHeader
		if err := c.BindHeader(&header); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Missing Authorization header",
			})
			return
		}

		apiKey := os.Getenv("API_KEY")

		if header.Authorization != apiKey {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Wait, that's illegal",
			})
			return
		}

		c.Next()
	}
}
