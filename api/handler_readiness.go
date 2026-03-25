package main

import "net/http"

// A handler function to handle HTTP requests in a standard way,
// it takes in two parameters: an http.ResponseWriter and a pointer *http.Request
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, struct{}{})
}
