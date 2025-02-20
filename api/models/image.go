package models

import (
	"time"
)

type Image struct {
	ID        uint      `gorm:"primaryKey"` // Unique identifier for the image
	ImageURL  string    `gorm:"not null"`   // URL of the image
	CreatedAt time.Time // Timestamp when the image was created
	UpdatedAt time.Time // Timestamp when the image was last updated
	UserID    uint      `gorm:"not null;index"`    // ID of the user who uploaded the image
	User      *User     `gorm:"foreignKey:UserID"` // User who uploaded the image
	PostID    uint      `gorm:"not null;index"`    // ID of the post this image belongs to
	Post      *Post     `gorm:"foreignKey:PostID"` // Post this image belongs to
}
