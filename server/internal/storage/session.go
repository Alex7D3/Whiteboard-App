package storage

import (
	"drawing-api/internal/model"
	"github.com/google/uuid"
	"context"
)

type SessionStorage interface {
	Create(ctx context.Context, session *model.Session) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Session, error)
	GetByUserID(ctx context.Context, id int64) (*model.Session, error)
}
