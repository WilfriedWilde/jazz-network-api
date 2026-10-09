package main

import (
	"database/sql"
	"jazz-network-api/httpapi"
	m "jazz-network-api/middleware"
	"jazz-network-api/repository"
	"jazz-network-api/service"
	"net/http"
)

type handlers struct {
	organizer *httpapi.OrganizerHandler
	session   *httpapi.SessionHandler
}

func newRouter(db *sql.DB) http.Handler {
	handlers := getHandlers(db)
	sessionRepository := repository.NewSessionRepository(db)
	requireAuth := m.AuthMiddleware(sessionRepository)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/signup", handlers.session.Signup)
	mux.HandleFunc("POST /api/login", handlers.session.Login)
	mux.HandleFunc("POST /api/logout", handlers.session.Logout)

	mux.Handle("GET /api/organizers", requireAuth(http.HandlerFunc(handlers.organizer.ReadAll)))
	mux.Handle("GET /api/organizers/{id}", requireAuth(http.HandlerFunc(handlers.organizer.Read)))
	mux.Handle("POST /api/organizers", requireAuth(http.HandlerFunc(handlers.organizer.Create)))
	mux.Handle("DELETE /api/organizers/{id}", requireAuth(http.HandlerFunc(handlers.organizer.Delete)))
	mux.Handle("PATCH /api/organizers/{id}", requireAuth(http.HandlerFunc(handlers.organizer.Update)))

	return mux
}

func getHandlers(db *sql.DB) handlers {
	organizerRepository := repository.NewOrganizerRepository(db)
	organizerService := service.NewOrganizerService(organizerRepository)

	sessionRepository := repository.NewSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepository)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	return handlers{
		organizer: &httpapi.OrganizerHandler{Service: organizerService},
		session: &httpapi.SessionHandler{
			SessionService: sessionService,
			UserService:    userService,
		},
	}
}
