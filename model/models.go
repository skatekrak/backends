package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Lang struct {
	IsoCode   string         `gorm:"primaryKey" json:"isoCode"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	ImageURL  string         `json:"imageUrl"`

	Sources []Source
}

type Source struct {
	Model

	RefreshedAt *time.Time `json:"refreshedAt"`
	Order       int        `gorm:"index" json:"order"`
	SourceType  string     `json:"sourceType"`
	LangIsoCode string     `json:"lang"`
	Title       string     `json:"title"`
	ShortTitle  string     `json:"shortTitle"`
	IconURL     string     `json:"iconUrl"`
	CoverURL    string     `json:"coverUrl"`
	Description string     `json:"description"`
	SkateSource bool       `gorm:"default:true" json:"skateSource"`
	WebsiteURL  string     `json:"websiteUrl"`
	PublishedAt *time.Time `json:"publishedAt"`
	SourceID    string     `gorm:"unique,index" json:"sourceId"` // Vimeo, Youtube or Feedly ID, depending on the type

	Contents []Content
}

type Content struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	SourceID     uint      `json:"sourceId"`
	ContentID    string    `gorm:"unique,index" json:"contentId"` // Youtube or Vimeo ID or Feedly ID
	PublishedAt  time.Time `json:"publishedAt"`
	Title        string    `json:"title"`
	ContentURL   string    `json:"contentUrl"` // Youtube or Vimeo video url or article URL
	ThumbnailURL string    `json:"thumbnailUrl"`
	RawSummary   string    `json:"rawSummary"`
	Summary      string    `json:"summary"`
	RawContent   string    `json:"rawContent"`
	Content      string    `json:"content"`
	Author       *string   `json:"author"` // For feedly article
	Type         string    `json:"type"`
}

func (c *Content) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()
	return
}
