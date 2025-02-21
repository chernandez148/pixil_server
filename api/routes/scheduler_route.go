package routes

import (
	"net/http"
	"pixi/api/controllers"

	"github.com/gin-gonic/gin"
)

func SchedulerRoutes(r *gin.Engine) {
	r.GET("/trigger-scheduler", func(c *gin.Context) {
		controllers.PublishScheduledPosts() // Manually trigger the function
		c.JSON(http.StatusOK, gin.H{"message": "Scheduler triggered successfully"})
	})
}
