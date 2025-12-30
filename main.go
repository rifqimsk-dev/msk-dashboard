package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rifqimsk-dev/msk-dashboard/db"
	"github.com/rifqimsk-dev/msk-dashboard/handlers"
)

func main() {
	database := db.OpenDB()
	defer database.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/gender", handlers.GenderHandler(database))
	mux.HandleFunc("/api/total-entry", handlers.TotalEntryHandler(database))
	mux.HandleFunc("/api/total-amount", handlers.TotalAmountHandler(database))
	mux.HandleFunc("/api/toj", handlers.TojHandler(database))
	mux.HandleFunc("/api/type-motor", handlers.TypeMotorHandler(database))
	mux.HandleFunc("/api/pekerjaan", handlers.PekerjaanHandler(database))

    fmt.Println("Server run on port 8080")
	log.Fatal(http.ListenAndServe(":8080", CorsMiddleware(mux)))
}

func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        allowedOrigins := []string{
            "http://localhost:3000",      
			"http://127.0.0.1:5500",
			"http://127.0.0.1:5501",
			"http://127.0.0.1:8000",
            "https://msk-dashboard-front.vercel.app",
        }

        origin := r.Header.Get("Origin")
        for _, o := range allowedOrigins {
            if o == origin {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
                break
            }
        }

        // Untuk preflight request
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

