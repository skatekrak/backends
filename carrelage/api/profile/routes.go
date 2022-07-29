package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/carrelage/auth"
	"github.com/skatekrak/carrelage/services"
	"github.com/skatekrak/utils/middlewares"
	"gorm.io/gorm"
)

type UpdateProfileBody struct {
	Username *string `json:"username" validate:"omitempty,username"`
	Bio      *string `json:"bio"`
	Stance   *string `json:"stance" validate:"omitempty,oneof=regular goofy"`
}

func Route(app *fiber.App, db *gorm.DB) {
	profilesService := services.NewProfilesService(db)
	controller := &Controller{
		profilesService: profilesService,
	}

	router := app.Group("profiles")

	router.Get("/me", auth.Logged(nil), controller.GetMe)
	router.Get("/:profileID", controller.Loader(), controller.Get)
	router.Patch(
		"/:profileID",
		auth.Logged(nil),
		controller.Loader(),
		controller.IsProfileOwner(),
		middlewares.JSONHandler[UpdateProfileBody](),
		controller.Update,
	)
}
