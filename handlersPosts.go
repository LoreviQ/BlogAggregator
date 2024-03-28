package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/LoreviQ/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	// GET POSTS FROM DB
	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  cfg.noPosts,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// RESPONSE
	type response struct {
		ID          uuid.UUID      `json:"id"`
		CreatedAt   time.Time      `json:"created_at"`
		UpdatedAt   time.Time      `json:"updated_at"`
		Title       sql.NullString `json:"title"`
		Url         string         `json:"url"`
		Description sql.NullString `json:"description"`
		PublishedAt sql.NullTime   `json:"published_at"`
		FeedID      uuid.UUID      `json:"feed_id"`
	}
	responseSlice := make([]response, 0, len(posts))
	for _, post := range posts {
		responseSlice = append(responseSlice, response{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			Url:         post.Url,
			Description: post.Description,
			PublishedAt: post.PublishedAt,
			FeedID:      post.FeedID,
		})
	}
	respondWithJSON(w, http.StatusOK, responseSlice)
}
