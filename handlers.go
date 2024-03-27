package main

import (
	"net/http"
)

func (cfg *apiConfig) getReadiness(w http.ResponseWriter, r *http.Request) {
	responseStruct := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	respondWithJSON(w, 200, responseStruct)
}

func (cfg *apiConfig) getError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	// REQUEST
	requestStruct := struct {
		Name string `json:"name"`
	}{}
	_, err := decodeRequest(w, r, requestStruct)
	if err != nil {
		respondWithError(w, 500, "failed to decode")
		return
	}

	// RESPONSE
	responseStruct := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	respondWithJSON(w, 200, responseStruct)
}
