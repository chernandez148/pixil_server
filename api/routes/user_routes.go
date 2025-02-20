// routes/user_routes.go
package routes

import (
	"pixi/api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")

	// Add POST route for user login
	router.POST("/login", controllers.Login)

	// Add POST route for creating a new user
	router.POST("/register", controllers.CreateUser)

	userGroup.GET("", controllers.GetUsers)

	// Add GET route for retrieving a user by ID
	userGroup.GET("/:id", controllers.GetUser) // Assuming you pass the user ID as a URL parameter

	// Add PATCH route for updating a user
	userGroup.PATCH("/:id", controllers.UpdateUser) // Assuming you pass the user ID as a URL parameter

	// Add DELETE route for deleting a user by ID
	userGroup.DELETE("/:id", controllers.DeleteUser) // Assuming you pass the user ID as a URL parameter
}
