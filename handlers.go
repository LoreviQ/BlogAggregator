package main

import (
	"net/http"
	"time"

	"github.com/LoreviQ/BlogAggregator/internal/database"
	"github.com/google/uuid"
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

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	// REQUEST
	request, err := decodeRequest(w, r, struct {
		Name string `json:"name"`
	}{})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to decode request body")
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
		respondWithError(w, http.StatusBadRequest, "failed to decode request body")
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

func (cfg *apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve feeds from database")
		return
	}

	// RESPONSE
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Url       string    `json:"url"`
		UserID    uuid.UUID `json:"user_id"`
	}
	responseSlice := make([]response, 0, len(feeds))
	for _, feed := range feeds {
		responseSlice = append(responseSlice, response{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserID:    feed.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, responseSlice)
}

func (cfg *apiConfig) postFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	// REQUEST
	request, err := decodeRequest(w, r, struct {
		FeedID string `json:"feed_id"`
	}{})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}
	feedID, err := uuid.Parse(request.FeedID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Feed ID")
		return
	}

	// CREATE FEED FOLLOW
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// RESPONSE
	respondWithJSON(w, http.StatusOK, struct {
		ID        uuid.UUID `json:"id"`
		FeedID    uuid.UUID `json:"feed_id"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        feedFollow.ID,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
	})
}

func (cfg *apiConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	// GET FFID FROM PATH
	FFid, err := uuid.Parse(r.PathValue("feedFollowID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// DELETE FF
	err = cfg.DB.DeleteFeedFollow(r.Context(), FFid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
