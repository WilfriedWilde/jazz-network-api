package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	serverAddress     = ":8080"
	readHeaderTimeout = time.Second * 5
	shutDownTimeout   = time.Second * 10
)

func server(ctx context.Context, db *sql.DB) error {
	router := newRouter(db)
	handler := corsMiddleWare(router)

	server := &http.Server{
		Addr:              serverAddress,
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           handler,
	}

	return serve(ctx, server)
}

func serve(ctx context.Context, server *http.Server) error {
	fmt.Printf("Server listening and serving HTTP on http://localhost%s\n", serverAddress)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
	case <-ctx.Done():
		log.Println("shutdown request")
	}

	shutdDownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutDownTimeout)
	defer cancelShutdown()

	if err := server.Shutdown(shutdDownCtx); err != nil {
		shutdownErr := fmt.Errorf("gracefully shut down HTTP server: %w", err)

		if closeErr := server.Close(); closeErr != nil {
			return errors.Join(shutdownErr, fmt.Errorf("force-close HTTP server: %w", closeErr))
		}

		return shutdownErr
	}

	err := <-serverErr

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server stopped unexpectedly: %w", err)
	}

	log.Printf("HTTP server stopped")

	return nil
}
