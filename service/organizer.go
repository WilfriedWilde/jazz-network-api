package service

import (
	"context"
	"fmt"
	m "jazz-network-api/middleware"
	"jazz-network-api/model"
	"jazz-network-api/repository"
	"jazz-network-api/validation"
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
	if err := validation.ValidateOrganizerInput(input); err != nil {
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

func (s *OrganizerService) Delete(ctx context.Context, organizerID int64, force bool) error {
	tx, err := s.repository.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx delete organizer: %w", err)
	}
	defer func() error {
		return tx.Rollback()
	}()

	// check if organizer is referenced by tour_stop
	tourStopIDs, err := s.repository.ReadTourStopIDsByOrganizerID(ctx, tx, organizerID)
	if err != nil {
		return convertToServiceError("read", "tour_stops IDs by organizer ID", err)
	}

	if len(tourStopIDs) > 0 {
		// on delete force:
		// reset referencing tour_stops to organizer_id = NULL and type = free_day
		// delete tour_stop referencers: tour_stop_finances, travels and accommodations
		if force {
			if err := s.repository.ResetOrganizerTourStops(ctx, tx, tourStopIDs); err != nil {
				return convertToServiceError("update", "tour_stops (organizer)", err)
			}

			if err := s.repository.DeleteOrganizerTourStopFinances(ctx, tx, tourStopIDs); err != nil {
				return convertToServiceError("delete", "tour_stop_finances by tour_stops IDs (organizer)", err)
			}

			if err := s.repository.DeleteOrganizerTravels(ctx, tx, tourStopIDs); err != nil {
				return convertToServiceError("delete", "travels by tour_stops IDs (organizer)", err)
			}

			if err := s.repository.DeleteOrganizerAccommodations(ctx, tx, tourStopIDs); err != nil {
				return convertToServiceError("delete", "accommodations by tour_stops IDs (organizer)", err)
			}
		} else {
			return ErrOrganizerConflict
		}

	}

	// delete organizer
	if err := s.repository.Delete(ctx, tx, organizerID); err != nil {
		return convertToServiceError("delete", "organizer", err)
	}

	return tx.Commit()
}

func (s *OrganizerService) Update(ctx context.Context, id int64, update model.OrganizerUpdate) (model.Organizer, error) {
	if err := validation.ValidateOrganizerUpdate(update); err != nil {
		return model.Organizer{}, err
	}

	exists, err := s.repository.Exists(ctx, s.repository.DB, id)
	if err != nil {
		return model.Organizer{}, convertToServiceError("check existence", "organizer", err)
	}

	if !exists {
		return model.Organizer{}, ErrNotFound
	}

	organizer, err := s.repository.Update(ctx, s.repository.DB, id, update)
	if err != nil {
		return model.Organizer{}, convertToServiceError("update", "organizer", err)
	}

	return organizer, nil
}
