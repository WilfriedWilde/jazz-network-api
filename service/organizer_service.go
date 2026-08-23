package service

import (
	"context"
	"fmt"
	"jazz-network-api/model"
	"jazz-network-api/repository"
)

type OrganizerService struct {
	repository *repository.OrganizerRepository
}

func NewOrganizerService(r *repository.OrganizerRepository) *OrganizerService {
	return &OrganizerService{repository: r}
}

func (s *OrganizerService) Read(ctx context.Context, id int64) (model.Organizer, error) {
	organizer, err := s.repository.Read(ctx, id)
	if err != nil {
		return model.Organizer{}, formatServiceError(err, "organizer")
	}

	return organizer, nil
}

func (s *OrganizerService) ReadAll(ctx context.Context) ([]model.Organizer, error) {
	organizers, err := s.repository.ReadAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("read all organizers: %w", err)
	}

	return organizers, nil
}
