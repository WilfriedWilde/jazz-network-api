package service

import (
	"context"
	"fmt"
	m "jazz-network-api/middleware"
	"jazz-network-api/model"
	"jazz-network-api/repository"
)

type TourService struct {
	repository *repository.TourRepository
}

func NewTourService(r *repository.TourRepository) *TourService {
	return &TourService{repository: r}
}

func (s *TourService) Create(ctx context.Context, input model.TourInput) (model.Tour, error) {
	if err := validateTourInput(input); err != nil {
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
	for _, s := range tourStops {
		byTourID[s.TourID] = append(byTourID[s.TourID], s)
	}

	for _, t := range tours {
		t.Stops = byTourID[t.ID]
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
