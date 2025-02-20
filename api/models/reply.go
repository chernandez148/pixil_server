package models

import (
	"time"
)

type Reply struct {
	ID           uint      `gorm:"primaryKey"`
	PostID       uint      `gorm:"not null"`
	UserID       uint      `gorm:"not null"`
	User         *User     `gorm:"foreignKey:UserID"`
	AuthorID     uint      `gorm:"not null"`
	Author       *User     `gorm:"foreignKey:AuthorID"`
	CommentID    *uint     `gorm:"default:null"` // Nullable, links reply to a specific comment
	Content      string    `gorm:"not null"`
	AudioContent string    `gorm:"default:null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"` // Automatically set timestamp
	UpdatedAt    time.Time `gorm:"autoUpdateTime"` // Automatically update timestamp
}
