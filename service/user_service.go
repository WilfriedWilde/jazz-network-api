package service

import (
	"context"
	"fmt"
	"jazz-network-api/model"
	"jazz-network-api/repository"

	"golang.org/x/crypto/bcrypt"
)

const userLevel = 0

var dummyHash = []byte("$2a$10$CwTycUXWue0Thq9StjUM0uJ8vHZlHUqR8DUoxfawR9.p0AR5xLDVy")

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repository: r}
}

func (s *UserService) Create(ctx context.Context, input model.UserInput) (model.User, error) {
	tx, err := s.repository.DB.BeginTx(ctx, nil)
	if err != nil {
		return model.User{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	if existing, _ := s.repository.GetByUsername(ctx, tx, input.Username); existing != nil {
		return model.User{}, convertToServiceError("read by username", "user", ErrUsernameTaken)
	}

	if existing, _ := s.repository.GetByEmail(ctx, tx, input.Email); existing != nil {
		return model.User{}, convertToServiceError("read by email", "user", ErrEmailTaken)
	}

	if err := validatePassword(input.Password); err != nil {
		return model.User{}, convertToServiceError("validate", "password", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, convertToServiceError("hash", "password", err)
	}

	user, err := s.repository.Create(ctx, tx, model.UserCreate{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(passwordHash),
		AccessLevel:  userLevel,
	})
	if err != nil {
		return model.User{}, convertToServiceError("create", "user", err)
	}

	if err := tx.Commit(); err != nil {
		return model.User{}, fmt.Errorf("commit transaction: %w", err)
	}

	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, input model.UserInput) (model.User, error) {
	user, err := s.repository.GetByEmail(ctx, s.repository.DB, input.Email)
	if err != nil {
		bcrypt.CompareHashAndPassword(dummyHash, []byte(input.Password))
		return model.User{}, convertToServiceError("authenticate", "user", ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return model.User{}, convertToServiceError("authenticate", "user", ErrInvalidCredentials)
	}

	return *user, nil
}
