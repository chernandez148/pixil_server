package controllers

import (
	"net/http"
	"pixi/api/config"
	"pixi/api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get comments by post id
func GetCommentsByPostID(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil || postID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Declare a slice to hold comments
	db := config.GetDB()
	var comments []models.Comment

	// Retrieve all comments for the given post from the database
	if err := db.Preload("User").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	// Return the comments as JSON
	c.JSON(http.StatusOK, comments)
}

// CreateComment handles the creation of a new comment
func CreateComment(c *gin.Context) {
	var newComment models.Comment

	// Get userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	// Assign the UserID from the context to the new post
	newComment.UserID = userID.(uint)

	// Bind the request body to newComment object
	if err := c.ShouldBindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate required fields
	if newComment.UserID == 0 || newComment.PostID == 0 || newComment.Author == "" || newComment.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID, Post ID, and Content are required"})
		return
	}

	// Get the DB connection
	db := config.GetDB()

	// Ensure the post exists before creating a comment
	var post models.Post
	if err := db.First(&post, newComment.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Save the new comment to the database
	if err := db.Create(&newComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving comment", "details": err.Error()})
		return
	}

	// Respond with the created comment
	c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully", "comment": newComment})
}

// UpdateComment updates an existing comment
func UpdateComment(c *gin.Context) {
	db := config.GetDB()
	commentID := c.Param("id")

	// Bind the input JSON data to a map to handle specific fields
	var input struct {
		Author  *string `json:"Author"` // Nullable Author
		Content string  `json:"Content"`
	}

	// Bind incoming JSON request data
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Find the comment by ID
	var existingComment models.Comment
	if err := db.First(&existingComment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Update the fields if they are provided in the request
	if input.Content != "" {
		existingComment.Content = input.Content
	}

	// If an Author is provided, update it; otherwise, leave it as it is
	if input.Author != nil {
		existingComment.Author = *input.Author
	}

	// Save the updated comment to the database
	if err := db.Save(&existingComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating comment", "details": err.Error()})
		return
	}

	// Respond with the updated comment
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"comment": existingComment,
	})
}

// DeleteComment deletes a comment by ID
func DeleteComment(c *gin.Context) {
	db := config.GetDB()
	var comment models.Comment
	commentID := c.Param("id")

	// Find the comment by ID
	if err := db.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Delete the comment
	if err := db.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting comment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
