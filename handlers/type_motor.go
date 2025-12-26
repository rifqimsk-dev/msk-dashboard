package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rifqimsk-dev/msk-dashboard/models"
)

func TypeMotorHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT TypeMotorDescription, COUNT(*) AS total
			FROM sales_service_data
			WHERE TypeMotorDescription IS NOT NULL
			GROUP BY TypeMotorDescription
			ORDER BY total DESC
			LIMIT 10
		`)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var result []models.TypeMotorStat
		for rows.Next() {
			var row models.TypeMotorStat
			if err := rows.Scan(&row.Type, &row.Total); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			result = append(result, row)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
