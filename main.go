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

	// r := gin.Default()
	// r.Use(middlewares.ErrorHandler())
	// lang.Route(r, db)
	// source.Route(r, db)

	// if err := r.Run(); err != nil {
	// 	log.Fatalf("error: %s", err.Error())
	// }

	app.Listen(":8080")
}
