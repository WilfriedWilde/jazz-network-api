package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

type Organizer struct {
	Type           string
	Name           string
	Country        string
	City           string
	Address        string
	Websites       []string
	Emails         []string
	PhoneNumbers   []string
	ContactPersons []string
	Schedule       string
	Notes          string
	Latitude       float64
	Longitude      float64
}

const organizerTableStatement = `
	CREATE TABLE IF NOT EXISTS organizers (
		id INTEGER PRIMARY KEY,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		country TEXT NOT NULL,
		city TEXT NOT NULL,
		address TEXT NOT NULL,
		websites TEXT,
		emails TEXT,
		phone_numbers TEXT,
		contact_persons TEXT,
		schedule TEXT,
		notes TEXT,
		latitude REAL,
		longitude REAL
);`

func CreateOrganizer(db *sql.DB, o Organizer) error {
	_, err := db.Exec(`
		INSERT INTO organizers (type, name, country, city, address, websites, emails, phone_numbers, contact_persons, schedule, notes, latitude, longitude) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.Type, o.Name, o.Country, o.City, o.Address,
		strings.Join(o.Websites, ";"),
		strings.Join(o.Emails, ";"),
		strings.Join(o.PhoneNumbers, ";"),
		strings.Join(o.ContactPersons, ";"),
		o.Schedule, o.Notes, o.Latitude, o.Longitude,
	)

	if err != nil {
		return fmt.Errorf("error creating organizer: %s", err.Error())
	}

	fmt.Printf("Organizer created successfully:\n %+v\n\n", o)

	return nil
}
