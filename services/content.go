package services

import (
	"fmt"

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
	var contents []model.Content

	tx := s.db.Model(contents).Order("contents.published_at desc").Session(&gorm.Session{})

	if len(sourceTypes)+len(sources) > 0 {
		tx = tx.
			Joins("JOIN sources ON sources.id = contents.source_id")
	}

	if len(sourceTypes) > 0 {
		tx = tx.
			Where("sources.source_type in ?", sourceTypes).
			Session(&gorm.Session{})
	}

	if len(sources) > 0 {
		tx = tx.Where("sources.id in ?", sources)
	}

	pagination := &database.Pagination{
		PerPage: 50,
		Page:    page,
	}

	tx = tx.
		Scopes(pagination.Scope()).
		Find(&contents)

	if tx.Error != nil {
		return pagination, tx.Error
	}

	pagination.Items = contents

	return pagination, nil
}

func (s *ContentService) FindFromSource(sourceID string, page int) (*database.Pagination, error) {
	var contents []model.Content

	pagination := &database.Pagination{
		PerPage: 50,
		Page:    page,
	}

	err := s.db.Find(&contents).
		Order("\"published_at desc\"").
		Joins(fmt.Sprintf("left join contents on sources.id = %s", sourceID)).
		Scopes(pagination.Scope()).
		Error

	pagination.Items = contents

	return pagination, err
}

func (s *ContentService) AddMany(contents []*model.Content, sources []*model.Source) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(contents, len(contents)).Error; err != nil {
			return err
		}

		for _, s := range sources {
			if err := tx.Save(&s).Error; err != nil {
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
