package storage

import (
	"drawing-api/internal/model"
	"github.com/jmoiron/sqlx"
	"context"
)

type PGUserStorage struct {
	db *sqlx.DB
}

func NewPGUserStorage(dbx *sqlx.DB) *PGUserStorage {
	return &PGUserStorage{db: dbx}
}

func (s *PGUserStorage) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64
	const query = "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING user_id"
	err := s.db.QueryRowContext(ctx, query,
		user.UserName, user.Email, user.PasswordHash,
	).Scan(&id)
	return id, err
}

func (s *PGUserStorage) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PGUserStorage) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := s.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
