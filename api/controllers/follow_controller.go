package controllers

import (
	"errors"
	"log"
	"net/http"
	"pixi/api/config"
	"pixi/api/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetFollowingsByUserID(c *gin.Context) {
	db := config.GetDB()
	var followers []models.Follow
	userID := c.Param("id") // Get the userID from the URL parameter

	// Query the Follow table and filter by FolloweeID
	if err := db.Preload("Following").Omit("Follower").Where("follower_id = ?", userID).Find(&followers).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follow data not found"})
		return
	}

	c.JSON(http.StatusOK, followers)
}

func GetFollowersByUserID(c *gin.Context) {
	db := config.GetDB()
	var followings []models.Follow
	userID := c.Param("id")

	if err := db.Preload("Follower").Omit("Following").Where("following_id = ?", userID).Find(&followings).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follow data not found"})
		return
	}

	c.JSON(http.StatusOK, followings)
}

func CreateAFollow(c *gin.Context) {
	var newFollow models.Follow

	// Get the logged in user id from the context (the follower)
	followerID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Assignes the user id from the context as the FollowerID
	newFollow.FollowerID = followerID.(uint)

	// Binds the request body to newFollow object
	if err := c.ShouldBindJSON(&newFollow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate required fields
	if newFollow.FollowingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Following ID is required"})
		return
	}

	// Check if a user is trying to follow themselves
	if newFollow.FollowerID == newFollow.FollowingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A user cannot follow themselves"})
		return
	}

	// Get the DB connection
	db := config.GetDB()

	// Ensure Follower exists (login in user)
	var follower models.User
	if err := db.First(&follower, newFollow.FollowerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follower not found"})
		return
	}

	// Ensure Follower exists (login in user)
	var following models.User
	if err := db.First(&following, newFollow.FollowingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User to follow not found"})
		return
	}

	// Check if the relationship already exists
	var existingFollow models.Follow
	err := db.Where("follower_id = ? AND following_id = ?", newFollow.FollowerID, newFollow.FollowingID).First(&existingFollow).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Unexpected error while checking the relationship
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking follow relationship", "details": err.Error()})
		return
	}

	// If relationship already exists, return a conflict
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Follow relationship already exists"})
		return
	}

	// No existing relationship found, so we proceed to create the new follow relationship
	if err := db.Create(&newFollow).Error; err != nil {
		log.Printf("Failed to create follow relationship: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create follower", "details": err.Error()})
		return
	}

	var createdFollow models.Follow
	if err := db.Preload("Following").First(&createdFollow, newFollow.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving follow with associations", "details": err.Error()})
		return
	}

	// Respond with the created follower
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Follow created successfully",
		"follower": createdFollow})

}

// DeleteFollower handles the deletion of a follower
func DeleteFollower(c *gin.Context) {
	// Validate follow ID parameter
	followIDParam := c.Param("id")
	followID, err := strconv.Atoi(followIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid follow ID format"})
		return
	}

	// Get userID from the context (set by authentication middleware)
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get the DB connection
	db := config.GetDB()

	// Find the follower relationship
	var follower models.Follow
	result := db.Where("follower_id = ? AND id = ?", followerID, followID).First(&follower)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Follower relationship not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": result.Error.Error()})
		}
		return
	}

	// Delete the follower relationship from the database
	if err := db.Delete(&follower).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete follower relationship", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Follower relationship deleted successfully"})
}
