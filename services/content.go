package services

import (
	"github.com/skatekrak/scribe/database"
	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

type ContentService struct {
	db *gorm.DB
}

func NewContentService(db *gorm.DB) *ContentService {
	return &ContentService{db}
}

func (s *ContentService) Find(sourceTypes []string, sources []int, page int) (*database.Pagination, error) {
	pagination := &database.Pagination{
		PerPage: 50,
		Page:    page,
		Items:   []model.Content{},
	}

	tx := s.db.Model(pagination.Items).
		Order("contents.created_at desc").
		Joins("Source").
		Session(&gorm.Session{})

	if len(sourceTypes)+len(sources) > 0 {
		tx = tx.Joins("JOIN sources ON sources.id = contents.source_id")
	}

	if len(sourceTypes) > 0 {
		tx = tx.
			Where("sources.source_type in ?", sourceTypes).
			Session(&gorm.Session{})
	}

	if len(sources) > 0 {
		tx = tx.Where("sources.id in ?", sources)
	}

	tx = tx.
		Scopes(pagination.Scope()).
		Find(&pagination.Items)

	return pagination, tx.Error
}

func (s *ContentService) Get(id string) (model.Content, error) {
	var content model.Content
	err := s.db.Where("contents.id = ?", id).Joins("Source").Preload("Source.Lang").First(&content).Error
	return content, err
}

func (s *ContentService) AddMany(contents []*model.Content, sources []*model.Source) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(contents, len(contents)).Error; err != nil {
			return err
		}

		for i := range sources {
			if err := tx.Save(&sources[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ContentService) FindOneByContentID(contentID string) (model.Content, error) {
	var content model.Content
	err := s.db.Where("content_id = ?", contentID).First(&content).Error
	return content, err
}
