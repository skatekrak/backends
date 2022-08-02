package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" swaggertype:"string"`
}

type User struct {
	Model

	UserSubscription *UserSubscription `json:"subscription,omitempty"`
	// Profile          Profile          `json:"profile"`
}

type UserSubscription struct {
	Model

	UserID               string     `json:"user"`
	SubscriptionStatus   string     `json:"subscriptionStatus"`
	SubscriptionEndAt    *time.Time `json:"subscriptionEndAt"`
	SubscriptionStripeId string     `json:"subscriptionStripeId"`
}

type Profile struct {
	Model

	UserID            string `json:"user"`
	User              User   `json:"-"`
	Username          string `gorm:"uniqueIndex" json:"username"`
	ProfilePictureURL string `json:"profilePictureUrl"`
	Bio               string `json:"bio"`
	Stance            string `json:"stance"`
}
