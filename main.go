package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/skatekrak/scribe/api/content"
	"github.com/skatekrak/scribe/api/lang"
	"github.com/skatekrak/scribe/api/source"
	"github.com/skatekrak/scribe/database"
	_ "github.com/skatekrak/scribe/docs"
	"github.com/skatekrak/scribe/model"
	"github.com/skatekrak/scribe/refresh"
)

// @title         Scribe API
// @version       1.0
// @description   Document for the Scribe API
// @license.name  AGPLv3
// @host          localhost:8080
// @BasePath      /
// @Accept        json
// @Produce       json
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	db, err := database.Open("./local.db")
	if err != nil {
		log.Fatalf("unable to open database: %s", err)
	}

	if err = db.AutoMigrate(&model.Lang{}, &model.Source{}, &model.Content{}); err != nil {
		log.Fatalf("unable to migrate database: %s", err)
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	lang.Route(app, db)
	source.Route(app, db)
	content.Route(app, db)
	refresh.Route(app, db)

	app.Get("/docs/*", swagger.HandlerDefault)

	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatalln("Error listening")
		log.Fatalln(err)
	}

}
