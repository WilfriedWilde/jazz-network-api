package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"jazz-network-api/model"
	"time"
)

const (
	createSessionQuery = `
		INSERT INTO sessions(uuid, user_id, created_at, expires_at)
		VALUES(?, ?, ?, ?)
		RETURNING *
	`

	getUserBySessionQuery = `
		SELECT u.id, u.username, u.email, u.password_hash, u.access_level
		FROM sessions AS s
		JOIN users AS u ON u.id = s.user_id
		WHERE s.uuid = ? AND s.expires_at > ?
	`

	deleteSessionQuery = `
		DELETE FROM sessions
		WHERE uuid = ?
	`
)

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) Create(ctx context.Context, s model.Session) (model.Session, error) {
	var session model.Session

	err := r.DB.QueryRowContext(
		ctx, createSessionQuery,
		s.UUID,
		s.UserID,
		s.CreatedAt,
		s.ExpiresAt,
	).Scan(
		&session.UUID,
		&session.UserID,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	return session, err
}

func (r *SessionRepository) GetUserBySession(ctx context.Context, db DBTX, sessionUUID string) (model.User, error) {
	var user model.User

	if err := db.QueryRowContext(ctx, getUserBySessionQuery, sessionUUID, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.AccessLevel,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrNotFound
		}

		return model.User{}, fmt.Errorf("get user by session: %w", err)
	}

	return user, nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionUUID string) error {
	if _, err := r.DB.ExecContext(ctx, deleteSessionQuery, sessionUUID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}

		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
