package storage

import (
	"drawing-api/internal/model"
	"context"
	"time"
)

type SessionStorage interface {
	Create(ctx context.Context, session *model.Session) (string, error)
	GetByUserID(ctx context.Context, id int64) (*model.Session, error)
	RotateToken(ctx context.Context, expiry time.Duration, oldToken, newToken string) (*model.Session, error)
	Revoke(ctx context.Context, token string) error
}
