package routes

import (
	"pixi/app/controllers"
	"pixi/app/middleware"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.Engine) {
	postGroup := router.Group("/posts")

	// Add GET routes for retrieving posts
	postGroup.GET("", controllers.GetPost)

	// Add POST route for creating a new post
	postGroup.POST("", middleware.AuthRequired(), controllers.CreatePost)

	// Add a scheduled post
	postGroup.POST("/create-scheduled-post", middleware.AuthRequired(), controllers.CreateScheduledPost)

	// Add PATCH route for updating an existing post
	router.PATCH("/post/:id", controllers.UpdatePost)

	// Add GET route for retrieving a post by ID
	router.GET("/post/:id", controllers.GetPostByID)

	// Add DELETE route for deleting a post by ID
	postGroup.DELETE("/:id", controllers.DeletePost)
}
