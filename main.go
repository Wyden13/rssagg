package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

// initialize go module :go mod init github.com/yourusername/job-aggregator
// go mod vendor: create a local copy of external dependencies required to build the project
func main() {
	fmt.Println("Hello, World!")

	// Load environment variables from .env file
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		// portString = "8000"
		log.Fatal("PORT is not set in environment variables")
	}
	// Create a new router using chi
	router := chi.NewRouter()

	// Configure CORS settings for the router
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"}, // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedHeaders:   []string{"*"}, // Allow all headers
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/ready", readinessHandler)
	v1Router.Get("/error", errorHandler)

	// v1Router.HandleFunc("/ready", handlerReadiness) --- IGNORE ---
	// v1Router.HandleFunc("/error", errorHandler) --- IGNORE ---

	router.Mount("/v1", v1Router)

	// Configure CORS settings
	srv := &http.Server{
		Addr:    ":" + portString,
		Handler: router,
	}

	// Listen and serve(handle) incoming HTTP requests
	// Stop running right here, if there is an error while starting the server,
	// srv.ListenAndServe() will return an error, and we log it and exit the program using log.Fatal()
	log.Printf("Starting server on port: %s\n", portString)
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server will run on port: %s\n", portString)
}
