package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func GenderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var lk, pr, persenlk, persenpr float64

		err := db.QueryRow(`
			SELECT
				totalLk,
				totalPr,
				ROUND(totalLk * 100.0 / (totalLk + totalPr), 2) AS persenLk,
				ROUND(totalPr * 100.0 / (totalLk + totalPr), 2) AS persenPr
			FROM (
				SELECT
					COUNT(CASE WHEN JenisKelamin = 1 THEN 1 END) AS totalLk,
					COUNT(CASE WHEN JenisKelamin = 2 THEN 1 END) AS totalPr
				FROM sales_service_data
			) AS sub;
		`).Scan(&lk, &pr, &persenlk, &persenpr)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]float64{
			"laki_laki": lk,
			"perempuan": pr,
			"persenlk":  persenlk,
			"persenpr":  persenpr,
		})
	}
}
