package handler

import (
	"log"
	"net/http"
	"pixi/api/config"
	"pixi/api/controllers"
	"pixi/api/models"
	"pixi/api/routes"
	"pixi/api/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Exported handler function for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	log.Println("Successfully connected to the database.")

	controllers.PublishScheduledPosts()

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Scheduled posts published"))

	// Perform AutoMigrate for all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Save{},
		&models.Comment{},
		&models.Like{},
		&models.Follow{},
		&models.Reply{},
	)
	if err != nil {
		log.Fatal("Failed to run AutoMigrate:", err)
	}

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

	// Handle incoming HTTP requests using Gin
	rGin.ServeHTTP(w, r)
}
