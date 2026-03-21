package main

import (
	"time"

	"github.com/Wyden13/rssagg/db"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	// Email     string    `json:"email"`
	// Password  string    `json:"password"`
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
