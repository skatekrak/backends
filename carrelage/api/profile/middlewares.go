package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/utils/helpers"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
)

const PROFILE_LOADER = "profileID"

/**
Fetch the profile pass via :profileID in the URL.
Return 404 if not found, go to next middleware if found
**/
func (c *Controller) Loader() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID := ctx.Params(PROFILE_LOADER)

		user, err := c.profilesService.Get(userID)
		if err != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}

		ctx.Locals(PROFILE_LOADER, user)
		return ctx.Next()
	}
}

/**
Middleware that checks if connected is the owner of this profile.
If the user a moderator or above, it bypass it
**/
func (c *Controller) IsProfileOwner() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		profile := c.getProfile(ctx)

		sessionContainer := session.GetSessionFromRequestContext(ctx.UserContext())
		resp, err := userroles.GetRolesForUser(sessionContainer.GetUserID(), nil)
		if err != nil {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Couldn't get your roles",
			})
		}

		roles := resp.OK.Roles
		if helpers.Has(roles, "moderator") {
			ctx.Next()
		}

		if profile.UserID != sessionContainer.GetUserID() {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Can't go through, this isn't you",
			})
		}

		return ctx.Next()
	}
}
