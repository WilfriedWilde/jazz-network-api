package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed init.sql
var initMigration string

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func OpenDB(path string) (*sql.DB, error) {
	// connect to database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("error creating/opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error establishing connection with database: %w", err)
	}

	fmt.Println("Connected to database!")

	return db, nil
}

func CreateTables(db *sql.DB) error {
	if _, err := db.Exec(initMigration); err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}

	return nil
}
