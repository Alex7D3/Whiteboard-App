package storage

import (
	"drawing-api/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
	"context"
)

type PGSessionStorage struct {
	db *sqlx.DB
}

func NewPGSessionStorage(dbx *sqlx.DB) *PGSessionStorage {
	return &PGSessionStorage{db: dbx}
}


func (s *PGSessionStorage) Create(ctx context.Context, session *model.Session) (uuid.UUID, error) {
	var id uuid.UUID 
	const query = "INSERT INTO sessions (id, user_id, refresh_token_hash, expires_at) VALUES ($1, $2, $3, $4) RETURNING id"
	err := s.db.QueryRowContext(ctx, query,
		session.ID, session.UserID, session.RefreshTokenHash, session.ExpiresAt,
	).Scan(&id)
	return id, err
}

func (s *PGSessionStorage) GetByID(ctx context.Context, id uuid.UUID) (*model.Session, error) {
	var session model.Session
	err := s.db.GetContext(ctx, &session, "SELECT * FROM sessions WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *PGSessionStorage) GetByUserID(ctx context.Context, id int64) (*model.Session, error) {
	var session model.Session
	err := s.db.GetContext(ctx, &session, "SELECT * FROM sessions WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
