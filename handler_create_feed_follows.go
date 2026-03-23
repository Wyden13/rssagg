package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Wyden13/rssagg/db"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) createFeedFollowHandler(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		// UserID string `json:"user_id"`
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint("Failed to decoded JSON data: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), db.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint("Failed to create feed follow: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedFollowToAPIFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) getFeedFollowsHandler(w http.ResponseWriter, r *http.Request, user db.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint("Failed to get feed: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}
