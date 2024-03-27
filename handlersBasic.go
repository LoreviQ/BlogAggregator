package main

import (
	"net/http"
)

func (cfg *apiConfig) getReadiness(w http.ResponseWriter, r *http.Request) {
	// RESPONSE
	respondWithJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}

func (cfg *apiConfig) getError(w http.ResponseWriter, r *http.Request) {
	// RESPONSE
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
