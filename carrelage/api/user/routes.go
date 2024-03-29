package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/carrelage/auth"
	"github.com/skatekrak/carrelage/auth/roles"
	"github.com/skatekrak/carrelage/services"
	"gorm.io/gorm"
)

type UpdateUserBody struct {
	Roles []string `json:"roles" validate:"oneof=superadmin,admin,moderator,user"`
}

func Route(app *fiber.App, db *gorm.DB) {
	usersService := services.NewUsersService(db)
	controller := &Controller{
		usersService: usersService,
	}

	router := app.Group("users")

	router.Get("/me", auth.Logged(nil), controller.GetMe)
	router.Get("/:userID", auth.Logged(nil), auth.RequireRole(roles.MODERATOR), controller.Loader(), controller.Get)
	router.Patch("/:userID", auth.Logged(nil), auth.RequireRole(roles.ADMIN), controller.Loader(), controller.UpdateUser)
}
