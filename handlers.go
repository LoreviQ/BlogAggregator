package main

import (
	"net/http"
)

func (cfg *apiConfig) readinessHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}

func (cfg *apiConfig) errorHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}
