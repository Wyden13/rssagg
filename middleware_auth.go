package main

import (
	"fmt"
	"net/http"

	"github.com/Wyden13/rssagg/auth"
	"github.com/Wyden13/rssagg/db"
)

// authedHandler is a function type definition, used to handle HTTP requests with authorized User.
// It takes three arguments:
// - http.ResponseWriter: To send data back to the client.
// - *http.Request: To read data from the incoming request.
// - db.User: It's the authenticated user object already fetched from your database.
type authedHandler func(http.ResponseWriter, *http.Request, db.User)

// authMiddleware is a middleware function that checks for a valid API key in the request header
// and retrieves the corresponding user from the database.
func (apiCfg *apiConfig) authMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Failed to get API key from request header: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Failed to get user by API key: %v", err))
			return
		}
		handler(w, r, user)
	}
}
