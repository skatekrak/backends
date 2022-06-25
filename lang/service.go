package lang

import (
	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) FindAll() ([]model.Lang, error) {
	var langs []model.Lang
	err := s.db.Find(&langs).Error
	return langs, err
}

func (s *Service) Get(isoCode string) (model.Lang, error) {
	var lang model.Lang
	err := s.db.Unscoped().Where("iso_code = ?", isoCode).First(&lang).Error
	return lang, err
}

func (s *Service) Create(lang *model.Lang) error {
	return s.db.Create(&lang).Error
}

func (s *Service) Update(lang *model.Lang) error {
	return s.db.Save(&lang).Error
}

func (s *Service) Delete(lang *model.Lang) error {
	return s.db.Delete(&lang).Error
}
