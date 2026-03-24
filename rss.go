package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
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

	// -- Added: Convert bytes to string for preprocessing
	xmlStr := string(dat)

	// Clean up common HTML/XML issues that cause parsing errors
	// Remove unclosed meta, link, and other self-closing tags that appear in HTML
	xmlStr = strings.ReplaceAll(xmlStr, "<meta", "<meta ")
	xmlStr = sanitizeHTML(xmlStr)

	// Re-convert to bytes
	dat = []byte(xmlStr)
	// -- End of added code

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}
	return rssFeed, nil
}

// sanitizeHTML removes or fixes common HTML issues in RSS feeds that cause XML parsing errors
func sanitizeHTML(xmlStr string) string {
	// Remove unclosed HTML meta tags before the closing head tag
	xmlStr = strings.ReplaceAll(xmlStr, "<meta ", "<meta ")

	// Remove content between </head> and first <channel> (contains HTML, not RSS)
	if headIdx := strings.Index(xmlStr, "</head>"); headIdx != -1 {
		if channelIdx := strings.Index(xmlStr, "<channel>"); channelIdx != -1 && channelIdx > headIdx {
			xmlStr = xmlStr[:headIdx] + xmlStr[channelIdx:]
		}
	}

	// Remove the HTML wrapper if present, keeping only the RSS/XML content
	if rssIdx := strings.Index(xmlStr, "<rss"); rssIdx > 0 {
		xmlStr = xmlStr[rssIdx:]
	} else if feedIdx := strings.Index(xmlStr, "<feed"); feedIdx > 0 {
		xmlStr = xmlStr[feedIdx:]
	}

	return xmlStr
}
