package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type ApiConfig struct {
	port      string
	dbConnect string
}

func main() {
	godotenv.Load()
	cfg := ApiConfig{
		port:      os.Getenv("PORT"),
		dbConnect: os.Getenv("DB_CONNECT"),
	}
	fmt.Println(cfg.port)
	fmt.Println(cfg.dbConnect)

	server := initialiseServer(cfg, http.NewServeMux())

	log.Printf("Serving on port: %s\n", cfg.port)
	log.Panic(server.ListenAndServe())
}

func initialiseServer(cfg ApiConfig, mux *http.ServeMux) *http.Server {
	mux.HandleFunc("GET /v1/readiness", cfg.readinessHandler)
	mux.HandleFunc("GET /v1/err", cfg.errorHandler)

	corsMux := cfg.CorsMiddleware(mux)
	server := &http.Server{
		Addr:    ":" + cfg.port,
		Handler: corsMux,
	}
	return server
}

func (cfg *ApiConfig) CorsMiddleware(next http.Handler) http.Handler {
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
