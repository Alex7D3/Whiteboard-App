package storage

import (
	"drawing-api/internal/model"
	"context"
)

type UserStorage interface {
	Create(ctx context.Context, user *model.User) (int, error)
	GetById(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}
