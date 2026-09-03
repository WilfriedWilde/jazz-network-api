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
	organizer, err := s.repository.Read(ctx, s.repository.DB, id)
	if err != nil {
		return model.Organizer{}, convertToServiceError("read", "organizer", err)
	}

	return organizer, nil
}

func (s *OrganizerService) ReadAll(ctx context.Context) ([]model.Organizer, error) {
	organizers, err := s.repository.ReadAll(ctx, s.repository.DB)
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

	user, err := s.repository.Read

	organizer, err := s.repository.Create(ctx, tx, input.ToOrganizerCreate())
	if err != nil {
		return model.Organizer{}, convertToServiceError("create", "organizer", err)
	}

	isAdmin := user.AccessLevel == 1

	if !isAdmin {
		if err := s.repository.CreateUserOrganizerLink(ctx, tx, model.UserOrganizerCreate{
			UserID:      input.UserID,
			OrganizerID: organizer.ID,
		}); err != nil {
			return model.Organizer{}, convertToServiceError("create", "user_organizer", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return model.Organizer{}, fmt.Errorf("commit transaction: %w", err)
	}

	return organizer, nil
}
