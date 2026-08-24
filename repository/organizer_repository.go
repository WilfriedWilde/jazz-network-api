package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"jazz-network-api/model"
)

const (
	readOrganizerByIDQuery = `
		SELECT *
		FROM organizers
		WHERE id = ?`

	readAllOrganizersQuery = `
		SELECT * 
		FROM organizers
		ORDER BY id ASC`

	createOrganizerQuery = `
		INSERT INTO organizers(type, name, country, city, address, websites, emails, phone_numbers, contact_persons, schedule, notes, latitude, longitude)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING *`
)

type OrganizerRepository struct {
	db *sql.DB
}

func NewOrganizerRepository(db *sql.DB) *OrganizerRepository {
	return &OrganizerRepository{db: db}
}

func (r *OrganizerRepository) Read(ctx context.Context, id int64) (model.Organizer, error) {
	var organizer model.Organizer

	if err := r.db.QueryRowContext(ctx, readOrganizerByIDQuery, id).Scan(
		&organizer.ID, &organizer.Type, &organizer.Name,
		&organizer.Country, &organizer.City, &organizer.Address,
		&organizer.Websites, &organizer.Emails, &organizer.PhoneNumbers,
		&organizer.ContactPersons, &organizer.Schedule, &organizer.Notes,
		&organizer.Latitude, &organizer.Longitude,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Organizer{}, errors.Join(err, ErrNotFound)
		}
		return model.Organizer{}, fmt.Errorf("read organizer by id: %w", err)
	}

	return organizer, nil
}

func (r *OrganizerRepository) ReadAll(ctx context.Context) ([]model.Organizer, error) {
	rows, err := r.db.QueryContext(ctx, readAllOrganizersQuery)
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
			&o.Latitude, &o.Longitude,
		); err != nil {
			return nil, fmt.Errorf("error scanning organizer: %w", err)
		}

		organizers = append(organizers, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organizers rows: %w", err)
	}

	return organizers, nil
}

func (r *OrganizerRepository) Create(ctx context.Context, o model.Organizer) (model.Organizer, error) {
	var organizer model.Organizer

	if err := r.db.QueryRowContext(ctx, createOrganizerQuery,
		o.Type, o.Name, o.Country, o.City, o.Address,
		o.Websites, o.Emails, o.PhoneNumbers, o.ContactPersons,
		o.Schedule, o.Notes, o.Latitude, o.Longitude,
	).Scan(
		&organizer.ID, &organizer.Type, &organizer.Name,
		&organizer.Country, &organizer.City, &organizer.Address,
		&organizer.Websites, &organizer.Emails, &organizer.PhoneNumbers,
		&organizer.ContactPersons, &organizer.Schedule, &organizer.Notes,
		&organizer.Latitude, &organizer.Longitude,
	); err != nil {
		return model.Organizer{}, fmt.Errorf("error creating organizer: %w", err)
	}

	return organizer, nil
}
