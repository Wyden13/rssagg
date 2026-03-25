# RSS Aggregator - Function Documentation

This document provides a comprehensive reference for all functions in the RSS Aggregator project.

## Table of Contents

1. [HTTP Handlers](#http-handlers)
2. [Middleware](#middleware)
3. [Database Functions](#database-functions)
4. [Utility Functions](#utility-functions)
5. [Data Models](#data-models)

---

## HTTP Handlers

HTTP handlers are functions that process incoming requests and send responses back to clients. They follow the standard Go HTTP handler pattern.

### User Handlers

#### `createUserHandler(w http.ResponseWriter, r *http.Request)`

**Location:** `handler_user.go`

Creates a new user account in the system.

**Request Body:**

```json
{
  "name": "username"
}
```

**Response (200 OK):**

```json
{
  "id": "uuid",
  "name": "username",
  "created_at": "2026-03-25T10:00:00Z",
  "api_key": "auto-generated-key"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid JSON body
- `500 Internal Server Error` - Database error

**Access:** Public (no authentication required)

---

#### `handlerGetUser(w http.ResponseWriter, r *http.Request, user db.User)`

**Location:** `handler_user.go`

Retrieves the authenticated user's information.

**Authentication:** Required (API key in Authorization header)

**Response (200 OK):**

```json
{
  "id": "uuid",
  "name": "username",
  "created_at": "2026-03-25T10:00:00Z",
  "api_key": "user-api-key"
}
```

**Error Responses:**

- `403 Forbidden` - Invalid or missing API key

---

#### `getPostsForUserHandler(w http.ResponseWriter, r *http.Request, user db.User)`

**Location:** `handler_user.go`

Retrieves all posts from feeds followed by the authenticated user.

**Authentication:** Required (API key in Authorization header)

**Response (200 OK):**

```json
[
  {
    "id": "uuid",
    "title": "Post Title",
    "description": "Post description",
    "url": "https://example.com/post",
    "published_at": "2026-03-25T10:00:00Z",
    "feed_id": "uuid"
  }
]
```

**Error Responses:**

- `403 Forbidden` - Invalid or missing API key
- `500 Internal Server Error` - Database error

---

### Feed Handlers

#### `createFeedHandler(w http.ResponseWriter, r *http.Request, user db.User)`

**Location:** `handler_feed.go`

Creates a new RSS feed subscription for the authenticated user.

**Authentication:** Required

**Request Body:**

```json
{
  "name": "Feed Name",
  "url": "https://example.com/rss.xml"
}
```

**Response (200 OK):**

```json
{
  "id": "uuid",
  "name": "Feed Name",
  "url": "https://example.com/rss.xml",
  "user_id": "uuid",
  "created_at": "2026-03-25T10:00:00Z",
  "updated_at": "2026-03-25T10:00:00Z"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid JSON body or failed to create feed
- `403 Forbidden` - Invalid authentication

---

#### `getFeedsHandler(w http.ResponseWriter, r *http.Request)`

**Location:** `handler_feed.go`

Retrieves all RSS feeds in the system.

**Authentication:** Not required

**Response (200 OK):**

```json
[
  {
    "id": "uuid",
    "name": "Feed Name",
    "url": "https://example.com/rss.xml",
    "user_id": "uuid",
    "created_at": "2026-03-25T10:00:00Z",
    "updated_at": "2026-03-25T10:00:00Z"
  }
]
```

**Error Responses:**

- `400 Bad Request` - Database error

---

### Feed Follow Handlers

#### `createFeedFollowHandler(w http.ResponseWriter, r *http.Request, user db.User)`

**Location:** `handler_create_feed_follows.go`

Creates a feed follow relationship (user subscribes to a feed).

**Authentication:** Required

**Request Body:**

```json
{
  "feed_id": "uuid"
}
```

**Response (200 OK):**

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "feed_id": "uuid",
  "created_at": "2026-03-25T10:00:00Z",
  "updated_at": "2026-03-25T10:00:00Z"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid JSON body
- `403 Forbidden` - Invalid authentication
- `500 Internal Server Error` - Database error

---

#### `getFeedFollowsHandler(w http.ResponseWriter, r *http.Request, user db.User)`

**Location:** `handler_create_feed_follows.go`

Retrieves all feeds that the authenticated user follows.

**Authentication:** Required

**Response (200 OK):**

```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "feed_id": "uuid",
    "created_at": "2026-03-25T10:00:00Z",
    "updated_at": "2026-03-25T10:00:00Z"
  }
]
```

**Error Responses:**

- `403 Forbidden` - Invalid authentication
- `500 Internal Server Error` - Database error

---

### Health Check Handlers

#### `readinessHandler(w http.ResponseWriter, r *http.Request)`

**Location:** `handler_readiness.go`

Health check endpoint to verify the server is running and accessible.

**Authentication:** Not required

**Response (200 OK):**

```json
{
  "status": "ok"
}
```

---

#### `errorHandler(w http.ResponseWriter, r *http.Request)`

**Location:** `handler_err.go`

Intentional error endpoint for testing error handling (debug only).

**Authentication:** Not required

**Response (500 Internal Server Error):**

```json
{
  "error": "Internal server error"
}
```

---

## Middleware

Middleware functions process requests before they reach handlers.

### Authentication Middleware

#### `authMiddleware(handler authedHandler) http.HandlerFunc`

**Location:** `middleware_auth.go`

Validates API key from the request header and retrieves the associated user from the database.

**How it works:**

1. Extracts API key from `Authorization: Bearer <api_key>` header
2. Looks up user in database using the API key
3. Passes authenticated user to the handler
4. Returns 403 if API key is invalid or missing

**Type Definition:**

```go
type authedHandler func(http.ResponseWriter, *http.Request, db.User)
```

**Example Usage:**

```go
v1Router.Get("/users", apiCfg.authMiddleware(apiCfg.handlerGetUser))
```

---

## Database Functions

Database functions are generated by sqlc from SQL queries. They interact with the PostgreSQL database.

### User Operations

#### `CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)`

Creates a new user in the database.

**Parameters:**

- `id`: UUID (auto-generated)
- `username`: String (unique)
- `created_at`: Timestamp (auto-set to current time)
- `updated_at`: Timestamp (auto-set to current time)
- `api_key`: String (auto-generated SHA256 hash)

**Returns:** Created User object or error

---

#### `GetUserByAPIKey(ctx context.Context, apiKey string) (db.User, error)`

Retrieves a user by their API key.

**Parameters:**

- `apiKey`: User's unique API key

**Returns:** User object or error if not found

---

### Feed Operations

#### `CreateFeed(ctx context.Context, arg db.CreateFeedParams) (db.Feed, error)`

Creates a new RSS feed subscription.

**Parameters:**

- `id`: UUID (auto-generated)
- `created_at`: Timestamp
- `updated_at`: Timestamp
- `name`: Feed name
- `url`: Feed URL
- `user_id`: ID of user creating the feed

**Returns:** Created Feed object or error

---

#### `GetFeeds(ctx context.Context) ([]db.Feed, error)`

Retrieves all feeds in the system.

**Returns:** Slice of Feed objects or error

---

#### `GetNextFeedsToFetch(ctx context.Context, limit int32) ([]db.Feed, error)`

Retrieves feeds that haven't been fetched recently (used by scraper).

**Parameters:**

- `limit`: Maximum number of feeds to return

**Returns:** Slice of Feed objects or error

---

#### `MarkFeedAsFetched(ctx context.Context, feedID uuid.UUID) (db.Feed, error)`

Updates a feed's `last_fetched_at` timestamp.

**Parameters:**

- `feedID`: ID of feed to mark as fetched

**Returns:** Updated Feed object or error

---

### Feed Follow Operations

#### `CreateFeedFollow(ctx context.Context, arg db.CreateFeedFollowParams) (db.FeedFollow, error)`

Creates a user-feed subscription relationship.

**Parameters:**

- `id`: UUID (auto-generated)
- `user_id`: User ID
- `feed_id`: Feed ID
- `created_at`: Timestamp
- `updated_at`: Timestamp

**Returns:** Created FeedFollow object or error

---

#### `GetFeedFollowsByUserID(ctx context.Context, userID uuid.UUID) ([]db.FeedFollow, error)`

Retrieves all feeds followed by a specific user.

**Parameters:**

- `userID`: ID of user

**Returns:** Slice of FeedFollow objects or error

---

### Post Operations

#### `CreatePost(ctx context.Context, arg db.CreatePostParams) (db.Post, error)`

Creates a new post (RSS item) in the database.

**Parameters:**

- `id`: UUID
- `created_at`: Timestamp
- `updated_at`: Timestamp
- `title`: Post title
- `description`: Post content (nullable)
- `published_at`: Post publication date (nullable)
- `url`: Post URL
- `feed_id`: Associated feed ID

**Returns:** Created Post object or error

---

#### `GetPostsByFeedID(ctx context.Context, feedID uuid.UUID) ([]db.Post, error)`

Retrieves all posts from a specific feed.

**Parameters:**

- `feedID`: ID of feed

**Returns:** Slice of Post objects or error

---

## Utility Functions

### Response Functions

#### `respondWithError(w http.ResponseWriter, code int, message string)`

**Location:** `handler_err.go`

Sends a JSON error response with the given HTTP status code.

**Response Format:**

```json
{
  "error": "error message"
}
```

---

#### `respondWithJSON(w http.ResponseWriter, code int, payload interface{})`

**Location:** `json.go`

Sends a JSON response with the given HTTP status code.

**Response Format:**

```json
{
  // ... payload data
}
```

---

### Authentication Functions

#### `GetAPIKey(headers http.Header) (string, error)`

**Location:** `auth/auth.go`

Extracts and validates the API key from the `Authorization` header.

**Expected Header Format:** `Authorization: Bearer <api_key>`

**Returns:** API key string or error if missing/malformed

---

### Scraper Functions

#### `startScraping(queries *db.Queries, concurrency int, timeBetweenRequest time.Duration)`

**Location:** `scraper.go`

Starts the background RSS feed scraper that periodically fetches feeds and stores posts.

**Parameters:**

- `queries`: Database connection
- `concurrency`: Number of concurrent goroutines
- `timeBetweenRequest`: Time to wait between scrape cycles

**Behavior:**

- Runs continuously in background (infinite loop)
- Fetches feeds concurrently
- Marks feeds as fetched
- Stores new posts in database

---

#### `scrapeFeed(queries *db.Queries, wg *sync.WaitGroup, feed db.Feed)`

**Location:** `scraper.go`

Scrapes a single RSS feed, fetches its items, and stores them as posts.

**Parameters:**

- `queries`: Database connection
- `wg`: WaitGroup for synchronization
- `feed`: Feed object to scrape

**Process:**

1. Mark feed as fetched
2. Fetch RSS XML from feed URL
3. Parse items from RSS
4. Create post records for each item
5. Log errors or success

---

#### `urlToFeed(url string) (RSSFeed, error)`

**Location:** `rss.go`

Fetches and parses an RSS feed from a URL.

**Parameters:**

- `url`: URL of RSS feed

**Returns:** Parsed RSSFeed struct or error

**Handles:**

- HTTP requests with 10-second timeout
- XML parsing
- Malformed HTML/XML sanitization

---

### Data Conversion Functions

#### `databaseUserToAPIUser(dbUser db.User) User`

**Location:** `models.go`

Converts database User model to API response format.

**Conversion:**

- `id` (UUID) → `id` (string)
- Database fields → JSON-serializable fields

---

#### `databaseFeedToAPIFeed(dbFeed db.Feed) Feed`

**Location:** `models.go`

Converts database Feed model to API response format.

---

#### `databaseFeedsToAPIFeeds(dbFeeds []db.Feed) []Feed`

**Location:** `models.go`

Converts slice of database Feed models to API response format.

---

## Data Models

### API Models (Response Types)

#### `User`

**Location:** `models.go`

```go
type User struct {
  ID        string    `json:"id"`
  Name      string    `json:"name"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  ApiKey    string    `json:"api_key"`
}
```

---

#### `Feed`

**Location:** `models.go`

```go
type Feed struct {
  ID        string    `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  Name      string    `json:"name"`
  Url       string    `json:"url"`
  UserID    string    `json:"user_id"`
}
```

---

#### `FeedFollow`

**Location:** `models.go`

```go
type FeedFollow struct {
  ID        string    `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  UserID    string    `json:"user_id"`
  FeedID    string    `json:"feed_id"`
}
```

---

#### `Post`

**Location:** `models.go`

```go
type Post struct {
  ID          string    `json:"id"`
  CreatedAt   time.Time `json:"created_at"`
  UpdatedAt   time.Time `json:"updated_at"`
  Title       string    `json:"title"`
  Description string    `json:"description"`
  PublishedAt time.Time `json:"published_at"`
  Url         string    `json:"url"`
  FeedID      string    `json:"feed_id"`
}
```

---

### RSS Models

#### `RSSFeed`

**Location:** `rss.go`

```go
type RSSFeed struct {
  Channel struct {
    Title       string
    Link        string
    Description string
    Language    string
    Items       []RSSItem
  }
}
```

---

#### `RSSItem`

**Location:** `rss.go`

```go
type RSSItem struct {
  Title       string
  Link        string
  Description string
  PubDate     string
}
```

---

### Request Models

#### `createUserParams`

```json
{
  "name": "username"
}
```

---

#### `createFeedParams`

```json
{
  "name": "Feed Name",
  "url": "https://example.com/rss.xml"
}
```

---

#### `createFeedFollowParams`

```json
{
  "feed_id": "uuid"
}
```

---

## Function Call Graph

### User Registration Flow

```
POST /v1/users
  ├─ createUserHandler()
  ├─ Decode JSON body
  ├─ CreateUser() [database]
  └─ respondWithJSON()
```

### Authenticated Request Flow

```
GET /v1/users (with Authorization header)
  ├─ authMiddleware()
  ├─ GetAPIKey() [auth package]
  ├─ GetUserByAPIKey() [database]
  ├─ handlerGetUser()
  └─ respondWithJSON()
```

### Feed Scraping Flow

```
startScraping() [background goroutine]
  ├─ GetNextFeedsToFetch() [database]
  ├─ scrapeFeed() [for each feed, concurrent]
  ├─ MarkFeedAsFetched() [database]
  ├─ urlToFeed()
  ├─ CreatePost() [database, for each item]
  └─ Log results
```

---

## Summary Table

| Function                | Type       | Auth | Purpose                   |
| ----------------------- | ---------- | ---- | ------------------------- |
| createUserHandler       | Handler    | No   | Create new user           |
| handlerGetUser          | Handler    | Yes  | Get authenticated user    |
| getPostsForUserHandler  | Handler    | Yes  | Get user's feed posts     |
| createFeedHandler       | Handler    | Yes  | Create feed subscription  |
| getFeedsHandler         | Handler    | No   | List all feeds            |
| createFeedFollowHandler | Handler    | Yes  | Follow a feed             |
| getFeedFollowsHandler   | Handler    | Yes  | Get user's followed feeds |
| readinessHandler        | Handler    | No   | Health check              |
| errorHandler            | Handler    | No   | Error testing             |
| authMiddleware          | Middleware | -    | Authenticate requests     |
| startScraping           | Scraper    | -    | Background feed fetcher   |
| scrapeFeed              | Scraper    | -    | Fetch single feed         |
| urlToFeed               | Utility    | -    | Parse RSS XML             |
| respondWithError        | Utility    | -    | Send error response       |
| respondWithJSON         | Utility    | -    | Send JSON response        |

---

## Notes

- All timestamps are in UTC timezone
- All UUIDs are randomly generated using `github.com/google/uuid`
- API keys are generated using SHA256 hashing
- Database operations use context for cancellation support
- Concurrent scraping uses `sync.WaitGroup` for synchronization
