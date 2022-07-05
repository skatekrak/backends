package model

import (
	"time"

	"github.com/google/uuid"
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
	Title       string
	ShortTitle  string
	IconURL     string
	CoverURL    string
	Description string
	SkateSource bool `gorm:"default:true"`
	WebsiteURL  string
	PublishedAt *time.Time
	SourceID    string `gorm:"unique,index"` // Vimeo, Youtube or Feedly ID, depending on the type

	Contents []Content
}

type Content struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	SourceID     uint
	ContentID    string `gorm:"unique,index"` // Youtube or Vimeo ID or Feedly ID
	PublishedAt  time.Time
	Title        string
	ContentURL   string // Youtube or Vimeo video url or article URL
	ThumbnailURL string
	RawSummary   string
	Summary      string
	RawContent   string
	Content      string
	Author       *string // For feedly article
	Type         string
}

func (c *Content) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()
	return
}
