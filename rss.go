package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// TODO: implement this function to fetch and parse the RSS feed in XML format from the given URL.
// The function should return an RSSFeed struct containing the parsed data or an error if something goes wrong.
func urlToFeed(url string) (RSSFeed, error) {
	// Create an HTTP client with a timeout to avoid hanging indefinitely.
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	// Fetch the RSS feed from the URL.
	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}

	// Ensure the response body is closed after we're done with it to prevent resource leaks.
	defer resp.Body.Close()

	// Get all the data from response body
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}
	return rssFeed, nil
}
