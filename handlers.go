package main

import (
	"net/http"
	"time"

	"github.com/LoreviQ/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getReadiness(w http.ResponseWriter, r *http.Request) {
	// RESPONSE
	responseStruct := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	respondWithJSON(w, 200, responseStruct)
}

func (cfg *apiConfig) getError(w http.ResponseWriter, r *http.Request) {
	// RESPONSE
	respondWithError(w, 500, "Internal Server Error")
}

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	// REQUEST
	requestStruct := struct {
		Name string `json:"name"`
	}{}
	request, err := decodeRequest(w, r, requestStruct)
	if err != nil {
		respondWithError(w, 500, "failed to decode")
		return
	}

	// CREATE USER
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      request.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	// RESPONSE
	responseStruct := struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
	}{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
	}
	respondWithJSON(w, http.StatusOK, responseStruct)
}

func (cfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request, user database.User) {
	// RESPONSE
	responseStruct := struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		ApiKey    string    `json:"api_key"`
	}{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
	respondWithJSON(w, http.StatusOK, responseStruct)
}
