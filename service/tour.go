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

	// TODO: everything should be in tx
	// first create tour, then corresponding tourstops
	user, ok := m.UserFromContext(ctx)
	if !ok {
		return model.Tour{}, fmt.Errorf("get user from context")
	}

	// TODO: implement ToTourCreate()
	createTour := input.ToTourCreate()
	if user.AccessLevel != adminAccessLevel {
		createTour.UserID = &user.ID
	}

	tour, err := s.repository.Create(ctx, s.repository.DB, createTour)
	if err != nil {
		return model.Tour{}, convertToServiceError("create", "tour", err)
	}

	return tour, nil
}
