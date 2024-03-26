package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON[T any](w http.ResponseWriter, responseCode int, body T) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(responseCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, responseCode int, errorMsg string) {
	type ReturnVals struct {
		Error string `json:"error"`
	}
	data, err := json.Marshal(ReturnVals{Error: errorMsg})
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(responseCode)
	w.Write(data)
}
