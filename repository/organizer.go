package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"jazz-network-api/model"

	"github.com/mattn/go-sqlite3"
)

const (
	readOrganizerByIDQuery = `
		SELECT id, type, name, country, city, address, websites, emails, phone_numbers, contact_persons, schedule, notes, latitude, longitude, public_visibility
		FROM organizers
		WHERE id = ?
		AND (
			public_visibility = TRUE
			OR user_id = ?
			OR (user_id IS NULL AND ?)
		)
	`

	readAllOrganizersQuery = `
		SELECT id, type, name, country, city, address, websites, emails, phone_numbers, contact_persons, schedule, notes, latitude, longitude, public_visibility 
		FROM organizers
		WHERE public_visibility = TRUE
			OR user_id = ?
			OR (user_id IS NULL AND ?) 
		ORDER BY id ASC
	`

	createOrganizerQuery = `
		INSERT INTO organizers(user_id, type, name, country, city, address, websites, emails, phone_numbers, contact_persons, schedule, notes, latitude, longitude, public_visibility)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, type, name, country, city, address, websites, emails, phone_numbers, contact_persons, schedule, notes, latitude, longitude, public_visibility
	`

	deleteOrganizerQuery = `
		DELETE FROM organizers
		WHERE id = ?
	`

	organizerHasTourStopsQuery = `
		SELECT EXISTS(
			SELECT 1
			FROM tour_stops
			WHERE organizer_id = ?
		)
	`
)

type OrganizerRepository struct {
	DB *sql.DB
}

func NewOrganizerRepository(db *sql.DB) *OrganizerRepository {
	return &OrganizerRepository{DB: db}
}

func (r *OrganizerRepository) Read(ctx context.Context, db DBTX, organizerID, userID int64, isAdmin bool) (model.Organizer, error) {
	var organizer model.Organizer

	if err := db.QueryRowContext(ctx, readOrganizerByIDQuery, organizerID, userID, isAdmin).Scan(
		&organizer.ID, &organizer.Type, &organizer.Name,
		&organizer.Country, &organizer.City, &organizer.Address,
		&organizer.Websites, &organizer.Emails, &organizer.PhoneNumbers,
		&organizer.ContactPersons, &organizer.Schedule, &organizer.Notes,
		&organizer.Latitude, &organizer.Longitude, &organizer.PublicVisibility,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Organizer{}, ErrNotFound
		}

		return model.Organizer{}, fmt.Errorf("read organizer by id: %w", err)
	}

	return organizer, nil
}

func (r *OrganizerRepository) ReadAll(ctx context.Context, db DBTX, userID int64, isAdmin bool) ([]model.Organizer, error) {
	rows, err := db.QueryContext(ctx, readAllOrganizersQuery, userID, isAdmin)
	if err != nil {
		return nil, fmt.Errorf("read all organizers: %w", err)
	}
	defer rows.Close()

	organizers := make([]model.Organizer, 0)

	for rows.Next() {
		o := model.Organizer{}

		if err := rows.Scan(
			&o.ID, &o.Type, &o.Name,
			&o.Country, &o.City, &o.Address,
			&o.Websites, &o.Emails, &o.PhoneNumbers,
			&o.ContactPersons, &o.Schedule, &o.Notes,
			&o.Latitude, &o.Longitude, &o.PublicVisibility,
		); err != nil {
			return nil, fmt.Errorf("scan organizer: %w", err)
		}

		organizers = append(organizers, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate over organizers rows: %w", err)
	}

	return organizers, nil
}

func (r *OrganizerRepository) Create(ctx context.Context, db DBTX, o model.OrganizerCreate) (model.Organizer, error) {
	var organizer model.Organizer

	if err := db.QueryRowContext(ctx, createOrganizerQuery,
		o.UserID, o.Type, o.Name, o.Country, o.City, o.Address,
		o.Websites, o.Emails, o.PhoneNumbers, o.ContactPersons,
		o.Schedule, o.Notes, o.Latitude, o.Longitude,
		o.PublicVisibility,
	).Scan(
		&organizer.ID, &organizer.Type, &organizer.Name,
		&organizer.Country, &organizer.City, &organizer.Address,
		&organizer.Websites, &organizer.Emails, &organizer.PhoneNumbers,
		&organizer.ContactPersons, &organizer.Schedule, &organizer.Notes,
		&organizer.Latitude, &organizer.Longitude, &organizer.PublicVisibility,
	); err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			return model.Organizer{}, ErrDuplicate
		}

		return model.Organizer{}, fmt.Errorf("create organizer: %w", err)
	}

	return organizer, nil
}

func (r *OrganizerRepository) Delete(ctx context.Context, db DBTX, organizerID int64) error {
	if _, err := db.ExecContext(ctx, deleteOrganizerQuery, organizerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}

		return fmt.Errorf("delete organizer: %w", err)
	}

	return nil
}

func (r *OrganizerRepository) OrganizerHasTourStops(ctx context.Context, db DBTX, organizerID int64) (bool, error) {
	var exists bool
	if err := db.QueryRowContext(ctx, organizerHasTourStopsQuery, organizerID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check organizer has tour stops", err)
	}

	return exists, nil
}
