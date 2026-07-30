package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./jazznetwork.db?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("error creating/opening jazznetwork.db: %s", err.Error())
	}

	fmt.Println("Connected to jazznetwork.db!")

	var sqliteVersion string
	if err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion); err != nil {
		log.Printf("error querying sqlite version: %s", err.Error())
	}

	fmt.Println(sqliteVersion)

	if err := CreateTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	var errs []error

	tablesStatements := map[string]string{
		"organizers": `
			CREATE TABLE IF NOT EXISTS organizers (
				id INTEGER PRIMARY KEY,
				user_id INTEGER REFERENCES users(id),
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
			);`,
		"organizer contacts": `
			CREATE TABLE IF NOT EXISTS organizer_contacts (
				id INTEGER PRIMARY KEY,
				user_id INTEGER REFERENCES users(id),
				organizer_id INTEGER REFERENCES organizers(id),
				method TEXT NOT NULL CHECK (method IN ('email', 'phone', 'in_person')),
				status TEXT CHECK (status IN ('not_tried', 'contacted', 'confirmed', 'declined')),
				last_attempted_at TEXT,
				notes TEXT
			);`,
		"organizer notes": `
			CREATE TABLE IF NOT EXISTS organizer_notes (
				id INTEGER PRIMARY KEY,
				user_id INTEGER REFERENCES users(id),
				organizer_id INTEGER REFERENCES organizers(id),
				body TEXT NOT NULL,
				created_at TEXT
			);`,
		"users": `
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY,
				full_name TEXT NOT NULL
			);`,
		"tours": `
			CREATE TABLE IF NOT EXISTS tours (
				id INTEGER PRIMARY KEY,
				user_id INTEGER REFERENCES users(id),
				name TEXT NOT NULL
		);`,
	}

	for name, statement := range tablesStatements {
		if _, err := db.Exec(statement); err != nil {
			errs = append(errs, fmt.Errorf("error creating table %s: %s", name, err.Error()))
		}
	}

	return errors.Join(errs...)
}
