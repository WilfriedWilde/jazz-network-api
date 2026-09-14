package service

import (
	"context"
	"fmt"
	m "jazz-network-api/middleware"
	"jazz-network-api/model"
	"jazz-network-api/repository"
)

type OrganizerService struct {
	repository *repository.OrganizerRepository
}

func NewOrganizerService(r *repository.OrganizerRepository) *OrganizerService {
	return &OrganizerService{repository: r}
}

func (s *OrganizerService) Read(ctx context.Context, organizerID int64) (model.Organizer, error) {
	user, ok := m.UserFromContext(ctx)
	if !ok {
		return model.Organizer{}, fmt.Errorf("get user from context")
	}

	isAdmin := user.AccessLevel == adminAccessLevel

	organizer, err := s.repository.Read(ctx, s.repository.DB, organizerID, user.ID, isAdmin)
	if err != nil {
		return model.Organizer{}, convertToServiceError("read", "organizer", err)
	}

	return organizer, nil
}

func (s *OrganizerService) ReadAll(ctx context.Context) ([]model.Organizer, error) {
	user, ok := m.UserFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("get user from context")
	}

	isAdmin := user.AccessLevel == adminAccessLevel

	organizers, err := s.repository.ReadAll(ctx, s.repository.DB, user.ID, isAdmin)
	if err != nil {
		return nil, convertToServiceError("read all", "organizers", err)
	}

	return organizers, nil
}

func (s *OrganizerService) Create(ctx context.Context, input model.OrganizerInput) (model.Organizer, error) {
	if err := validateOrganizerInput(input); err != nil {
		return model.Organizer{}, err
	}

	user, ok := m.UserFromContext(ctx)
	if !ok {
		return model.Organizer{}, fmt.Errorf("get user from context")
	}

	createOrganizer := input.ToOrganizerCreate()
	if user.AccessLevel != adminAccessLevel {
		createOrganizer.UserID = &user.ID
	}

	organizer, err := s.repository.Create(ctx, s.repository.DB, createOrganizer)
	if err != nil {
		return model.Organizer{}, convertToServiceError("create", "organizer", err)
	}

	return organizer, nil
}
