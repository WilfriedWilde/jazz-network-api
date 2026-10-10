package service

import "jazz-network-api/repository"

type TourService struct {
	repository *repository.TourRepository
}

func NewTourService(r *repository.TourRepository) *TourService {
	return &TourService{repository: r}
}
