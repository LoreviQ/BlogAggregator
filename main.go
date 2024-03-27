package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/LoreviQ/BlogAggregator/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	port string
	DB   *database.Queries
}

func main() {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Panic(err)
	}

	cfg := apiConfig{
		port: os.Getenv("PORT"),
		DB:   database.New(db),
	}

	server := initialiseServer(cfg, http.NewServeMux())

	log.Printf("Serving on port: %s\n", cfg.port)
	log.Panic(server.ListenAndServe())
}

func initialiseServer(cfg apiConfig, mux *http.ServeMux) *http.Server {
	mux.HandleFunc("GET /v1/readiness", cfg.readinessHandler)
	mux.HandleFunc("GET /v1/err", cfg.errorHandler)

	corsMux := cfg.CorsMiddleware(mux)
	server := &http.Server{
		Addr:    ":" + cfg.port,
		Handler: corsMux,
	}
	return server
}

func (cfg *apiConfig) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
