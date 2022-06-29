package model

import (
	"time"

	"gorm.io/gorm"
)

type Lang struct {
	IsoCode   string `gorm:"primaryKey" json:"isoCode"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ImageURL  string         `json:"imageURL"`

	Sources []Source
}

type Source struct {
	gorm.Model
	RefreshedAt *time.Time
	Order       int `gorm:"index"`
	SourceType  string
	LangIsoCode string
	FeedlyID    string
	Title       string
	ShortTitle  string
	IconURL     string
	CoverURL    string
	Description string
	SkateSource bool `gorm:"default:true"`
	WebsiteURL  string
	PublishedAt *time.Time
	SourceID    string `gorm:"unique,index"` // Vimeo, Youtube or Feedly ID, depending on the type
}
