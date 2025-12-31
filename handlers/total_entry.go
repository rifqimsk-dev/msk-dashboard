package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func TotalEntryHandler(db *sql.DB) http.HandlerFunc {
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

		now := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.Local)

		// bulan sekarang
		startCurrent := now
		endCurrent := now.AddDate(0, 1, 0)

		// bulan sebelumnya
		startPrev := now.AddDate(0, -1, 0)
		endPrev := now

		var currentMonth, previousMonth int

		err := db.QueryRow(`
			SELECT
				SUM(CASE
					WHEN SaleOrderdate >= ? AND SaleOrderdate < ?
					THEN 1 ELSE 0 END
				) AS bulan_sekarang,

				SUM(CASE
					WHEN SaleOrderdate >= ? AND SaleOrderdate < ?
					THEN 1 ELSE 0 END
				) AS bulan_sebelumnya
			FROM rmsk17
		`,
			startCurrent,
			endCurrent,
			startPrev,
			endPrev,
		).Scan(&currentMonth, &previousMonth)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{
			"total_entry":            currentMonth,
			"total_entry_sebelumnya": previousMonth,
		})
	}
}
