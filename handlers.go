package main

import "net/http"

func (cfg *ApiConfig) readinessHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct {
		Status string `json:"status"`
	}{Status: "ok"})
}
