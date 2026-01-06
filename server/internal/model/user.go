package model

import (
	"time"
)

type User struct {
	ID           int64     `json:"id"         db:"id"`
	UserName     string    `json:"username"   db:"username"`
	Email        string    `json:"email"      db:"email"`
	PasswordHash string    `json:"-"          db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewUser(username string, email string, passwordhash string) *User {
	return &User{
		UserName: username,
		Email: email,
		PasswordHash: passwordhash,
	}
}

type UserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID          int64  `json:"id"`
	UserName    string `json:"username"`
	Email       string `json:"email"`
}

type LoginResponse struct {
	User        *UserResponse `json:"user"`
	AccessToken string       `json:"access_token"`
}

func NewUserResponse(id int64, username string, email string) *UserResponse {
	return &UserResponse{
		ID: id,
		UserName: username,
		Email: email,
	}
}

func NewLoginResponse(id int64, username string, email string, accessToken string) *LoginResponse {
	return &LoginResponse{
		User: NewUserResponse(id, username, email),
		AccessToken: accessToken,
	}
}
