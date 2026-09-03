package main

import (
	"database/sql"
	"jazz-network-api/httpapi"
	"jazz-network-api/repository"
	"jazz-network-api/service"
	"net/http"
)

type handlers struct {
	organizer *httpapi.OrganizerHandler
}

func newRouter(db *sql.DB) http.Handler {
	handlers := getHandlers(db)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/organizers/{id}", handlers.organizer.Read)
	mux.HandleFunc("GET /api/organizers", handlers.organizer.ReadAll)
	mux.HandleFunc("POST /api/organizers", handlers.organizer.Create)

	return mux
}

func getHandlers(db *sql.DB) handlers {
	organizerRepository := repository.NewOrganizerRepository(db)
	organizerService := service.NewOrganizerService(organizerRepository)

	sessionRepository := repository.NewSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepository)

	return handlers{
		organizer: &httpapi.OrganizerHandler{Service: organizerService},
	}
}
