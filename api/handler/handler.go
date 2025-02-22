package handler

import (
	"net/http"
	"pixi/routes"
	"pixi/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler function for Vercel
func Handler(db *gorm.DB) http.HandlerFunc {
	// Initialize Gin router
	rGin := gin.Default()

	// Set up CORS options
	corsOptions := cors.New(cors.Config{
		AllowOrigins:     []string{utils.GetEnv("FRONTEND_URL", "https://er8ooes-anonymous-8081.exp.direct")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	rGin.Use(corsOptions)

	// Set the db instance in the Gin context
	rGin.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Register routes
	routes.UserRoutes(rGin)
	routes.PostRoutes(rGin)
	routes.SavedRoutes(rGin)
	routes.CommentRoutes(rGin)
	routes.LikeRoutes(rGin)
	routes.FollowRoutes(rGin)
	routes.ReplyRoutes(rGin)
	routes.SchedulerRoutes(rGin)

	// Return the Gin router as an http.HandlerFunc
	return func(w http.ResponseWriter, r *http.Request) {
		rGin.ServeHTTP(w, r)
	}
}
