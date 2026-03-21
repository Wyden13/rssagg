package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Wyden13/rssagg/db"
	"github.com/google/uuid"
)

// Handler for creating a new user
// apiCfg *apiConfig: a pointer to the apiCfg struct,
// which contains the database queries and other configuration needed to handle the request
func (apiCfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"name"`
		// CreatedAt time.Time `json:"created_at"`
		// UpdatedAt time.Time `json:"updated_at"`
	}
	// Parse the JSON request into the parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to decode request body: %v", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		Username:  params.Username,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, databaseUserToAPIUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user db.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToAPIUser(user))
}
