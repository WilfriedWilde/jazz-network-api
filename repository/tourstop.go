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
	existTourStopQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM tours
			WHERE id = ?
		)
	`

	updateTourStopQuery = `
		UPDATE tour_stops
		SET %s
		WHERE id = ?
		RETURNING *
	`
)

type TourStopRepository struct {
	DB *sql.DB
}

func NewTourStopRepository(db *sql.DB) *TourStopRepository {
	return &TourStopRepository{DB: db}
}

func (r *TourStopRepository) Exists(ctx context.Context, db DBTX, id int64) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, existTourStopQuery, id).Scan(&exists)
	return exists, err
}

func (r *TourStopRepository) ExistsOrganizer(ctx context.Context, db DBTX, organizerID int64) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, existOrganizerQuery, organizerID).Scan(&exists)
	return exists, err
}

func (r *TourStopRepository) Update(ctx context.Context, tx *sql.Tx, tourStopID int64, u model.TourStopUpdate) (model.TourStop, error) {
	var ts model.TourStop
	var placeholders []string
	var args []any

	addColumnValue := func(col string, val any) {
		placeholders = append(placeholders, col+" = ?")
		args = append(args, val)
	}

	if u.OrganizerID != nil {
		addColumnValue("organizer_id", *u.OrganizerID)
	}

	if u.Type != nil {
		addColumnValue("type", *u.Type)
	}

	if u.Country != nil {
		addColumnValue("country", *u.Country)
	}

	if u.City != nil {
		addColumnValue("city", *u.City)
	}

	if len(args) == 0 {
		return model.TourStop{}, ErrNothingToUpdate
	}

	query := fmt.Sprintf(updateTourStopQuery, strings.Join(placeholders, ", "))

	if err := tx.QueryRowContext(ctx, query, append(args, tourStopID)...).Scan(
		&ts.ID, &ts.OrganizerID, &ts.Type, &ts.Date, &ts.Country, &ts.City, &ts.Notes,
	); err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintCheck {
			switch {
			case strings.Contains(err.Error(), "check_location") == true:
				return model.TourStop{}, ErrTourStopConstraintLocation

			case strings.Contains(err.Error(), "check_organizer") == true:
				return model.TourStop{}, ErrTourStopConstraintOrganizer
			}
		}

		return model.TourStop{}, fmt.Errorf("update tour_stop: %w", err)
	}

	return ts, nil
}
