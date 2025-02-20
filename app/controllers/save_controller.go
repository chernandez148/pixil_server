package controllers

import (
	"net/http"
	"pixi/app/config"
	"pixi/app/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateSave handles the creation of a new save
func CreateSave(c *gin.Context) {
	var newSave models.Save

	// Get userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	// Bind the request body to newLike object
	if err := c.ShouldBindJSON(&newSave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Assign the UserID from the context to the new like
	newSave.UserID = userID.(uint)

	// Validate required fields
	if newSave.UserID == 0 || newSave.PostID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and Post ID are required"})
		return
	}

	// Get the DB connection
	db := config.GetDB()

	// Ensure the post exists
	var post models.Post
	if err := db.First(&post, newSave.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Prevent duplicate likes
	var existingSave models.Like
	if err := db.Where("user_id = ? AND post_id = ?", newSave.UserID, newSave.PostID).First(&existingSave).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Save already exists"})
		return
	}

	// Save the new like to the database
	if err := db.Create(&newSave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving post", "details": err.Error()})
		return
	}

	// Respond with the created like
	c.JSON(http.StatusCreated, gin.H{"message": "Like created successfully", "save": newSave})
}

// DeleteSave handles the deletion of a saved post
func DeleteSave(c *gin.Context) {
	// Get saveID from the URL parameter
	saveID := c.Param("saveID")

	// Get userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	// Validate saveID and userID
	if saveID == "" || userID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SaveID and UserID are required"})
		return
	}

	// Convert saveID to uint
	saveIDUint, err := strconv.ParseUint(saveID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SaveID"})
		return
	}

	// Get the DB connection
	db := config.GetDB()

	// Find the save by ID and UserID
	var save models.Save
	if err := db.Where("id = ? AND user_id = ?", saveIDUint, userID).First(&save).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Saved post not found"})
		return
	}

	// Delete the save
	if err := db.Delete(&save).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete saved post", "details": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Saved post deleted successfully"})
}
