package repository

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetOrganizersFromCSV() ([]Organizer, error) {
	file, err := os.Open("spreadsheet/jazz_network.csv")
	if err != nil {
		return nil, fmt.Errorf("error opening jazz_network.csv: %s\n", err.Error())
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading jazz_network.csv: %s\n", err.Error())
	}

	organizers := make([]Organizer, len(records)-1)

	for rowID, row := range records[1:] {
		latitude, err := strconv.ParseFloat(row[11], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing latitude: %s\n", err.Error())
		}
		longitude, err := strconv.ParseFloat(row[12], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing longitude: %s\n", err.Error())
		}

		organizers[rowID] = Organizer{
			Type:           strings.TrimSpace(strings.ToLower(row[0])),
			Name:           strings.TrimSpace(strings.ToLower(row[1])),
			Country:        strings.TrimSpace(strings.ToLower(row[2])),
			City:           strings.TrimSpace(strings.ToLower(row[3])),
			Address:        strings.TrimSpace(strings.ToLower(row[4])),
			Websites:       strings.TrimSpace(strings.ToLower(row[5])),
			Emails:         strings.TrimSpace(strings.ToLower(row[6])),
			PhoneNumbers:   strings.TrimSpace(strings.ToLower(row[7])),
			ContactPersons: strings.TrimSpace(strings.ToLower(row[8])),
			Schedule:       strings.TrimSpace(strings.ToLower(row[9])),
			Notes:          strings.TrimSpace(row[10]),
			Latitude:       latitude,
			Longitude:      longitude,
		}
	}

	return organizers, nil
}

func CreateOrganizer(db *sql.DB, o Organizer) error {
	_, err := db.Exec(`
		INSERT INTO organizers (type, name, country, city, address, websites, emails, phone_numbers, contact_persons, schedule, notes, latitude, longitude) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.Type, o.Name, o.Country, o.City, o.Address,
		o.Websites, o.Emails, o.PhoneNumbers, o.ContactPersons,
		o.Schedule, o.Notes, o.Latitude, o.Longitude,
	)

	if err != nil {
		return fmt.Errorf("error creating organizer: %s\n", err.Error())
	}

	fmt.Printf("Organizer created successfully:\n %+v\n\n", o)

	return nil
}

func SelectAllOrganizers(ctx context.Context, db *sql.DB) ([]Organizer, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT * FROM organizers`,
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting organizer by id: %s", err.Error())
	}
	defer rows.Close()

	organizers := make([]Organizer, 0)

	for rows.Next() {
		o := Organizer{}

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

func (o *Organizer) PrepareForJSON() OrganizerJSON {
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
