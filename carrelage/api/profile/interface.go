package profile

import (
	"time"

	"github.com/skatekrak/carrelage/models"
)

type GetProfileResponse struct {
	ID                string    `json:"id"`
	CreatedAt         time.Time `json:"createdAt"`
	Username          string    `json:"username"`
	ProfilePictureURL string    `json:"profilePictureURL,omitempty"`
	Bio               string    `json:"bio,omitempty"`
	Stance            string    `json:"stance,omitempty"`
}

func GetProfileResponseFrom(profile *models.Profile) *GetProfileResponse {
	response := &GetProfileResponse{
		ID:                profile.ID,
		CreatedAt:         profile.CreatedAt,
		Username:          profile.Username,
		ProfilePictureURL: profile.ProfilePictureURL,
		Bio:               profile.Bio,
		Stance:            profile.Stance,
	}

	return response
}
