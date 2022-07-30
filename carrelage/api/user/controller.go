package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/carrelage/services"
	"github.com/skatekrak/utils/middlewares"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
)

type Controller struct {
	usersService *services.UsersService
}

func (c *Controller) Get(ctx *fiber.Ctx) error {
	user := c.getUser(ctx)

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (c *Controller) GetMe(ctx *fiber.Ctx) error {
	sessionContainer := session.GetSessionFromRequestContext(ctx.UserContext())
	userID := sessionContainer.GetUserID()

	user, err := c.usersService.Get(userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Well, that's awkward, you don't exists",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (c *Controller) UpdateUser(ctx *fiber.Ctx) error {
	user := c.getUser(ctx)
	body := ctx.Locals(middlewares.BODY).(UpdateUserBody)

	errs := map[string]string{}
	for _, role := range body.Roles {
		resp, err := userroles.AddRoleToUser(user.ID, role, nil)
		if err != nil {
			errs[role] = err.Error()
		}

		if resp.UnknownRoleError != nil {
			errs[role] = "This role doesn't exists"
		}
	}

	if len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(errs)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "roles updated",
	})
}
