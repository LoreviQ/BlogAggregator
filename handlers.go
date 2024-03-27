package main

import (
	"net/http"
	"time"

	"github.com/LoreviQ/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getReadiness(w http.ResponseWriter, r *http.Request) {
	// RESPONSE
	respondWithJSON(w, 200, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}

func (cfg *apiConfig) getError(w http.ResponseWriter, r *http.Request) {
	// RESPONSE
	respondWithError(w, 500, "Internal Server Error")
}

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	// REQUEST
	request, err := decodeRequest(w, r, struct {
		Name string `json:"name"`
	}{})
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
	respondWithJSON(w, http.StatusOK, struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
	}{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
	})
}

func (cfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request, user database.User) {
	// RESPONSE
	respondWithJSON(w, http.StatusOK, struct {
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
	})
}

func (cfg *apiConfig) postFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	// REQUEST
	request, err := decodeRequest(w, r, struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}{})
	if err != nil {
		respondWithError(w, 500, "failed to decode")
		return
	}

	// CREATE FEED
	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      request.Name,
		Url:       request.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// RESPONSE
	respondWithJSON(w, http.StatusOK, struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserID    uuid.UUID `json:"user_id"`
	}{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	})
}
