package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rifqimsk-dev/msk-dashboard/models"
)

func PekerjaanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT pekerjaan, COUNT(*) AS total 
			FROM sales_service_data 
			WHERE pekerjaan IS NOT NULL 
			AND TRIM(pekerjaan) <> '' 
			GROUP BY pekerjaan 
			ORDER BY total DESC 
			LIMIT 5
		`)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var result []models.JenisPekerjaan
		for rows.Next() {
			var row models.JenisPekerjaan
			if err := rows.Scan(&row.Pekerjaan, &row.Total); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			result = append(result, row)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
