package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func respondWithJSON[T any](w http.ResponseWriter, responseCode int, body T) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(&body)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(responseCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, responseCode int, errorMsg string) {
	responseStruct := struct {
		Error string `json:"error"`
	}{
		Error: errorMsg,
	}
	data, err := json.Marshal(&responseStruct)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(responseCode)
	w.Write(data)
}

func decodeRequest[T any](w http.ResponseWriter, r *http.Request, _ T) (T, error) {
	var request T
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return request, err
	}
	return request, nil
}

func convertNulltime(time sql.NullTime) *time.Time {
	if time.Valid {
		return &time.Time
	}
	return nil
}
