package config

import (
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
