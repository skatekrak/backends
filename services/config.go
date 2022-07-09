package services

import (
	"database/sql"
	"errors"

	"github.com/skatekrak/scribe/model"
	"gorm.io/gorm"
)

type ConfigKey = string

const (
	FeedlyToken ConfigKey = "feedly_token"
)

var keys = [...]ConfigKey{FeedlyToken}

type ConfigService struct {
	db *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{db}
}

// Add db entry if needed for each of the config keys
func (s *ConfigService) InitSetup() error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		batches := []model.Config{}

		for _, k := range keys {
			// Only add keys that are not already here
			if _, err := s.Get(k); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					batches = append(batches, model.Config{
						Key: k,
					})
				}
			}
		}

		return tx.CreateInBatches(batches, len(batches)).Error
	})
}

func (s *ConfigService) Get(key ConfigKey) (sql.NullString, error) {
	var configKey *model.Config

	if err := s.db.Where("key = ?", key).First(&configKey).Error; err != nil {
		return sql.NullString{}, err
	}

	return configKey.Value, nil
}

func (s *ConfigService) Set(key ConfigKey, value *string) error {

	v := sql.NullString{}

	if value != nil {
		v = sql.NullString{
			String: *value,
			Valid:  true,
		}
	}

	return s.db.Create(&model.Config{
		Key:   key,
		Value: v,
	}).Error
}
