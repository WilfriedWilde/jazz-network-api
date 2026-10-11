package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"jazz-network-api/model"
	"strings"

	"github.com/mattn/go-sqlite3"
)

const (
	createTourQuery = `
		INSERT INTO tours (user_id, name, start_date, end_date)
		VALUES(?, ?, ?, ?)
		RETURNING id, name, start_date, end_date
	`

	createToursTourStopsQuery = `
		INSERT INTO tour_stops (tour_id, type, date)
		VALUES(?, ?, ?, ?)
		RETURNING *
	`

	readAllToursQuery = `
		SELECT id, name, start_date, end_date
		FROM tours
		WHERE user_id = ?
		ORDER BY id ASC
	`

	readAllTourStopsByUserIDQuery = `
		SELECT ts.id, ts.tour_id, ts.organizer_id, ts.type, ts.date, ts.country, ts.city, ts.notes
		FROM tour_stops AS ts
		JOIN tours AS t
		ON t.id = ts.tour_id
		WHERE t.user_id = ?
		ORDER BY ts.date ASC
	`

	readTourByIDQuery = `
		SELECT id, name, start_date, end_date
		FROM tours
		WHERE id = ?
	`

	readTourStopsByTourIDQuery = `
		SELECT *
		FROM tour_stops
		WHERE tour_id = ?
		ORDER BY date ASC
	`

	existTourQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM tours
			WHERE id = ?
		)
	`

	deleteOutOfRangeTourStopsQuery = `
		DELETE FROM tour_stops
		WHERE tour_id = ? AND (date < ? OR date > ?)
	`

	createMissingTourStopsQuery = `
		INSERT INTO tour_stops (tour_id, type, date)
		VALUES(?, ?, ?)
	`

	updateTourQuery = `
		UPDATE tours
		SET %s
		WHERE id = ?
		RETURNING id, name, start_date, end_date
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
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return model.Tour{}, ErrDuplicate
		}

		return model.Tour{}, fmt.Errorf("create tour: %w", err)
	}

	return tour, nil
}

func (r *TourRepository) CreateTourStops(ctx context.Context, tx *sql.Tx, numTourStops int, ts model.TourStopCreate) ([]model.TourStop, error) {
	tourStops := make([]model.TourStop, numTourStops)

	stmt, err := tx.PrepareContext(ctx, createToursTourStopsQuery)
	if err != nil {
		return nil, fmt.Errorf("prepare tx statement: %w", err)
	}
	defer stmt.Close()

	for i := range numTourStops {
		var s model.TourStop

		date, err := getStopDate(ts.Date, i)
		if err != nil {
			return nil, fmt.Errorf("get stop date: %w", err)
		}

		if err := stmt.QueryRowContext(ctx, ts.TourID, ts.Type, date).Scan(
			&s.ID, &s.TourID, &s.OrganizerID,
			&s.Type, &s.Date, &s.Country,
			&s.City, &s.Notes,
		); err != nil {
			return nil, fmt.Errorf("scan tour_stop: %w", err)
		}

		tourStops[i] = s
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
			&s.Type, &s.Date, &s.Country,
			&s.City, &s.Notes,
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
			&s.Type, &s.Date, &s.Country,
			&s.City, &s.Notes,
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

func (r *TourRepository) Exists(ctx context.Context, db DBTX, id int64) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, existTourQuery, id).Scan(&exists)
	return exists, err
}

func (r *TourRepository) DeleteOutOfRangeTourStops(ctx context.Context, tx *sql.Tx, tourID int64, t model.TourUpdate) error {
	if _, err := tx.ExecContext(ctx, deleteOutOfRangeTourStopsQuery, tourID, *t.Start, *t.End); err != nil {
		return fmt.Errorf("delete out of range tour_stops: %w", err)
	}

	return nil
}

func (r *TourRepository) CreateMissingTourStops(ctx context.Context, tx *sql.Tx, tourID int64, missing []string) error {
	stmt, err := tx.PrepareContext(ctx, createMissingTourStopsQuery)
	if err != nil {
		return fmt.Errorf("prepare tx statement: %w", err)
	}
	defer stmt.Close()

	for _, date := range missing {
		if _, err := stmt.ExecContext(ctx, tourID, "unplanned", date); err != nil {
			return fmt.Errorf("create missing tour_stop: %w", err)
		}
	}

	return nil
}

func (r *TourRepository) Update(ctx context.Context, tx *sql.Tx, tourID int64, u model.TourUpdate) (model.Tour, error) {
	var tour model.Tour
	var placeholders []string
	var args []any

	addColumnValue := func(col string, val any) {
		placeholders = append(placeholders, col+" = ?")
		args = append(args, val)
	}

	if u.Name != nil {
		addColumnValue("name", *u.Name)
	}

	if u.Start != nil {
		addColumnValue("start_date", *u.Start)
		addColumnValue("end_date", *u.End)
	}

	if len(args) == 0 {
		return model.Tour{}, ErrNothingToUpdate
	}

	query := fmt.Sprintf(updateTourQuery, strings.Join(placeholders, ", "))

	if err := tx.QueryRowContext(ctx, query, append(args, tourID)...).Scan(
		&tour.ID, &tour.Name, &tour.Start, &tour.End,
	); err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return model.Tour{}, ErrDuplicate
		}

		return model.Tour{}, fmt.Errorf("update tour: %w", err)
	}

	return tour, nil
}
