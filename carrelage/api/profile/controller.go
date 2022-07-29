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
