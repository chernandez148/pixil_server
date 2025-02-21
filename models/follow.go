package models

import (
	"time"
)

type Follow struct {
	ID          uint      `gorm:"primaryKey"`                           // Primary key for the follow record
	FollowerID  uint      `gorm:"index;not null"`                       // Foreign key for the follower
	FollowingID uint      `gorm:"index;not null"`                       // Foreign key for the followee
	Follower    User      `gorm:"foreignKey:FollowerID;references:ID"`  // Relationship to User via FollowerID
	Following   User      `gorm:"foreignKey:FollowingID;references:ID"` // Relationship to User via FolloweeID
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
