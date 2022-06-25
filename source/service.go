package source

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

func (s *Service) FindAll(types []string) ([]model.Source, error) {
	var sources []model.Source
	query := s.db

	if len(types) > 0 {
		query = query.Where("source_type IN ?", types)
	}

	err := query.Find(&sources).Error
	return sources, err
}

func (s *Service) Get(id string) (model.Source, error) {
	var source model.Source
	err := s.db.Where("id = ?", id).First(&source).Error
	return source, err
}

func (s *Service) GetBySourceID(sourceID string) (model.Source, error) {
	var source model.Source
	err := s.db.Where("source_id = ?", sourceID).First(&source).Error
	return source, err
}

func (s *Service) GetNextOrder() (int, error) {
	var sources []model.Source
	if err := s.db.Order("\"order desc\"").Limit(1).Find(&sources).Error; err != nil {
		return 0, err
	}

	if len(sources) > 0 {
		return sources[0].Order + 1, nil
	}
	return 0, nil
}

func (s *Service) Create(source *model.Source) error {
	return s.db.Create(&source).Error
}

func (s *Service) Update(source *model.Source) error {
	return s.db.Save(&source).Error
}

func (s *Service) Delete(source *model.Source) error {
	return s.db.Delete(&source).Error
}
