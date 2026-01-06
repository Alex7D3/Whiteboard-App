package model

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID               uuid.UUID `db:"id"`
	UserID           int64     `db:"user_id"`
	RefreshTokenHash string    `db:"refresh_token_hash"`
	ExpiresAt        time.Time `db:"expires_at"`
	CreatedAt        time.Time `db:"created_at"`
}
