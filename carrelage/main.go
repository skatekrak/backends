package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/skatekrak/carrelage/api/profile"
	"github.com/skatekrak/carrelage/api/user"
	"github.com/skatekrak/carrelage/auth"
	_ "github.com/skatekrak/carrelage/docs"
	"github.com/skatekrak/carrelage/models"
	"github.com/skatekrak/carrelage/services"
	"github.com/skatekrak/utils/database"
	"github.com/skatekrak/utils/middlewares"
	"github.com/supertokens/supertokens-golang/supertokens"
	"gorm.io/gorm"
)

// @title         Carrelage API
// @version       3.0
// @description   Docs for Carrelage API
// @license.name  APGLv3
// @host          https://carrelage.api.skatekrak.com
// @BasePath      /
// @Accept        json
// @Produce       json
func main() {
	db, err := database.Open(os.Getenv("POSTGRESQL_ADDON_URI"))
	if err != nil {
		log.Fatalf("unable to open database: %s", err)
	}

	if err = db.AutoMigrate(&models.User{}, &models.UserSubscription{}, &models.Profile{}); err != nil {
		log.Fatalf("unable to migrate database: %s", err)
	}

	authService := services.NewAuthService(db)
	auth.InitSuperTokens(authService)
	auth.SetupRoles()

	middlewares.RegisterCustomValidator()

	app := fiber.New()

	// Setup middlewares
	app.Use(logger.New())

	corsHeaders := []string{"content-type"}
	corsHeaders = append(corsHeaders, supertokens.GetAllCORSHeaders()...)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ORIGINS"),
		AllowHeaders:     strings.Join(corsHeaders, ","),
		AllowCredentials: true,
	}))

	app.Use(compress.New())
	app.Use(recover.New())

	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))

	setupRoute(app, db)

	// Start server
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatalln("Error listening")
		log.Fatalln(err)
	}
}

func setupRoute(app *fiber.App, db *gorm.DB) {
	user.Route(app, db)
	profile.Route(app, db)

	app.Get("/docs/*", swagger.HandlerDefault)
}
