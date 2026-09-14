package service

import (
	"context"
	"jazz-network-api/model"
	"jazz-network-api/repository"
	"time"
	"uuid"
)

type SessionService struct {
	repository *repository.SessionRepository
}

const SessionExpiryTime = time.Hour * 24

func NewSessionService(r *repository.SessionRepository) *SessionService {
	return &SessionService{repository: r}
}

func (s *SessionService) Create(ctx context.Context, userID int64) (model.Session, error) {
	session := model.Session{
		UUID:      uuid.New().String(),
		UserID:    userID,
		CreatedAt: time.Now().Format(time.DateTime),
		ExpiresAt: time.Now().Add(SessionExpiryTime).Format(time.DateTime),
	}

	return s.repository.Create(ctx, session)
}

func (s *SessionService) Delete(ctx context.Context, sessionUUID string) error {
	if err := s.repository.Delete(ctx, sessionUUID); err != nil {
		return convertToServiceError("delete", "session", err)
	}

	return nil
}
