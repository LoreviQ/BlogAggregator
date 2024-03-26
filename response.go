package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type JSONResponse interface {
	convertToJson() ([]byte, error)
}

func respondWithJSON(w http.ResponseWriter, responseCode int, body JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	data, err := body.convertToJson()
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
