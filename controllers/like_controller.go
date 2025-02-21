package controllers

import (
	"net/http"
	"pixi/config"
	"pixi/models"

	"github.com/gin-gonic/gin"
)

// CreateLike handles the creation of a new like
func CreateLike(c *gin.Context) {
	var newLike models.Like

	// Get userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	// Bind the request body to newLike object
	if err := c.ShouldBindJSON(&newLike); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Assign the UserID from the context to the new like
	newLike.UserID = userID.(uint)

	// Validate required fields
	if newLike.UserID == 0 || newLike.PostID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and Post ID are required"})
		return
	}

	// Get the DB connection
	db := config.GetDB()

	// Ensure the post exists
	var post models.Post
	if err := db.First(&post, newLike.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Prevent duplicate likes
	var existingLike models.Like
	if err := db.Where("user_id = ? AND post_id = ?", newLike.UserID, newLike.PostID).First(&existingLike).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Like already exists"})
		return
	}

	// Save the new like to the database
	if err := db.Create(&newLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving like", "details": err.Error()})
		return
	}

	// Respond with the created like
	c.JSON(http.StatusCreated, gin.H{"message": "Like created successfully", "like": newLike})
}

// DeleteLike handles the deletion of a like
func DeleteLike(c *gin.Context) {
	// Get PostID and UserID from query parameters or request body
	postID := c.Param("postID")
	userID := c.Param("userID")

	if postID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PostID and UserID are required"})
		return
	}

	// Get the DB connection
	db := config.GetDB()

	// Find the like by PostID and UserID
	var like models.Like
	if err := db.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		return
	}

	// Delete the like
	if err := db.Delete(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete like", "details": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Like deleted successfully"})
}
