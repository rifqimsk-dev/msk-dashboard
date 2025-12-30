package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rifqimsk-dev/msk-dashboard/models"
)

func TypeMotorHandler(db *sql.DB) http.HandlerFunc {
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

		month := req.Month
		year := req.Year

		prevMonth := month - 1
		prevYear := year
		if month == 1 {
			prevMonth = 12
			prevYear = year - 1
		}


		rows, err := db.Query(`
			SELECT 
				t.TypeMotorDescription,
				COUNT(*) AS total,
				COALESCE(prev.total_sebelumnya, 0) AS total_sebelumnya
			FROM (
				SELECT TypeMotorDescription, SaleOrderdate
				FROM rmsk17
				WHERE TypeMotorDescription IS NOT NULL
				AND YEAR(SaleOrderdate) = ?
				AND MONTH(SaleOrderdate) = ?
			) t
			LEFT JOIN (
				SELECT TypeMotorDescription, COUNT(*) AS total_sebelumnya
				FROM rmsk17
				WHERE TypeMotorDescription IS NOT NULL
				AND YEAR(SaleOrderdate) = ?
				AND MONTH(SaleOrderdate) = ?
				GROUP BY TypeMotorDescription
			) prev ON prev.TypeMotorDescription = t.TypeMotorDescription
			GROUP BY t.TypeMotorDescription, prev.total_sebelumnya
			ORDER BY total DESC
			LIMIT 10;
		`, year, month, prevYear, prevMonth)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var result []models.TypeMotorStat
		for rows.Next() {
			var row models.TypeMotorStat
			if err := rows.Scan(&row.Type, &row.Total, &row.TotalSebelumnya); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			result = append(result, row)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
