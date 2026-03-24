package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Wyden13/rssagg/db"
)

// startScraper starts the scraper that periodically fetches RSS feeds and stores new posts in the database.
func startScraping(
	db *db.Queries, // The database connection to use for storing scraped data
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
		feeds, err := db.GetNextFeedToFetch(
			context.Background(),
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
			go scrapeFeed(wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *db.Queries, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
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
		log.Print("Found post:", item.Title)
	}
	log.Printf("Fetched feed %v with %v items", feed.Name, len(rssFeed.Channel.Items))
}
