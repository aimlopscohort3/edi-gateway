// Placeholder for database_activity.go
package workflows

import (
	"context"
	"log"

	"gorm.io/gorm"
)

// SaveToDatabaseActivity saves a transaction to the database.
func SaveToDatabaseActivity(ctx context.Context, db *gorm.DB, transaction Transaction) error {
	log.Printf("Saving transaction %s to the database", transaction.ID)
	if err := db.Create(&transaction).Error; err != nil {
		log.Printf("Error saving transaction %s: %v", transaction.ID, err)
		return err
	}
	log.Printf("Transaction %s saved successfully", transaction.ID)
	return nil
}