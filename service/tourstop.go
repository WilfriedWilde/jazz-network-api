package service

import "jazz-network-api/repository"

type TourStopService struct {
	repository *repository.TourStopRepository
}

func NewTourStopService(r *repository.TourStopRepository) *TourStopService {
	return &TourStopService{repository: r}
}
