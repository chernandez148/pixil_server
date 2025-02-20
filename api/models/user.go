package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt" // Import bcrypt package
	"gorm.io/gorm"
)

// User represents a user record in the database
type User struct {
	ID           uint      `gorm:"primaryKey"`
	FullName     string    `gorm:"not null"`
	Username     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	Password     string    `gorm:"not null"`
	Bio          string    `gorm:"default:''"`
	ProfileImage string    `gorm:"default:''"`
	Posts        []Post    `gorm:"foreignKey:UserID"`
	Saves        []Save    `gorm:"foreignKey:UserID"`
	Likes        []Like    `gorm:"foreignKey:UserID"`
	Followees    []Follow  `gorm:"foreignKey:FollowingID"` // Use Follow as the model here
	Followings   []Follow  `gorm:"foreignKey:FollowerID"`  // Use Follow as the model here
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// HashPassword hashes the user's password before storing it
func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// BeforeCreate is a GORM hook that gets called before a user record is created
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Hash the password before saving the user
	err = user.HashPassword()
	if err != nil {
		return err
	}
	return nil
}

// Example method that uses gorm to save the user
func (user *User) Save(db *gorm.DB) error {
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// CheckPassword compares the provided password with the stored hashed password
func (user *User) CheckPassword(providedPassword string) error {
	// Compare the provided password with the stored hashed password using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}
