package models

import (
	"time"
)

type Post struct {
	ID             uint      `gorm:"primaryKey"` // Unique identifier for the post
	Caption        string    `gorm:"not null"`   // Caption for the post
	Description    string    `gorm:"default:''"` // Optional description for the post
	ImageURL       string    `gorm:"not null"`   // URL of the post's image
	ScheduledAt    time.Time `json:"ScheduledAt"`
	CreatedAt      time.Time // Timestamp when the post was created
	UpdatedAt      time.Time // Timestamp when the post was last updated
	UserID         uint      `gorm:"not null;index"`      // ID of the user who created the post
	User           *User     `gorm:"foreignKey:UserID"`   // User who created the post
	AllowComments  bool      `gorm:"default:false"`       // Whether comments are allowed on the post
	HideLikeCounts bool      `gorm:"default:false"`       // Whether to hide like counts on the post
	IsPrivate      bool      `gorm:"default:false"`       // Whether the post is private
	IsScheduled    bool      `gorm:"default:false"`       // Whether the post is scheduled
	Status         string    `gorm:"default:'scheduled'"` // Add a Status field to track the post's status (draft, published, etc.)
	Comments       []Comment `gorm:"foreignKey:PostID"`   // List of comments on the post
	Likes          []Like    `gorm:"foreignKey:PostID"`   // List of likes on the post
}
