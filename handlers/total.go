package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func TotalEntryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var total int
		err := db.QueryRow(`SELECT COUNT(*) FROM sales_service_data`).Scan(&total)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]int{"total_entry": total})
	}
}

func TotalAmountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var total float64
		err := db.QueryRow(`SELECT SUM(TotalAmount) FROM sales_service_data`).Scan(&total)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]float64{"total_amount": total})
	}
}
