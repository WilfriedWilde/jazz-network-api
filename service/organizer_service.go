package service

import (
	"context"
	"fmt"
	"jazz-network-api/model"
	"jazz-network-api/repository"
	r "jazz-network-api/repository"
)

type OrganizerService struct {
	repository *r.OrganizerRepository
}

func NewOrganizerService(r *repository.OrganizerRepository) *OrganizerService {
	return &OrganizerService{repository: r}
}

func (s *OrganizerService) ReadAll(ctx context.Context) ([]model.Organizer, error) {
	organizers, err := s.repository.ReadAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("read all organizers: %w", err)
	}

	return organizers, nil
}
