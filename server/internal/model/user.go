package model

import (
	"time"
)

type User struct {
	ID           int       `json:"id"         db:"user_id"`
	UserName     string    `json:"username"   db:"username"`
	Email        string    `json:"email"      db:"email"`
	PasswordHash string    `json:"-"          db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewUser(username string, email, passwordhash string) *User {
	return &User{
		UserName: username,
		Email: email,
		PasswordHash: passwordhash,
	}
}

type UserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	ID           int    `json:"id"`
	UserName     string `json:"username"`
	Email        string `json:"email"`
}

func NewUserResponse(id int, username string, email string) *UserResponse {
	return &UserResponse{
		ID: id,
		UserName: username,
		Email: email,
	}
}
