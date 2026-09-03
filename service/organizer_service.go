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
		return model.Organizer{}, convertToServiceError("read", "organizer", err)
	}

	return organizer, nil
}

func (s *OrganizerService) ReadAll(ctx context.Context) ([]model.Organizer, error) {
	organizers, err := s.repository.ReadAll(ctx)
	if err != nil {
		return nil, convertToServiceError("read all", "organizers", err)
	}

	return organizers, nil
}

func (s *OrganizerService) Create(ctx context.Context, input model.OrganizerCreateInput) (model.Organizer, error) {
	tx, err := s.repository.DB.BeginTx(ctx, nil)
	if err != nil {
		return model.Organizer{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	organizer, err := s.repository.Create(ctx, input.ToOrganizerCreate())
	if err != nil {
		convertToServiceError("create", "organizer", err)
	}

	tx.Commit()

	return organizer, nil
}
