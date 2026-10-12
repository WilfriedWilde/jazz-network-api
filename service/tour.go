package service

import (
	"context"
	"fmt"
	m "jazz-network-api/middleware"
	"jazz-network-api/model"
	"jazz-network-api/repository"
	"jazz-network-api/validation"
)

type TourService struct {
	repository *repository.TourRepository
}

func NewTourService(r *repository.TourRepository) *TourService {
	return &TourService{repository: r}
}

func (s *TourService) Create(ctx context.Context, input model.TourInput) (model.Tour, error) {
	if err := validation.ValidateTourInput(input); err != nil {
		return model.Tour{}, err
	}

	tx, err := s.repository.DB.BeginTx(ctx, nil)
	if err != nil {
		return model.Tour{}, fmt.Errorf("begin transaction")
	}
	defer func() error {
		return tx.Rollback()
	}()

	user, ok := m.UserFromContext(ctx)
	if !ok {
		return model.Tour{}, fmt.Errorf("get user from context")
	}

	tourCreate := model.TourCreate{
		UserID: user.ID,
		Name:   input.Name,
		Start:  input.Start,
		End:    input.End,
	}

	tour, err := s.repository.Create(ctx, tx, tourCreate)
	if err != nil {
		return model.Tour{}, convertToServiceError("create", "tour", err)
	}

	numTourStops := getNumTourStops(tour)

	tour.Stops, err = s.repository.CreateTourStops(ctx, tx, numTourStops, model.TourStopCreate{
		TourID: tour.ID,
		Type:   "unplanned",
		Date:   tour.Start,
	})
	if err != nil {
		return model.Tour{}, convertToServiceError("create", "tour's tour_stops", err)
	}

	if err := tx.Commit(); err != nil {
		return model.Tour{}, fmt.Errorf("commit transaction")
	}

	return tour, nil
}

func (s *TourService) ReadAll(ctx context.Context) ([]model.Tour, error) {
	user, ok := m.UserFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("get user from context")
	}

	tours, err := s.repository.ReadAll(ctx, s.repository.DB, user.ID)
	if err != nil {
		return nil, convertToServiceError("read all", "tours", err)
	}

	tourStops, err := s.repository.ReadAllTourStopsByUserID(ctx, s.repository.DB, user.ID)
	if err != nil {
		return nil, convertToServiceError("read all", "tour_stops by user ID", err)
	}

	byTourID := make(map[int64][]model.TourStop, len(tours))
	for _, stop := range tourStops {
		byTourID[stop.TourID] = append(byTourID[stop.TourID], stop)
	}

	for i := range tours {
		tours[i].Stops = byTourID[tours[i].ID]
	}

	return tours, nil
}

func (s *TourService) Read(ctx context.Context, tourID int64) (model.Tour, error) {
	tour, err := s.repository.Read(ctx, s.repository.DB, tourID)
	if err != nil {
		return model.Tour{}, convertToServiceError("read", "tour", err)
	}

	tour.Stops, err = s.repository.ReadTourStopsByTourID(ctx, s.repository.DB, tourID)
	if err != nil {
		return model.Tour{}, convertToServiceError("read", "tour_stops by tourID", err)
	}

	return tour, nil
}

func (s *TourService) Update(ctx context.Context, tourID int64, update model.TourUpdate) (model.Tour, error) {
	if err := validation.ValidateTourUpdate(update); err != nil {
		return model.Tour{}, err
	}

	tx, err := s.repository.DB.BeginTx(ctx, nil)
	if err != nil {
		return model.Tour{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() error {
		return tx.Rollback()
	}()

	exists, err := s.repository.Exists(ctx, tx, tourID)
	if err != nil {
		return model.Tour{}, convertToServiceError("check existence", "tour", err)
	}

	if !exists {
		return model.Tour{}, ErrNotFound
	}

	// if updating time range, add/remove corresponding tour_stops
	if update.Start != nil && update.End != nil {
		// delete stops outside of new range
		if err := s.repository.DeleteOutOfRangeTourStops(ctx, tx, tourID, update); err != nil {
			return model.Tour{}, convertToServiceError("delete", "out of range tour_stops", err)
		}

		// get missing dates and insert them
		stops, err := s.repository.ReadTourStopsByTourID(ctx, tx, tourID)
		if err != nil {
			return model.Tour{}, convertToServiceError("read", "tour_stops by tour ID", err)
		}

		existingDates := make([]string, len(stops))
		for i, s := range stops {
			existingDates[i] = s.Date
		}

		missingDates, err := getMissingDates(existingDates, *update.Start, *update.End)
		if err != nil {
			return model.Tour{}, fmt.Errorf("get missing dates: %w")
		}

		if len(missingDates) > 0 {
			if err := s.repository.CreateMissingTourStops(ctx, tx, tourID, missingDates); err != nil {
				return model.Tour{}, convertToServiceError("create", "missing tour_stops", err)
			}
		}
	}

	// update tour
	tour, err := s.repository.Update(ctx, tx, tourID, update)
	if err != nil {
		return model.Tour{}, convertToServiceError("update", "tour", err)
	}

	tour.Stops, err = s.repository.ReadTourStopsByTourID(ctx, tx, tourID)
	if err != nil {
		return model.Tour{}, convertToServiceError("read", "tour_stops by tourID", err)
	}

	if err := tx.Commit(); err != nil {
		return model.Tour{}, fmt.Errorf("commit transaction")
	}

	return tour, nil
}

func (s *TourService) Delete(ctx context.Context, tourID int64, force bool) error {
	tx, err := s.repository.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx delete tour: %w", err)
	}
	defer func() error {
		return tx.Rollback()
	}()

	// check if there is any tour's tour_stops with type != 'unplanned'
	plannedTourStops, err := s.repository.CheckPlannedTourStopsByTourID(ctx, tx, tourID)
	if err != nil {
		return convertToServiceError("check", "planned tour_stops by tour ID", err)
	}

	if plannedTourStops && !force {
		return ErrTourConflict
	}

	// delete tour
	if err := s.repository.Delete(ctx, tx, tourID); err != nil {
		return convertToServiceError("delete", "tour", err)
	}

	return tx.Commit()
}
