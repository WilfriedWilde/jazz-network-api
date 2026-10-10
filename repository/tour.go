package repository

import "database/sql"

type TourRepository struct {
	DB *sql.DB
}

func NewTourRepository(db *sql.DB) *TourRepository {
	return &TourRepository{DB: db}
}
