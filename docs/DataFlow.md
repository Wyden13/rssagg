# RSS Aggregator - Data Flow

| Step | Components           | Action                                   |
| ---- | -------------------- | ---------------------------------------- |
| 1    | Web Server           | User adds https://example.com/rss.       |
| 2    | Database             | Store Feed URL + User ID.                |
| 3    | Scraper (Go Routine) | Periodically fetch XML from Example.com. |
| 4    | Parser               | Convert XML → []Feed structs.            |
| 5    | Database             | Save only new feeds to the feeds table.  |
| 6    | Client               | User requests latest news via JSON.      |
