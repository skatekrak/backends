package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/skatekrak/scribe/database"
	"github.com/skatekrak/scribe/lang"
	"github.com/skatekrak/scribe/middlewares"
	"github.com/skatekrak/scribe/model"
	"github.com/skatekrak/scribe/source"
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

	if err = db.AutoMigrate(&model.Lang{}, &model.Source{}); err != nil {
		log.Fatalf("unable to migrate database: %s", err)
	}

	r := gin.Default()
	r.Use(middlewares.ErrorHandler())
	lang.Route(r, db)
	source.Route(r, db)

	if err := r.Run(); err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}
