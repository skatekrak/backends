package services

import (
	"log"

	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SourceService struct {
	db *gorm.DB
}

func NewSourceService(db *gorm.DB) *SourceService {
	return &SourceService{db}
}

func (s *SourceService) FindAll(types []string) ([]*model.Source, error) {
	var sources []*model.Source
	query := s.db.Session(&gorm.Session{})

	if len(types) > 0 {
		query = query.Where("source_type IN ?", types)
	}

	err := query.Find(&sources).Error
	return sources, err
}

func (s *SourceService) Get(id string) (model.Source, error) {
	var source model.Source
	err := s.db.Where("id = ?", id).First(&source).Error
	return source, err
}

func (s *SourceService) GetBySourceID(sourceID string) (model.Source, error) {
	var source model.Source
	err := s.db.Where("source_id = ?", sourceID).First(&source).Error
	return source, err
}

func (s *SourceService) GetNextOrder() (int, error) {
	var sources []model.Source
	if err := s.db.Order("\"order\" desc").Limit(1).Find(&sources).Error; err != nil {
		return 0, err
	}

	if len(sources) > 0 {
		return sources[0].Order + 1, nil
	}

	return 0, nil
}

func (s *SourceService) Create(source *model.Source) error {
	return s.db.Create(&source).Error
}

func (s *SourceService) Update(source *model.Source) error {
	return s.db.Save(&source).Error
}

func (s *SourceService) Delete(source *model.Source) error {
	return s.db.Unscoped().Delete(&source).Error
}

func (s *SourceService) AddMany(sources []*model.Source) error {
	return s.db.CreateInBatches(sources, len(sources)).Error
}

func (s *SourceService) UpdateOrder(updates map[int]map[string]interface{}) ([]model.Source, error) {
	var sources []model.Source
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for id, update := range updates {
			log.Println("update", id, update)
			var s []model.Source
			if err := tx.Model(&s).Clauses(clause.Returning{}).Where("id = ?", id).Updates(update).Error; err != nil {
				return err
			}
			sources = append(sources, s...)
		}

		return nil
	})

	return sources, err
}
