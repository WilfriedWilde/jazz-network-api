package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

	readTourStopIDsByOrganizerID = `
		SELECT id
		FROM tour_stops
		WHERE organizer_id = ?
	`

	resetOrganizerTourStopsQuery = `
		UPDATE tour_stops
		SET organizer_id = NULL, type = 'free_day'
		WHERE id IN (%s)
	`

	deleteOrganizerTourStopFinancesQuery = `
		DELETE FROM tour_stop_finances
		WHERE tour_stop_id IN (%s)
	`

	deleteOrganizerTravelsQuery = `
		DELETE FROM travels
		WHERE from_stop_id IN (%s) OR to_stop_id IN (%s)
	`

	deleteOrganizerAccommodationsQuery = `
		DELETE FROM accommodations
		WHERE tour_stop_id IN (%s)
	`

	existOrganizerQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM organizers
			WHERE id = ?
		)
	`

	updateOrganizerQuery = `
		UPDATE organizers
		SET %s
		WHERE id = ?
		RETURNING id, type, name, country, city, address, websites, emails, phone_numbers, contact_persons, schedule, notes, latitude, longitude, public_visibility
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

func (r *OrganizerRepository) Exists(ctx context.Context, db DBTX, id int64) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, existOrganizerQuery, id).Scan(&exists)
	return exists, err
}

func (r *OrganizerRepository) Update(ctx context.Context, db DBTX, id int64, o model.OrganizerUpdate) (model.Organizer, error) {
	var organizer model.Organizer
	var placeholders []string
	var args []any

	addColumnValue := func(col string, val any) {
		placeholders = append(placeholders, col+" = ?")
		args = append(args, val)
	}

	if o.Type != nil {
		addColumnValue("type", *o.Type)
	}

	if o.Name != nil {
		addColumnValue("name", *o.Name)
	}

	if o.Country != nil {
		addColumnValue("country", *o.Country)
	}

	if o.City != nil {
		addColumnValue("city", *o.City)
	}

	if o.Address != nil {
		addColumnValue("address", *o.Address)
	}

	if o.Websites != nil {
		addColumnValue("websites", strings.Join(o.Websites, ";"))
	}

	if o.Emails != nil {
		addColumnValue("emails", strings.Join(o.Emails, ";"))
	}

	if o.PhoneNumbers != nil {
		addColumnValue("phone_numbers", strings.Join(o.PhoneNumbers, ";"))
	}

	if o.ContactPersons != nil {
		addColumnValue("contact_persons", strings.Join(o.ContactPersons, ";"))
	}

	if o.Schedule != nil {
		addColumnValue("schedule", o.Schedule)
	}

	if o.Notes != nil {
		addColumnValue("notes", o.Notes)
	}

	if o.Latitude != nil {
		addColumnValue("notes", o.Latitude)
	}

	if o.Longitude != nil {
		addColumnValue("longitude", o.Longitude)
	}

	if o.PublicVisibility != nil {
		addColumnValue("public_visibility", o.PublicVisibility)
	}

	if len(args) == 0 {
		return model.Organizer{}, ErrNothingToUpdate
	}

	query := fmt.Sprintf(updateOrganizerQuery, strings.Join(placeholders, ", "))

	if err := db.QueryRowContext(ctx, query, append(args, id)...).Scan(
		&organizer.ID, &organizer.Type, &organizer.Name,
		&organizer.Country, &organizer.City, &organizer.Address,
		&organizer.Websites, &organizer.Emails, &organizer.PhoneNumbers,
		&organizer.ContactPersons, &organizer.Schedule, &organizer.Notes,
		&organizer.Latitude, &organizer.Longitude, &organizer.PublicVisibility,
	); err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			return model.Organizer{}, ErrDuplicate
		}

		return model.Organizer{}, fmt.Errorf("update organizer: %w", err)
	}

	return organizer, nil
}

func (r *OrganizerRepository) Delete(ctx context.Context, db DBTX, organizerID int64) error {
	result, err := db.ExecContext(ctx, deleteOrganizerQuery, organizerID)
	if err != nil {
		return fmt.Errorf("delete organizer: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete organizer (rows affected): %w", err)
	}

	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *OrganizerRepository) ReadTourStopIDsByOrganizerID(ctx context.Context, db DBTX, organizerID int64) ([]int64, error) {
	var ids []int64

	rows, err := db.QueryContext(ctx, readTourStopIDsByOrganizerID, organizerID)
	if err != nil {
		return nil, fmt.Errorf("read tour_stop IDs by organizer ID: %w", err)
	}

	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan tour_stop ID: %w", err)
		}

		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tour_stops: %w", err)
	}

	return ids, nil
}

func (r *OrganizerRepository) ResetOrganizerTourStops(ctx context.Context, db DBTX, tourStopIDs []int64) error {
	query, args := buildIDQueryArguments(resetOrganizerTourStopsQuery, tourStopIDs)

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("reset organizer tour_stops: %w", err)
	}

	return nil
}

func (r *OrganizerRepository) DeleteOrganizerTourStopFinances(ctx context.Context, db DBTX, tourStopIDs []int64) error {
	query, args := buildIDQueryArguments(deleteOrganizerTourStopFinancesQuery, tourStopIDs)

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("delete tour_stop_finances (organizer): %w", err)
	}

	return nil
}

func (r *OrganizerRepository) DeleteOrganizerTravels(ctx context.Context, db DBTX, tourStopIDs []int64) error {
	placeholders := make([]string, len(tourStopIDs))
	args := make([]any, len(tourStopIDs))

	for i, id := range tourStopIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(deleteOrganizerTravelsQuery, strings.Join(placeholders, ", "), strings.Join(placeholders, ", "))

	if _, err := db.ExecContext(ctx, query, append(args, args...)...); err != nil {
		return fmt.Errorf("delete travels (organizer): %w", err)
	}

	return nil
}

func (r *OrganizerRepository) DeleteOrganizerAccommodations(ctx context.Context, db DBTX, tourStopIDs []int64) error {
	query, args := buildIDQueryArguments(deleteOrganizerAccommodationsQuery, tourStopIDs)

	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("delete accommodations (organizer): %w", err)
	}

	return nil
}
