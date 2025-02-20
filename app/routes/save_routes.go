package routes

import (
	"pixi/app/controllers"
	"pixi/app/middleware"

	"github.com/gin-gonic/gin"
)

// LikeRoutes defines the routes for handeling like-related operations
func SavedRoutes(router *gin.Engine) {
	saveGroup := router.Group("/posts/:postID/saves")

	// POST route to create a new like on a specific post
	saveGroup.POST("", middleware.AuthRequired(), controllers.CreateSave)

	// DELETE route to remove a like by ID
	router.DELETE("/saves/:saveID", middleware.AuthRequired(), controllers.DeleteSave)
}
