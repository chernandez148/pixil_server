package config

import (
	"encoding/json"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global DB variable
var DB *gorm.DB

// ConnectDB initializes and returns a DB connection
func ConnectDB() (*gorm.DB, error) {
	dsn := "host=" + os.Getenv("PGHOST") +
		" user=" + os.Getenv("PGUSER") +
		" password=" + os.Getenv("PGPASSWORD") +
		" dbname=" + os.Getenv("PGDATABASE") +
		" sslmode=require"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return DB, nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}

// Handler is the Vercel serverless function entry point
func Handler(w http.ResponseWriter, r *http.Request) {
	// Example: Attempt to connect to the database
	db, err := ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to the database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// You need to handle the two return values from db.DB() correctly
	sqlDB, err := db.DB() // db.DB() returns *sql.DB and an error
	if err != nil {
		http.Error(w, "Failed to get raw SQL DB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer sqlDB.Close() // Make sure to close the connection

	// Do something with db (e.g., fetch some data, etc.)
	// For demonstration, let's return a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Connected to the database", "status": "success"}
	json.NewEncoder(w).Encode(response)
}
