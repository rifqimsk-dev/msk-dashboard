package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rifqimsk-dev/msk-dashboard/models"
)

func PekerjaanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Month int `json:"month"`
			Year  int `json:"year"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		if req.Month < 1 || req.Month > 12 {
			http.Error(w, "Month harus 1-12", http.StatusBadRequest)
			return
		}

		if req.Year < 2000 || req.Year > 2100 {
			http.Error(w, "Year tidak valid", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`
			SELECT pekerjaan, COUNT(*) AS total 
			FROM rmsk17 
			WHERE pekerjaan IS NOT NULL 
			AND MONTH(SaleOrderdate) = ?
			AND YEAR(SaleOrderdate) = ?
			AND TRIM(pekerjaan) <> '' 
			GROUP BY pekerjaan 
			ORDER BY total DESC 
			LIMIT 5
		`, req.Month, req.Year)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var result []models.JenisPekerjaan
		for rows.Next() {
			var row models.JenisPekerjaan
			if err := rows.Scan(&row.Pekerjaan, &row.Total); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			result = append(result, row)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
