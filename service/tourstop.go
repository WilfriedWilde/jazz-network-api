package service

import (
	"context"
	"fmt"
	m "jazz-network-api/middleware"
	"jazz-network-api/model"
	"jazz-network-api/repository"
	"jazz-network-api/validation"
)

type TourStopService struct {
	repository *repository.TourStopRepository
}

func NewTourStopService(r *repository.TourStopRepository) *TourStopService {
	return &TourStopService{repository: r}
}

func (s *TourStopService) Update(ctx context.Context, tourStopID int64, update model.TourStopUpdate) (model.TourStop, error) {
	if err := validation.ValidateTourStopUpdate(update); err != nil {
		return model.TourStop{}, err
	}

	tx, err := s.repository.DB.BeginTx(ctx, nil)
	if err != nil {
		return model.TourStop{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() error {
		return tx.Rollback()
	}()

	tourStopExists, err := s.repository.Exists(ctx, tx, tourStopID)
	if err != nil {
		return model.TourStop{}, convertToServiceError("check existence", "tour_stop", err)
	}

	if !tourStopExists {
		fmt.Println("tour not found")
		return model.TourStop{}, ErrNotFound
	}

	if update.OrganizerID != nil {
		user, ok := m.UserFromContext(ctx)
		if !ok {
			return model.TourStop{}, fmt.Errorf("get user from context")
		}

		isAdmin := user.AccessLevel == adminAccessLevel

		organizer, err := s.repository.ReadOrganizer(ctx, tx, *update.OrganizerID, user.ID, isAdmin)
		if err != nil {
			return model.TourStop{}, convertToServiceError("read", "tour_stop organizer", err)
		}

		*update.Country = organizer.Country
		*update.City = organizer.City
	}

	tourStop, err := s.repository.Update(ctx, tx, tourStopID, update)
	if err != nil {
		return model.TourStop{}, convertToServiceError("update", "tour_stop", err)
	}

	if err := tx.Commit(); err != nil {
		return model.TourStop{}, fmt.Errorf("commit transaction")
	}

	return tourStop, nil
}
