package service

import (
	"jazz-network-api/repository"
)

type SessionService struct {
	repository *repository.SessionRepository
}

func NewSessionService(r *repository.SessionRepository) *SessionService {
	return &SessionService{repository: r}
}
