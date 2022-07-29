package services

import (
	"github.com/skatekrak/carrelage/models"
	"gorm.io/gorm"
)

type ProfilesService struct {
	db *gorm.DB
}

func NewProfilesService(db *gorm.DB) *ProfilesService {
	return &ProfilesService{db}
}

// Get profile based on given id
func (s *ProfilesService) Get(id string) (*models.Profile, error) {
	var profile models.Profile
	err := s.db.Where("id = ?", id).First(&profile).Error
	return &profile, err
}

func (s *ProfilesService) GetFromUserID(id string) (*models.Profile, error) {
	var profile models.Profile

	err := s.db.Where("user_id = ?", id).First(&profile).Error
	return &profile, err
}

func (s *ProfilesService) Update(profile *models.Profile) error {
	return s.db.Save(&profile).Error
}

func (s *ProfilesService) IsUsernameAvailable(username string) (bool, error) {
	var totalResults int64
	err := s.db.Model(&models.Profile{}).Where("username = ?", username).Count(&totalResults).Error

	return totalResults <= 0, err
}
