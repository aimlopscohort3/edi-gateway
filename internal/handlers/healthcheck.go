package handlers

import (
	"fmt"
	"net/http"
	"gorm.io/gorm"
)

// HealthCheckHandler checks the health of the database and other services.
func HealthCheckHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			http.Error(w, "Database connection unhealthy", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	}
}