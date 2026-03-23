package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/Wyden13/rssagg/db"
)

func (apiCfg *apiConfig) createFeedHandler(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	// Decode the JSON body into the parameters struct
	decoder := json.NewDecoder(r.Body)

	// Initialize the parameters struct and decode the JSON into it
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint("invalid JSON body: ", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), db.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	respondWithJSON(w, http.StatusOK, databaseFeedToAPIFeed(feed))
}

// GetFeedsHandler is a function
func (apiCfg *apiConfig) getFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint("failed to get feeds: %v ", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedsToAPIFeeds(feeds))
}
