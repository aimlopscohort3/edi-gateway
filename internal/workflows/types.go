// Placeholder for types.go
package workflows

import "time"

// Transaction defines the structure for an EDI transaction.
type Transaction struct {
	ID       string    `json:"id" gorm:"primaryKey"`
	Date     time.Time `json:"date"`
	ShipTo   string    `json:"ship_to"`
	ItemList string    `json:"items"`
	Status   string    `json:"status"`
}