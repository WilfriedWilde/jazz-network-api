package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed init.sql
var initMigration string

func InitDB(ctx context.Context, path string) (*sql.DB, error) {
	// check if database exists
	exists := exists(path)

	// init database
	db, err := openDB(ctx, fmt.Sprintf("%s?_foreign_keys=on", path))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	if err := createTables(ctx, db); err != nil {
		return nil, err
	}

	if !exists {
		if err := seed(ctx, db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func openDB(ctx context.Context, path string) (*sql.DB, error) {
	// connect to database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("error creating/opening database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error establishing connection with database: %w", err)
	}

	fmt.Println("Connected to database!")

	return db, nil
}

func createTables(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, initMigration); err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}

	return nil
}
