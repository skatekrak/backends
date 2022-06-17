package source

import (
	"time"

	"gorm.io/gorm"
)

type Source struct {
	gorm.Model
	RefreshedAt time.Time
	Order       int `gorm:"index"`
	SourceType  string
	LangIsoCode string
	FeedlyID    string
	Title       string
	ShortTitle  string
	IconURL     string
	CoverURL    string
	Description string
	SkateSource bool
	WebsiteURL  string
	PublishedAt time.Time
	SourceID    string `gorm:"unique,index"`
}
