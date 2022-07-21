package services

import (
	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

type LangService struct {
	db *gorm.DB
}

func NewLangService(db *gorm.DB) *LangService {
	return &LangService{db}
}

func (s *LangService) FindAll() ([]model.Lang, error) {
	var langs []model.Lang
	err := s.db.Find(&langs).Error
	return langs, err
}

func (s *LangService) Get(isoCode string) (model.Lang, error) {
	var lang model.Lang
	err := s.db.Unscoped().Where("iso_code = ?", isoCode).First(&lang).Error
	return lang, err
}

func (s *LangService) Create(lang *model.Lang) error {
	return s.db.Create(&lang).Error
}

func (s *LangService) Update(lang *model.Lang) error {
	return s.db.Save(&lang).Error
}

func (s *LangService) Delete(lang *model.Lang) error {
	return s.db.Delete(&lang).Error
}
