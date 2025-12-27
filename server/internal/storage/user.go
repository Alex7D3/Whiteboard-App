package storage

import (
	"drawing-api/internal/model"
	"context"
)

type UserStorage interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}
