package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	// check if file exists
	fileExists := true
	_, err := os.Stat("repository/jazznetwork.db")
	if err != nil {
		fileExists = false
	}

	// connect to database
	db, err := sql.Open("sqlite3", "repository/jazznetwork.db?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("error creating/opening jazznetwork.db: %s\n", err.Error())
	}

	fmt.Println("Connected to jazznetwork.db!")

	var sqliteVersion string
	if err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion); err != nil {
		log.Printf("error querying sqlite version: %s\n", err.Error())
	}

	fmt.Println(sqliteVersion)

	if err := CreateTables(db); err != nil {
		return nil, err
	}

	// add all organizers from spreadsheet
	if !fileExists {
		organizers, err := GetOrganizers()
		if err != nil {
			return nil, err
		}

		for _, o := range organizers {
			if err := CreateOrganizer(db, o); err != nil {
				return nil, err
			}
		}
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	var errs []error

	tablesStatements := map[string]string{
		"organizers":         organizerTableStatement,
		"organizer contacts": contactsTableStatement,
		"organizer notes":    notesTableStatement,
		"users":              usersTableStatement,
		"tours":              toursTableStatement,
		"tour stops":         stopsTableStatement,
	}

	for name, statement := range tablesStatements {
		if _, err := db.Exec(statement); err != nil {
			errs = append(errs, fmt.Errorf("error creating table %s: %s", name, err.Error()))
		}
	}

	return errors.Join(errs...)
}
