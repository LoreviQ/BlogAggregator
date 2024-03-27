package main

import (
	"net/http"
)

func (cfg *apiConfig) getReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}

func (cfg *apiConfig) getError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	// REQUEST
	_, err := decodeRequest(w, r, struct {
		Name string `json:"name"`
	}{})
	if err != nil {
		respondWithError(w, 500, "failed to decode")
		return
	}

	respondWithJSON(w, 200, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}
