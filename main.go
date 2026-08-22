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

func main() {
	//s.UpdateSpreadSheet()

	// init database
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)

	// init server
	server := Server{
		db: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /organizers", server.HandleSelectAllOrganizers)

	fmt.Println("Running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleWare(mux)))
}
