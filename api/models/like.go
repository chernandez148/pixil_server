package models

import (
	"time"
)

type Like struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	PostID    uint      `gorm:"not null"`
	User      *User     `gorm:"foreignKey:UserID"`
	Post      *Post     `gorm:"foreignKey:PostID"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
