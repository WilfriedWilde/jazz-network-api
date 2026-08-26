package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jazz-network-api/database"
)

type Server struct {
	db *sql.DB
}

const (
	dbPath            = "jazznetwork.db"
	serverAddress     = ":8080"
	readHeaderTimeout = time.Second * 5
	shutDownTimeout   = time.Second * 10
)

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

	handler := corsMiddleWare(http.NewServeMux())

	server := &http.Server{
		Addr:              serverAddress,
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           handler,
	}

	fmt.Println("Running on http://localhost:8080")

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
	case <-signalCtx.Done():
		log.Println("shutdown request")
	}

	stopSignals()

	shutdDownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutDownTimeout)
	defer cancelShutdown()

	if err := server.Shutdown(shutdDownCtx); err != nil {
		shutdownErr := fmt.Errorf("gracefully shut down HTTP server: %w", err)

		if closeErr := server.Close(); closeErr != nil {
			return errors.Join(shutdownErr, fmt.Errorf("force-close HTTP server: %w", closeErr))
		}

		return shutdownErr
	}

	err = <-serverErr

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server stopped unexpectedly: %w", err)
	}

	log.Printf("HTTP server stopped")

	return nil
}
