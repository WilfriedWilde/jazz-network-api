package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"jazz-network-api/model"

	"github.com/mattn/go-sqlite3"
)

const (
	createUserQuery = `
		INSERT INTO users(username, email, password_hash, access_level)
		VALUES (?, ?, ?, ?)
		RETURNING *
	`

	readUserByUsername = `
		SELECT *
		FROM users
		WHERE username = ?
	`

	readUserByEmail = `
		SELECT *
		FROM users
		WHERE email = ?
	`
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, db DBTX, u model.UserCreate) (model.User, error) {
	var user model.User

	if err := db.QueryRowContext(
		ctx, createUserQuery,
		u.Username, u.Email, u.PasswordHash, u.AccessLevel,
	).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.PasswordHash, &user.AccessLevel,
	); err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			return model.User{}, ErrDuplicate
		}

		return model.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, db DBTX, username string) (*model.User, error) {
	var user model.User

	err := db.QueryRowContext(
		ctx, readUserByUsername, username,
	).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.PasswordHash, &user.AccessLevel,
	)

	return &user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, db DBTX, email string) (*model.User, error) {
	var user model.User

	err := db.QueryRowContext(
		ctx, readUserByEmail, email,
	).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.PasswordHash, &user.AccessLevel,
	)

	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}

	return &user, err
}
