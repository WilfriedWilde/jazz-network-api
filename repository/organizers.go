package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"jazz-network-api/model"
)

type OrganizerRepository struct {
	db *sql.DB
}

func NewOrganizerRepository(db *sql.DB) *OrganizerRepository {
	return &OrganizerRepository{db: db}
}

func (r *OrganizerRepository) CreateOrganizer(o model.Organizer) error {
	_, err := r.db.Exec(`
		INSERT INTO organizers (type, name, country, city, address, websites, emails, phone_numbers, contact_persons, schedule, notes, latitude, longitude) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.Type, o.Name, o.Country, o.City, o.Address,
		o.Websites, o.Emails, o.PhoneNumbers, o.ContactPersons,
		o.Schedule, o.Notes, o.Latitude, o.Longitude,
	)

	if err != nil {
		return fmt.Errorf("error creating organizer: %w", err)
	}

	fmt.Printf("Organizer created successfully:\n %+v\n\n", o)

	return nil
}

func SelectAllOrganizers(ctx context.Context, db *sql.DB) ([]model.Organizer, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT * FROM organizers`,
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting organizer by id: %s", err.Error())
	}
	defer rows.Close()

	organizers := make([]model.Organizer, 0)

	for rows.Next() {
		o := model.Organizer{}

		if err := rows.Scan(
			&o.ID,
			&o.Type,
			&o.Name,
			&o.Country,
			&o.City,
			&o.Address,
			&o.Websites,
			&o.Emails,
			&o.PhoneNumbers,
			&o.ContactPersons,
			&o.Schedule,
			&o.Notes,
			&o.Latitude,
			&o.Longitude,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %s", err.Error())
		}

		organizers = append(organizers, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error decoding rows: %s", err.Error())
	}

	return organizers, nil
}

// TODO: move this somewhere else
func PrepareForJSON(o model.Organizer) OrganizerJSON {
	return OrganizerJSON{
		ID:             o.ID,
		Type:           o.Type,
		Name:           o.Name,
		Country:        o.Country,
		City:           o.City,
		Address:        o.Address,
		Websites:       strings.Split(o.Websites, ";"),
		Emails:         strings.Split(o.Emails, ";"),
		PhoneNumbers:   strings.Split(o.PhoneNumbers, ";"),
		ContactPersons: strings.Split(o.ContactPersons, ";"),
		Schedule:       o.Schedule,
		Notes:          o.Notes,
		Latitude:       o.Latitude,
		Longitude:      o.Longitude,
	}
}
