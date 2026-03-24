package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Wyden13/rssagg/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *db.Queries
}

// initialize go module :go mod init github.com/yourusername/job-aggregator
// go mod vendor: create a local copy of external dependencies required to build the project
func main() {
	// feed, err := urlToFeed("https://www.wagslane.dev/index.xml")
	// if err != nil {
	// 	log.Fatal("Failed to fetch and parse RSS feed: %v", err)
	// }
	// fmt.Printf("Fetched feed: %+v\n", feed)

	// Load environment variables from .env file
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not set in environment variables")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set in environment variables")
	}

	// Connect to the database using the provided URL
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database: %v", err)
	}

	// queries, err := db.New(conn)
	// if err != nil {
	// 	log.Fatal("Failed to create database queries: %v", err)
	// }

	db := db.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

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
	v1Router.Post("/users", apiCfg.createUserHandler)
	v1Router.Get("/users", apiCfg.authMiddleware(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.authMiddleware(apiCfg.createFeedHandler))
	v1Router.Get("/feeds", apiCfg.getFeedsHandler)

	v1Router.Post("/feed_follows", apiCfg.authMiddleware(apiCfg.createFeedFollowHandler))
	v1Router.Get("/feed_follows", apiCfg.authMiddleware(apiCfg.getFeedFollowsHandler))
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
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server will run on port: %s\n", portString)
}
