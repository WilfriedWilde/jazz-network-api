package httpapi

import (
	"database/sql"
	"jazz-network-api/repository"
	"jazz-network-api/service"
	"net/http"
)

type handlers struct {
	organizer *OrganizerHandler
}

func newRouter(db *sql.DB) http.Handler {
	handlers := getHandlers(db)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/organizers/{id}", handlers.organizer.Read)
	mux.HandleFunc("GET /api/organizers", handlers.organizer.ReadAll)

	return mux
}

func getHandlers(db *sql.DB) handlers {
	organizerRepository := repository.NewOrganizerRepository(db)
	organizerService := service.NewOrganizerService(organizerRepository)

	return handlers{
		organizer: &OrganizerHandler{service: organizerService},
	}
}
