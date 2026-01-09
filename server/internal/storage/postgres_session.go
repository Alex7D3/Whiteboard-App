package storage

import (
	"drawing-api/internal/model"
	"drawing-api/internal/util"
	"github.com/jmoiron/sqlx"
	"database/sql"
	"fmt"
	"time"
	"context"
)

type PGSessionStorage struct {
	db *sqlx.DB
}

func NewPGSessionStorage(dbx *sqlx.DB) *PGSessionStorage {
	return &PGSessionStorage{db: dbx}
}


func (s *PGSessionStorage) Create(ctx context.Context, session *model.Session) (string, error) {
	var id string
	const query = "INSERT INTO sessions (token_hash, user_id,  expires_at) VALUES ($1, $2, $3) RETURNING token_hash"
	err := s.db.QueryRowContext(ctx, query,
		 session.TokenHash, session.UserID, session.ExpiresAt,
	).Scan(&id)
	return id, err
}

func (s *PGSessionStorage) GetByUserID(ctx context.Context, id int64) (*model.Session, error) {
	var session model.Session
	err := s.db.GetContext(ctx, &session, "SELECT * FROM sessions WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *PGSessionStorage) Revoke(ctx context.Context, token string) error {
	hash := util.HashToken(token)
	_, err := s.db.ExecContext(ctx, "UPDATE sessions SET revoked = TRUE WHERE token_hash = $1", hash)
	return err
}

func (s *PGSessionStorage) RotateToken(ctx context.Context, expiry time.Duration, oldToken, newToken string) (*model.Session, error) {
	errorStr := "Failed to rotate token: %v"

	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	defer tx.Rollback()

	if err != nil {
		return nil, fmt.Errorf(errorStr, err)
	}

	oldHash := util.HashToken(oldToken)
	fmt.Println(oldHash)

	// Find the session and lock down the record
	var session model.Session
	if err = tx.GetContext(ctx, &session, `
		SELECT *
		FROM sessions
		WHERE token_hash = $1 
		FOR UPDATE
	`, oldHash);
	err == sql.ErrNoRows || session.Revoked || time.Now().After(session.ExpiresAt) {
		fmt.Println(err)
		fmt.Println(session)
		return nil, fmt.Errorf(errorStr, "token has expired or was revoked")
	} else if err != nil {
		return nil, fmt.Errorf(errorStr, err)
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE sessions
		SET revoked = TRUE
		WHERE token_hash = $1
	`, oldHash);
	err != nil {
		return nil, fmt.Errorf(errorStr, err)
	}

	newHash := util.HashToken(newToken)
	if _, err = tx.ExecContext(ctx, `
		INSERT INTO sessions (token_hash, user_id,  expires_at)
		VALUES ($1, $2, $3) RETURNING token_hash
	`, newHash, session.UserID, time.Now().Add(expiry));
	err != nil {
		return nil, fmt.Errorf(errorStr, err)
	}


	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(errorStr, err)
	}

	return &session, nil
}
