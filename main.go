package main

import (
	"log"

	r "jazz-network-api/repository"
)

func main() {
	//s.UpdateSpreadSheet()

	// create/open database
	db, err := r.InitDB()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	if err = r.CreateTables(db); err != nil {

	}
}
