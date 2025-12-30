package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rifqimsk-dev/msk-dashboard/models"
)

func GenderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ========================
		// METHOD CHECK
		// ========================
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// ========================
		// REQUEST BODY
		// ========================
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

		// ========================
		// DATE RANGE
		// ========================
		loc, _ := time.LoadLocation("Asia/Jakarta")

		awalBulan := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, loc)
		akhirBulan := awalBulan.AddDate(0, 1, 0)
		awalBulanLalu := awalBulan.AddDate(0, -1, 0)

		// ========================
		// RESULT VARIABLES
		// ========================
		var (
			lk, pr             float64
			lkPrev, prPrev     float64
			persenLk           float64
			persenPr           float64
			persenLkPrev       float64
			persenPrPrev       float64
		)

		// ========================
		// QUERY
		// ========================
		err := db.QueryRow(`
			SELECT
				-- TOTAL BULAN INI
				COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin = 1 THEN 1 END), 0) AS lk,
				COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin = 2 THEN 1 END), 0) AS pr,

				-- TOTAL BULAN LALU
				COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin = 1 THEN 1 END), 0) AS lk_prev,
				COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin = 2 THEN 1 END), 0) AS pr_prev,

				-- PERSENTASE BULAN INI
				ROUND(
					COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin = 1 THEN 1 END),0)
					* 100 /
					NULLIF(
						COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin IN (1,2) THEN 1 END),0),
						0
					), 2
				) AS persen_lk,

				ROUND(
					COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin = 2 THEN 1 END),0)
					* 100 /
					NULLIF(
						COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin IN (1,2) THEN 1 END),0),
						0
					), 2
				) AS persen_pr,

				-- PERSENTASE BULAN LALU
				ROUND(
					COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin = 1 THEN 1 END),0)
					* 100 /
					NULLIF(
						COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin IN (1,2) THEN 1 END),0),
						0
					), 2
				) AS persen_lk_prev,

				ROUND(
					COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin = 2 THEN 1 END),0)
					* 100 /
					NULLIF(
						COALESCE(SUM(CASE WHEN SaleOrderdate >= ? AND SaleOrderdate < ? AND JenisKelamin IN (1,2) THEN 1 END),0),
						0
					), 2
				) AS persen_pr_prev

			FROM rmsk17
		`,
			// bulan ini
			awalBulan, akhirBulan,
			awalBulan, akhirBulan,

			// bulan lalu
			awalBulanLalu, awalBulan,
			awalBulanLalu, awalBulan,

			// persen bulan ini
			awalBulan, akhirBulan,
			awalBulan, akhirBulan,
			awalBulan, akhirBulan,
			awalBulan, akhirBulan,

			// persen bulan lalu
			awalBulanLalu, awalBulan,
			awalBulanLalu, awalBulan,
			awalBulanLalu, awalBulan,
			awalBulanLalu, awalBulan,
		).Scan(
			&lk, &pr,
			&lkPrev, &prPrev,
			&persenLk, &persenPr,
			&persenLkPrev, &persenPrPrev,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// ========================
		// RESPONSE
		// ========================
		resp := models.Gender{
			LakiLaki:           lk,
			Perempuan:          pr,
			TotalLkSebelumnya:  lkPrev,
			TotalPrSebelumnya:  prPrev,
			PersenLk:           persenLk,
			PersenPr:           persenPr,
			PersenLkSebelumnya: persenLkPrev,
			PersenPrSebelumnya: persenPrPrev,
		}

		json.NewEncoder(w).Encode(resp)

	}
}
