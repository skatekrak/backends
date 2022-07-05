package services

import (
	"fmt"
	"log"

	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

type ContentService struct {
	db *gorm.DB
}

func NewContentService(db *gorm.DB) *ContentService {
	return &ContentService{db}
}

func (s *ContentService) Find(sourceTypes []string, page int) ([]model.Content, error) {
	var contents []model.Content

	log.Println("sourceTypes", sourceTypes)

	tx := s.db.Order("\"published_at desc\"").Session(&gorm.Session{})

	if len(sourceTypes) > 0 {
		tx = tx.
			Joins("JOIN sources ON sources.id = contents.source_id").
			Where("sources.source_type in ?", sourceTypes).
			Session(&gorm.Session{})
	}

	tx = tx.Limit(50).
		Offset((page - 1) * 50).
		Find(&contents)

	if tx.Error != nil {
		return contents, tx.Error
	}

	return contents, nil
}

func (s *ContentService) FindFromSource(sourceID string, page int) ([]model.Content, error) {
	var contents []model.Content

	err := s.db.Find(&contents).
		Order("\"published_at desc\"").
		Joins(fmt.Sprintf("left join contents on sources.id = %s", sourceID)).
		Limit(50).
		Offset((page - 1) * 50).
		Error

	if err != nil {
		return contents, err
	}

	return contents, err
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

	// return s.db.CreateInBatches(contents, len(contents)).Error
}

func (s *ContentService) FindOneByContentID(contentID string) (model.Content, error) {
	var content model.Content
	err := s.db.Where("content_id = ?", contentID).First(&content).Error
	return content, err
}
