package model

import (
	"time"
)

type Session struct {
	TokenHash string    `db:"token_hash"`
	UserID    int64     `db:"user_id"`
	Revoked   bool      `db:"revoked"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
