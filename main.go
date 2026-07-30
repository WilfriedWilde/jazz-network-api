package main

import (
	"log"

	r "jazz-network-api/repository"
)

func main() {
	//s.UpdateSpreadSheet()

	// init database
	db, err := r.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}
