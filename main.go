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

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// init server
	server := Server{
		db: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /organizers", handlers.Organizers.ReadAll)

	fmt.Println("Running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleWare(mux)))
}
