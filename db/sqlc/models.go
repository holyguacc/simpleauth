// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"database/sql"
	"time"
)

type Post struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	PostDescription string    `json:"post_description"`
	AuthorName      string    `json:"author_name"`
	PostDate        time.Time `json:"post_date"`
}

type User struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type Verification struct {
	Username   string         `json:"username"`
	IsVerified sql.NullBool   `json:"is_verified"`
	VerifyKey  string         `json:"verify_key"`
	VerefiedOn time.Time      `json:"verefied_on"`
	ResetKey   sql.NullString `json:"reset_key"`
	ResetOn    time.Time      `json:"reset_on"`
	IsReset    sql.NullBool   `json:"is_reset"`
}
