package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/skatekrak/carrelage/services"
	"github.com/skatekrak/utils/helpers"
	"github.com/skatekrak/utils/middlewares"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

type Controller struct {
	profilesService *services.ProfilesService
}

// @Summary  Find a profile with its profileID
// @Tags     profiles
// @Success  200        {object}  profile.GetProfileResponse
// @Failure  404        {object}  api.JSONError
// @Failure  500        {object}  api.JSONError
// @Param    profileID  path      string  true  "Profile ID"
// @Router   /profiles/{profileID} [get]
func (c *Controller) Get(ctx *fiber.Ctx) error {
	profile := c.getProfile(ctx)

	response := &GetProfileResponse{
		ID:                profile.ID,
		CreatedAt:         profile.CreatedAt,
		Username:          profile.Username,
		ProfilePictureURL: profile.ProfilePictureURL,
		Bio:               profile.Bio,
		Stance:            profile.Stance,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// @Summary  Get the profile of the current authenticated user
// @Tags     profiles
// @Success  200  {object}  profile.GetProfileResponse
// @Failure  404  {object}  api.JSONError
// @Failure  500  {object}  api.JSONError
// @Router   /profiles/me [get]
func (c *Controller) GetMe(ctx *fiber.Ctx) error {
	sessionContainer := session.GetSessionFromRequestContext(ctx.UserContext())
	userID := sessionContainer.GetUserID()

	profile, err := c.profilesService.GetFromUserID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := GetProfileResponseFrom(profile)

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// @Summary  Update profile
// @Tags     profiles
// @Success  200        {object}  profile.GetProfileResponse
// @Failure  404        {object}  api.JSONError
// @Failure  500        {object}  api.JSONError
// @Param    body       body      profile.UpdateProfileBody  true  "Update body"
// @Param    profileID  path      string                     true  "Profile ID"
// @Router   /profiles/{profileID} [patch]
func (c *Controller) Update(ctx *fiber.Ctx) error {
	profile := c.getProfile(ctx)
	body := ctx.Locals(middlewares.BODY).(UpdateProfileBody)

	if body.Username != nil {
		available, err := c.profilesService.IsUsernameAvailable(*body.Username)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Couldn't check if username is available",
				"error":   err.Error(),
			})
		}

		if !available {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Username already taken",
			})
		}
	}

	profile.Username = helpers.SetIfNotNil(body.Username, profile.Username)
	profile.Bio = helpers.SetIfNotNil(body.Bio, profile.Bio)
	profile.Stance = helpers.SetIfNotNil(body.Stance, profile.Stance)

	if err := c.profilesService.Update(profile); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Couldn't save profile",
			"error":   err.Error(),
		})
	}

	response := GetProfileResponseFrom(profile)

	return ctx.Status(fiber.StatusOK).JSON(response)
}
