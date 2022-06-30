package content

import (
	"fmt"

	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) Find(sourceTypes []string, page int) ([]model.Content, error) {
	var contents []model.Content

	err := s.db.Find(&contents).
		Order("\"published_at desc\"").
		Joins(fmt.Sprintf("left join contents on sources.source_type in %s", sourceTypes)).
		Limit(50).
		Offset((page - 1) * 50).
		Error

	if err != nil {
		return contents, err
	}

	return contents, err
}

func (s *Service) FindFromSource(sourceID string, page int) ([]model.Content, error) {
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

func (s *Service) AddMany(sources []*model.Content) error {
	return s.db.CreateInBatches(sources, len(sources)).Error
}

func (s *Service) FindOneByContentID(contentID string) (model.Content, error) {
	var content model.Content
	err := s.db.Where("content_id = ?", contentID).First(&content).Error
	return content, err
}
