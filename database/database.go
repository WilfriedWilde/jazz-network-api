package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (*sql.DB, error) {
	// connect to database
	db, err := sql.Open("sqlite3", "repository/jazznetwork.db?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("error creating/opening jazznetwork.db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error establishing connection with database: %w", err)
	}

	fmt.Println("Connected to jazznetwork.db!")

	return db, nil
}

func CreateTables(db *sql.DB) error {
	return nil
}
