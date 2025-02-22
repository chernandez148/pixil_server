package main

import (
	"log"
	"net/http"
	"pixi/api/handler"
	"pixi/config"
	"pixi/models"

	"github.com/gin-gonic/gin"
)

// Exported Handler function for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Connect to the database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	log.Println("Successfully connected to the database.")

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

	// Initialize the Gin router and handle the request
	handler.Handler(db)(w, r)
}
