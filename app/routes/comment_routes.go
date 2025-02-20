package routes

import (
	"pixi/app/controllers"
	"pixi/app/middleware"

	"github.com/gin-gonic/gin"
)

// CommentRoutes defines the routes for handling comment-related operations.
func CommentRoutes(router *gin.Engine) {
	commentGroup := router.Group("/posts/:postID/comments")

	// POST route to create a new comment on a specific post
	commentGroup.POST("", middleware.AuthRequired(), controllers.CreateComment)

	// GET route to retrieve all comments for a specific post
	commentGroup.GET("", controllers.GetCommentsByPostID)

	// PATCH route to update an existing comment by ID
	router.PATCH("/comment/:commentID", controllers.UpdateComment)

	// DELETE route to remove a comment by ID
	router.DELETE("/comment/:commentID", controllers.UpdateComment)
}
