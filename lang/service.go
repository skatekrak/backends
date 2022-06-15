package lang

import "gorm.io/gorm"

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) FindAll() ([]Lang, error) {
	var langs []Lang
	err := s.db.Find(&langs).Error
	return langs, err
}

func (s *Service) Get(isoCode string) (Lang, error) {
	var lang Lang
	err := s.db.Where("iso_code = ?", isoCode).First(&lang).Error
	return lang, err
}

func (s *Service) Create(lang *Lang) error {
	return s.db.Create(&lang).Error
}

func (s *Service) Update(lang *Lang) error {
	return s.db.Save(&lang).Error
}

func (s *Service) Delete(lang *Lang) error {
	return s.db.Delete(&lang).Error
}
