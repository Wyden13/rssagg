package main

import (
	"time"

	"github.com/Wyden13/rssagg/db"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseUserToAPIUser(dbUser db.User) User {
	return User{
		ID:        dbUser.ID.String(),
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}

type Feed struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    string    `json:"user_id"`
}

func databaseFeedToAPIFeed(dbFeed db.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID.String(),
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		UserID:    dbFeed.UserID.String(),
	}
}
func databaseFeedsToAPIFeeds(dbFeeds []db.Feed) []Feed {
	feeds := []Feed{}
	for _, feed := range dbFeeds {
		feeds = append(feeds, databaseFeedToAPIFeed(feed))
	}
	return feeds
}
