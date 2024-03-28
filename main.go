package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LoreviQ/BlogAggregator/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	port     string
	DB       *database.Queries
	interval time.Duration
	noFeeds  int32
}

func main() {
	godotenv.Load()
	// Connect to DB
	db, err := sql.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Panic(err)
	}
	// Setup Config
	cfg := apiConfig{
		port:     os.Getenv("PORT"),
		DB:       database.New(db),
		interval: 20 * time.Second,
		noFeeds:  2,
	}
	// Start Scraper
	go cfg.startScraper()
	// Initialise Server
	server := initialiseServer(cfg, http.NewServeMux())
	// Serve Server
	log.Printf("Serving on port: %s\n", cfg.port)
	log.Panic(server.ListenAndServe())
}

func initialiseServer(cfg apiConfig, mux *http.ServeMux) *http.Server {
	mux.HandleFunc("GET /v1/readiness", cfg.getReadiness)
	mux.HandleFunc("GET /v1/err", cfg.getError)
	mux.HandleFunc("POST /v1/users", cfg.postUser)
	mux.HandleFunc("GET /v1/users", cfg.AuthMiddleware(cfg.getUser))
	mux.HandleFunc("POST /v1/feeds", cfg.AuthMiddleware(cfg.postFeed))
	mux.HandleFunc("GET /v1/feeds", cfg.getFeeds)
	mux.HandleFunc("GET /v1/feed_follows", cfg.AuthMiddleware(cfg.getFeedFollows))
	mux.HandleFunc("POST /v1/feed_follows", cfg.AuthMiddleware(cfg.postFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.deleteFeedFollow)

	corsMux := cfg.CorsMiddleware(mux)
	server := &http.Server{
		Addr:    ":" + cfg.port,
		Handler: corsMux,
	}
	return server
}
