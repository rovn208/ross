// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Follow struct {
	FollowingUserID int64     `json:"following_user_id"`
	FollowedUserID  int64     `json:"followed_user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	FullName       string    `json:"full_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Video struct {
	ID           int64       `json:"id"`
	Title        string      `json:"title"`
	StreamUrl    string      `json:"stream_url"`
	Description  pgtype.Text `json:"description"`
	ThumbnailUrl pgtype.Text `json:"thumbnail_url"`
	CreatedBy    int64       `json:"created_by"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
