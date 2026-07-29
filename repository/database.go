package repository

import (
	"database/sql"
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

}
