package services

import (
	"github.com/skatekrak/carrelage/models"
	"gorm.io/gorm"
)

type UsersService struct {
	db *gorm.DB
}

func NewUsersService(db *gorm.DB) *UsersService {
	return &UsersService{db}
}

func (s *UsersService) Get(id string) (*models.User, error) {
	var user models.User
	err := s.db.Where("id = ?", id).First(&user).Error
	return &user, err
}
