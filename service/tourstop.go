package service

import (
	"context"
	"fmt"
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

	tourExists, err := s.repository.Exists(ctx, tx, tourStopID)
	if err != nil {
		return model.TourStop{}, convertToServiceError("check existence", "tour_stop", err)
	}

	if !tourExists {
		return model.TourStop{}, ErrNotFound
	}

	if update.OrganizerID != nil {
		organizerExists, err := s.repository.ExistsOrganizer(ctx, tx, *update.OrganizerID)
		if err != nil {
			return model.TourStop{}, convertToServiceError("check existence", "tour_stop organizer", err)
		}

		if !organizerExists {
			return model.TourStop{}, ErrNotFound
		}
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
