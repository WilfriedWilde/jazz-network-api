package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"log"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./jazznetwork.db")
	if err != nil {
		return nil, fmt.Errorf("error creating/opening jazznetwork.db: %s", err.Error())
	}

	fmt.Println("Connected to jazznetwork.db!")

	var sqliteVersion string
	if err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion); err != nil {
		log.Printf("error querying sqlite version: %s", err.Error())
	}

	fmt.Println(sqliteVersion)

	return db, nil
}

func CreateTables(db *sql.DB) error {
	var errs []error

	tablesStatements := map[string]string{
		"organizers": `
			CREATE TABLE IF NOT EXISTS organizers (
				id INTEGER PRIMARY KEY,
				userID INTEGER,
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
				longitude REAL,
				status TEXT DEFAULT 'not contacted'
			);`,
		"statuses": `
			CREATE TABLE IF NOT EXISTS statuses (
				id INTEGER PRIMARY KEY,
				userID INTEGER,
				organizerID INTEGER,
				tourID INTEGER,
				contacted_in_person BOOLEAN,
				contacted_phone BOOLEAN,
				contacted_email BOOLEAN,
				timestamp DATETIME
			);`,
		"users": `
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY,
				tourID INTEGER,
				name TEXT NOT NULL
			);`,
		"tours": `
			CREATE TABLE IF NOT EXISTS tours (
				id INTEGER PRIMARY KEY,
				userID INTEGER,
				name TEXT NOT NULL
				FOREIGN KEY (userdID) REFERENCES users(id)
		);`,
	}

	for name, statement := range tablesStatements {
		if _, err := db.Exec(statement); err != nil {
			errs = append(errs, fmt.Errorf("error creating table %s: %s", name, err.Error()))
		}
	}

	return errors.Join(errs...)
}
