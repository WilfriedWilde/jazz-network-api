package repository

import (
	"context"
	"database/sql"
	"fmt"
	"jazz-network-api/model"
)

const (
	createTourQuery = `
		INSERT INTO tours (user_id, name, start, end)
		VALUES(?, ?, ?, ?)
		RETURNING id, name, start, end
	`

	createToursTourStopsQuery = `
		INSERT INTO tour_stops (tour_id, type, stop_order)
		VALUES(?, ?, ?)
		RETURNING *
	`
)

type TourRepository struct {
	DB *sql.DB
}

func NewTourRepository(db *sql.DB) *TourRepository {
	return &TourRepository{DB: db}
}

func (r *TourRepository) Create(ctx context.Context, db DBTX, t model.TourCreate) (model.Tour, error) {
	var tour model.Tour

	if err := db.QueryRowContext(ctx, createTourQuery, t.UserID, t.Name, t.Start, t.End).Scan(
		&tour.ID, &tour.Name, &tour.Start, &tour.End,
	); err != nil {
		return model.Tour{}, fmt.Errorf("create tour: %w", err)
	}

	return tour, nil
}

func (r *TourRepository) CreateTourStops(ctx context.Context, tx *sql.Tx, numTourStops int, tourID int64) ([]model.TourStop, error) {
	tourStops := make([]model.TourStop, numTourStops)

	stmt, err := tx.PrepareContext(ctx, createToursTourStopsQuery)
	if err != nil {
		return nil, fmt.Errorf("prepare tx statement: %w", err)
	}
	defer stmt.Close()

	for i := range numTourStops {
		var stop model.TourStop

		if err := stmt.QueryRowContext(ctx, tourID, "free_day", i).Scan(
			&stop.ID, &stop.TourID, &stop.OrganizerID,
			&stop.Type, &stop.StopOrder, &stop.Notes,
		); err != nil {
			return nil, fmt.Errorf("scan tour_stop: %w", err)
		}

		tourStops[i] = stop
	}

	return tourStops, nil
}
