package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"jazz-network-api/database"
)

type Server struct {
	db *sql.DB
}

var dbPath = "jazznetwork.db"

func main() {
	//s.UpdateSpreadSheet()

	// check if database exists
	exists := database.Exists(dbPath)

	// init database
	db, err := database.OpenDB(fmt.Sprintf("%s?_foreign_keys=on", dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)

	if err := database.CreateTables(db); err != nil {
		log.Fatal(err)
	}

	if !exists {
		if err := database.Seed(db); err != nil {
			log.Fatal(err)
		}
	}

	// init server
	server := Server{
		db: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /organizers", server.HandleSelectAllOrganizers)

	fmt.Println("Running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleWare(mux)))
}
