package controllers

import (
	"net/http"
	"pixi/config"
	"pixi/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRepliesByPostID retrieves all replies associated with a given comment ID
func GetRepliesByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	db := config.GetDB()
	var replies []models.Reply

	// Retrieve replies by post ID (assuming there's a `post_id` field)
	if err := db.Preload("Author").Where("post_id = ?", postID).Find(&replies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve replies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"replies": replies})
}

// CreateReply handles the creation of a new reply
func CreateReply(c *gin.Context) {
	var newReply models.Reply

	// Get userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	// Assign the UserID from the context to the new post
	newReply.UserID = userID.(uint)

	// Bind the request body to newReply object
	if err := c.ShouldBindJSON(&newReply); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate required fields
	if newReply.PostID == 0 || newReply.UserID == 0 || newReply.AuthorID == 0 || newReply.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID, User ID, Reply User ID, Comment ID, and Content are required"})
		return
	}

	// Get the DB connection
	db := config.GetDB()

	// Ensure the post, user, and comment exist before creating a reply
	var post models.Post
	var user models.User
	var comment models.Comment

	if err := db.First(&post, newReply.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := db.First(&user, newReply.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.First(&comment, newReply.CommentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Save the new reply to the database
	if err := db.Create(&newReply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving reply", "details": err.Error()})
		return
	}

	// Respond with the created reply
	c.JSON(http.StatusCreated, gin.H{"message": "Reply created successfully", "reply": newReply})
}

// UpdateReply updates the reply
func UpdateReply(c *gin.Context) {
	db := config.GetDB()
	replyID := c.Param("id")

	// Bind the input JSON data to a map to handle specific fields
	var updatedReply models.Reply

	// Bind incoming JSON request data to updatedReply object
	if err := c.ShouldBindJSON(&updatedReply); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Find the reply by ID
	var existingReply models.Reply
	if err := db.First(&existingReply, replyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reply not found"})
		return
	}

	// Update the reply with the new data
	if updatedReply.Content != "" {
		existingReply.Content = updatedReply.Content
	}

	// Save the updated reply to the database
	if err := db.Save(&existingReply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reply", "details": err.Error()})
		return
	}

	// Respond with the updated reply
	c.JSON(http.StatusOK, gin.H{"message": "Reply updated successfully", "reply": existingReply})
}

// Delete a reply
func DeleteReply(c *gin.Context) {
	db := config.GetDB()
	var reply models.Reply
	replyID := c.Param("id")

	// Find the reply by ID
	if err := db.First(&reply, replyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reply not found"})
		return
	}

	// Delete the reply
	if err := db.Delete(&reply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reply", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reply deleted successfully"})
}
