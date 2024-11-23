package services

import (
	"log"

	"edi-gateway/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDatabase initializes the database connection using the provided configuration.
func InitDatabase(cfg config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database connection initialized successfully.")
	return db
}