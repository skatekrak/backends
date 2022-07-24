package services

import (
	"github.com/lib/pq"
	"github.com/skatekrak/carrelage/models"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db}
}

func (s *AuthService) CreateUserAndProfileIfNotExists(id string) error {
	user := models.User{
		Model: models.Model{
			ID: id,
		},
	}

	profile := models.Profile{
		User: user,
	}

	err := s.db.Create(&profile).Error
	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == "unique_violation" {
		return nil
	}
	return err
}
