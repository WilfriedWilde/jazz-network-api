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

	tour.Stops, err = s.repository.CreateTourStops(ctx, tx, numTourStops, tour.ID)
	if err != nil {
		return model.Tour{}, convertToServiceError("create", "tour's tour_stops", err)
	}

	if err := tx.Commit(); err != nil {
		return model.Tour{}, fmt.Errorf("commit transaction")
	}

	return tour, nil
}
