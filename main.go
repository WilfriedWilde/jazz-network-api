package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"jazz-network-api/database"
	"jazz-network-api/httpapi"
)

const dbPath = "jazznetwork.db"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	signalCtx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stopSignals()

	//s.UpdateSpreadSheet()

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return httpapi.Server(signalCtx, db)
}
