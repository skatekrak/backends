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
	"github.com/skatekrak/scribe/api/refresh"
	"github.com/skatekrak/scribe/api/source"
	"github.com/skatekrak/scribe/database"
	_ "github.com/skatekrak/scribe/docs"
	"github.com/skatekrak/scribe/model"
	"github.com/skatekrak/scribe/services"
	"gorm.io/gorm"
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

	if err = db.AutoMigrate(&model.Lang{}, &model.Source{}, &model.Content{}, &model.Config{}); err != nil {
		log.Fatalf("unable to migrate database: %s", err)
	}

	setupConfig(db)

	app := fiber.New()
	setupRoutes(db, app)

	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatalln("Error listening")
		log.Fatalln(err)
	}
}

func setupConfig(db *gorm.DB) {
	// Setup necessary config key
	configService := services.NewConfigService(db)
	if err := configService.InitSetup(); err != nil {
		log.Fatalf("Unable to init config: %s", err)
	}
}

func setupRoutes(db *gorm.DB, app *fiber.App) {
	app.Use(logger.New())
	app.Use(recover.New())

	lang.Route(app, db)
	source.Route(app, db)
	content.Route(app, db)
	refresh.Route(app, db)

	app.Get("/docs/*", swagger.HandlerDefault)
}
