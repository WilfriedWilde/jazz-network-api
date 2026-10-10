package repository

import "database/sql"

type TourStopRepository struct {
	DB *sql.DB
}

func NewTourStopRepository(db *sql.DB) *TourStopRepository {
	return &TourStopRepository{DB: db}
}
