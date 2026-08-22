package service

import (
	"context"
	"fmt"
	"jazz-network-api/model"
	r "jazz-network-api/repository"
)

type OrganizerService struct {
	repository *r.OrganizerRepository
}

func (s *OrganizerService) ReadAllOrganizers(ctx context.Context) ([]model.Organizer, error) {
	organizers, err := s.repository.ReadAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("read all organizers: %w", err)
	}

	return organizers, nil
}
