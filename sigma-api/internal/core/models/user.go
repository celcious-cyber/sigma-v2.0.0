package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the system user (Admin, Guru, etc.)
type User struct {
	gorm.Model
	Name            string     `gorm:"type:string;not null" json:"name"`
	Email           string     `gorm:"type:string;uniqueIndex;not null" json:"email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	Password        string     `gorm:"type:string;not null" json:"-"` // Hidden from JSON
	Role            string     `gorm:"type:text" json:"role"`         // Stored as JSON string
	LastSeenAt      *time.Time `json:"last_seen_at,omitempty"`
	AvatarURL       *string    `json:"avatar_url,omitempty"`
	RememberToken   string     `gorm:"type:string" json:"-"`
}
