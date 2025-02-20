package routes

import (
	"pixi/app/controllers"
	"pixi/app/middleware"

	"github.com/gin-gonic/gin"
)

// FollowerRoutes registers the routes for followers
func FollowRoutes(r *gin.Engine) {
	// Group routes for followers
	followGroup := r.Group("/follow")

	// Use the authentication middleware for routes that require user authentication
	// followerGroup.Use(middleware.AuthRequired())

	r.GET("followers/:id", controllers.GetFollowersByUserID)

	r.GET("followings/:id", controllers.GetFollowingsByUserID)

	// Route to create a new follower
	followGroup.POST("", middleware.AuthRequired(), controllers.CreateAFollow)

	// Route to delete a follower by ID
	followGroup.DELETE("/:id", middleware.AuthRequired(), controllers.DeleteFollower)
}
