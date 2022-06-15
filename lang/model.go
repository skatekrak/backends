package lang

import (
	"time"

	"gorm.io/gorm"
)

type Lang struct {
	IsoCode string `gorm:"primaryKey" json:"isoCode"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ImageURL string `json:"imageURL"`
}