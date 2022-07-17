package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/skatekrak/scribe/api/content"
	"github.com/skatekrak/scribe/api/lang"
	"github.com/skatekrak/scribe/api/refresh"
	"github.com/skatekrak/scribe/api/source"
	"github.com/skatekrak/scribe/database"
	_ "github.com/skatekrak/scribe/docs"
	"github.com/skatekrak/scribe/jobs"
	"github.com/skatekrak/scribe/model"
	"github.com/skatekrak/scribe/services"
	"gorm.io/gorm"
)

// @title                       Scribe API
// @version                     1.0
// @description                 Document for the Scribe API
// @license.name                AGPLv3
// @host                        localhost:8080
// @BasePath                    /
// @Accept                      json
// @Produce                     json
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
func main() {
	db, err := database.Open(os.Getenv("POSTGRESQL_ADDON_URI"))
	if err != nil {
		log.Fatalf("unable to open database: %s", err)
	}

	if err = db.AutoMigrate(&model.Lang{}, &model.Source{}, &model.Content{}, &model.Config{}); err != nil {
		log.Fatalf("unable to migrate database: %s", err)
	}

	setupConfig(db)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ORIGINS"),
	}))
	app.Use(compress.New())
	setupRoutes(db, app)

	jobs.Setup(db)

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
