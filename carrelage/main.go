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
	"github.com/skatekrak/carrelage/models"
	"github.com/skatekrak/utils/database"
)

func main() {
	db, err := database.Open(os.Getenv("POSTGRESQL_ADDON_URI"))
	if err != nil {
		log.Fatalf("unable to open database: %s", err)
	}

	if err = db.AutoMigrate(&models.User{}, &models.UserSubscription{}, &models.Profile{}); err != nil {
		log.Fatalf("unable to migrate database: %s", err)
	}

	app := fiber.New()

	// Setup middlewares
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ORIGINS"),
	}))
	app.Use(compress.New())
	app.Use(recover.New())

	// Start server
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatalln("Error listening")
		log.Fatalln(err)
	}
}
