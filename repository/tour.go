package repository

import (
	"context"
	"database/sql"
	"errors"
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

	readAllToursQuery = `
		SELECT id, name, start, end
		FROM tours
		WHERE user_id = ?
		ORDER BY id ASC
	`

	readAllTourStopsByUserIDQuery = `
		SELECT *
		FROM tour_stops
		WHERE user_id = ?
		ORDER BY id ASC
	`

	readTourByIDQuery = `
		SELECT id, name, start, end
		FROM tours
		WHERE id = ?
	`

	readTourStopsByTourIDQuery = `
		SELECT *
		FROM tour_stops
		WHERE tour_id = ?
		ORDER BY id ASC
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

func (r *TourRepository) ReadAll(ctx context.Context, db DBTX, userID int64) ([]model.Tour, error) {
	rows, err := db.QueryContext(ctx, readAllToursQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("read all tours: %w", err)
	}
	defer rows.Close()

	tours := make([]model.Tour, 0)

	for rows.Next() {
		t := model.Tour{}

		if err := rows.Scan(
			&t.ID, &t.Name, &t.Start, &t.End,
		); err != nil {
			return nil, fmt.Errorf("scan tour: %w", err)
		}

		tours = append(tours, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate over tours rows: %w", err)
	}

	return tours, nil
}

func (r *TourRepository) ReadAllTourStopsByUserID(ctx context.Context, db DBTX, userID int64) ([]model.TourStop, error) {
	rows, err := db.QueryContext(ctx, readAllTourStopsByUserIDQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("read all tour_stops by user ID: %w", err)
	}
	defer rows.Close()

	stops := make([]model.TourStop, 0)

	for rows.Next() {
		s := model.TourStop{}

		if err := rows.Scan(
			&s.ID, &s.TourID, &s.OrganizerID,
			&s.Type, &s.StopOrder, &s.Notes,
		); err != nil {
			return nil, fmt.Errorf("scan tour_stop: %w", err)
		}

		stops = append(stops, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate over tour_stops rows: %w", err)
	}

	return stops, nil
}

func (r *TourRepository) Read(ctx context.Context, db DBTX, tourID int64) (model.Tour, error) {
	var tour model.Tour

	if err := db.QueryRowContext(ctx, readTourByIDQuery, tourID).Scan(
		&tour.ID, &tour.Name, &tour.Start, &tour.End,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Tour{}, ErrNotFound
		}

		return model.Tour{}, fmt.Errorf("read tour by id: %w", err)
	}

	return tour, nil
}

func (r *TourRepository) ReadTourStopsByTourID(ctx context.Context, db DBTX, tourID int64) ([]model.TourStop, error) {
	rows, err := db.QueryContext(ctx, readTourStopsByTourIDQuery, tourID)
	if err != nil {
		return nil, fmt.Errorf("read tour_stops by tour ID: %w", err)
	}
	defer rows.Close()

	stops := make([]model.TourStop, 0)

	for rows.Next() {
		s := model.TourStop{}

		if err := rows.Scan(
			&s.ID, &s.TourID, &s.OrganizerID,
			&s.Type, &s.StopOrder, &s.Notes,
		); err != nil {
			return nil, fmt.Errorf("scan tour_stop: %w", err)
		}

		stops = append(stops, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate over tour_stops rows: %w", err)
	}

	return stops, nil
}
