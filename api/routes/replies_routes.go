package routes

import (
	"pixi/api/controllers"
	"pixi/api/middleware"

	"github.com/gin-gonic/gin"
)

// RepliesRoutes defines the routes for handling reply-related operations
func ReplyRoutes(router *gin.Engine) {
	replyGroup := router.Group("/replies")

	// Apply authentication middleware to all routes in this group
	// replyGroup.Use(middleware.AuthRequired())

	// POST route to create a new reply on a specific comment
	replyGroup.POST("", middleware.AuthRequired(), controllers.CreateReply)

	// GET route to retrieve all replies by comment ID
	router.GET("/:postID/replies", controllers.GetRepliesByPostID)

	// PATCH route to update an existing reply by ID
	router.PATCH("/reply/:replyID", controllers.UpdateReply)

	// DELETE route to remove a reply by ID
	router.DELETE("/reply/:replyID", controllers.DeleteReply)
}
