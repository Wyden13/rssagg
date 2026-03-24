package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/Wyden13/rssagg/db"
	"github.com/google/uuid"
)

// startScraper starts the scraper that periodically fetches RSS feeds and stores new posts in the database.
func startScraping(
	database *db.Queries, // The database connection to use for storing scraped data
	concurrency int, // The number of concurrent scraper goroutines to run
	timeBetweenRequest time.Duration, // The time to wait between each request
) {
	log.Printf("Scraping on %v goroutines with %v between each request", concurrency, timeBetweenRequest)

	// Create a ticker that will trigger at the specified interval for making requests
	ticker := time.NewTicker(timeBetweenRequest)

	// Create a background worker that runs a task at a set interval
	// time.Ticker has a channel called "C", the ticker sends the current time onto that channel every time the interval elapses
	//
	for ; ; <-ticker.C {
		feeds, err := database.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("Error fetching feeds to scrape: %v", err)
			continue
		}
		wg := &sync.WaitGroup{}
		// For every feed that we need to scrape, we will start a new goroutine to scrape it concurrently
		// Every time we start a new goroutine, we will add 1 to the wait group counter using wg.Add(1)
		// For example:
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(database, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(queries *db.Queries, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()

	_, err := queries.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed %v as fetched: %v", feed.Name, err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed %v: %v", feed.Url, err)
		return
	}
	for _, item := range rssFeed.Channel.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{String: item.Description, Valid: true}
		}
		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing pubDate %v: %v", item.PubDate, err)
			return
		}
		_, err = queries.CreatePost(context.Background(),
			db.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Description: description,
				PublishedAt: sql.NullTime{Time: t, Valid: true},
				Url:         item.Link,
				FeedID:      feed.ID,
			})
		if err != nil {
			log.Printf("Error creating post: %v", err)
			continue
		}

		// log.Printf("Found post: %v on feed: %v", item.Title, feed.Name)
	}
	log.Printf("Fetched feed %v with %v items", feed.Name, len(rssFeed.Channel.Items))
}
