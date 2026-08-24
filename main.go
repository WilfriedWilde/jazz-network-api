package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"jazz-network-api/database"
	"jazz-network-api/httpapi"
)

type Server struct {
	db *sql.DB
}

var dbPath = "jazznetwork.db"

func main() {
	//spreadsheet.UpdateSpreadSheet()

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := httpapi.NewRouter(db)

	fmt.Println("Running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", httpapi.CorsMiddleWare(router)))
}
