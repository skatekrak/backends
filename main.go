package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/skatekrak/scribe/database"
	"github.com/skatekrak/scribe/lang"
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

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	lang.Route(app, db)
	source.Route(app, db)

	app.Listen(":8080")
}
