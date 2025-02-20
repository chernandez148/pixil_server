package models

import (
	"time"
)

type Comment struct {
	ID           uint      `gorm:"primaryKey"`
	Author       string    `gorm:"not null"`
	Content      string    `gorm:"not null"`
	AudioContent string    `gorm:"default:null"`
	Replies      []Reply   `gorm:"foreignKey:CommentID"` // Correct the foreign key
	CreatedAt    time.Time `gorm:"autoCreateTime"`       // Automatically set timestamp
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`       // Automatically update timestamp
	PostID       uint      `gorm:"not null"`
	UserID       uint      `gorm:"not null"`
	User         *User     `gorm:"foreignKey:UserID"`
}
