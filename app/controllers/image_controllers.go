package controllers

import (
	"net/http"
	"net/url"
	"pixi/app/config"
	"pixi/app/models"

	"github.com/gin-gonic/gin"
)

// UploadImage handles uploading an image and associating it with a post
func UploadImage(c *gin.Context) {
	var newImage models.Image

	// Get UserID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	// Assign the UserID from the context to the new image
	newImage.UserID = userID.(uint)

	// Bind the request to the newImage object
	if err := c.ShouldBindJSON(&newImage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate required fields
	if newImage.ImageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image URL is required"})
		return
	}

	// Validate the Image URL format
	if _, err := url.ParseRequestURI(newImage.ImageURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Image URL format"})
		return
	}

	// Save the new image to the database
	db := config.GetDB()
	if err := db.Create(&newImage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving image", "details": err.Error()})
		return
	}

	// Retrieve the saved image
	var createdImage models.Image
	if err := db.First(&createdImage, newImage.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving image with associations", "details": err.Error()})
		return
	}

	// Respond with the created image
	c.JSON(http.StatusCreated, gin.H{"message": "Image uploaded successfully", "image": createdImage})
}
