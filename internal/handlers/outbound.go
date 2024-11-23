package handlers

import (
	"fmt"
	"net/http"
	"gorm.io/gorm"
)

// OutboundHandler retrieves all transactions from the database and returns them.
func OutboundHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define a struct to represent the transaction
		type Transaction struct {
			ID       string `json:"id"`
			Date     string `json:"date"`
			ShipTo   string `json:"ship_to"`
			ItemList string `json:"items"`
			Status   string `json:"status"`
		}

		var transactions []Transaction

		// Fetch all transactions from the database
		if err := db.Find(&transactions).Error; err != nil {
			http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
			return
		}

		// Return transactions in plain text EDI 856 format
		for _, t := range transactions {
			edi := fmt.Sprintf("EDI 856: Shipment %s to %s on %s with items: %s\n",
				t.ID, t.ShipTo, t.Date, t.ItemList)
			fmt.Fprintln(w, edi)
		}
	}
}