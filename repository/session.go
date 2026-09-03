package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	getUserIDBySessionQuery = `
		SELECT user_id
		FROM sessions
		WHERE uuid = ?`
)

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (s *SessionRepository) GetUserIDBySession(ctx context.Context, db DBTX, uuid string) (int64, error) {
	var userID int64

	if err := db.QueryRowContext(ctx, getUserIDBySessionQuery, uuid).Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userID, ErrNotFound
		}

		return userID, fmt.Errorf("read user id by session: %w", err)
	}

	return userID, nil
}
