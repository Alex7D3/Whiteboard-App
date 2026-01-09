package model

import (
	"time"
)

type User struct {
	ID           int64     `json:"id"         db:"id"`
	Username     string    `json:"username"   db:"username"`
	Email        string    `json:"email"      db:"email"`
	PasswordHash string    `json:"-"          db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type LoginResponse struct {
	User        *User  `json:"user"`
	AccessToken string `json:"access_token"`
}
