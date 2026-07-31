package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	r "jazz-network-api/repository"
)

type Server struct {
	db *sql.DB
}

func main() {
	//s.UpdateSpreadSheet()

	// init database
	database, err := r.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// init server
	server := Server{
		db: database,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /organizers", server.HandleSelectAllOrganizers)

	fmt.Println("Running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleWare(mux)))
}
