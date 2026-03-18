-- name: CreateUser :one
INSERT INTO users (id, username, email, password_hash, created_at, updated_at) 
VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;