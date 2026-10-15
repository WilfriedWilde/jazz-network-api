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
	session   *httpapi.SessionHandler
	organizer *httpapi.OrganizerHandler
	tour      *httpapi.TourHandler
	tourStop  *httpapi.TourStopHandler
}

func newRouter(db *sql.DB) http.Handler {
	handlers := getHandlers(db)
	sessionRepository := repository.NewSessionRepository(db)
	requireAuth := m.AuthMiddleware(sessionRepository)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/signup", handlers.session.Signup)
	mux.HandleFunc("POST /api/login", handlers.session.Login)
	mux.HandleFunc("POST /api/logout", handlers.session.Logout)
	mux.Handle("GET /api/me", requireAuth(http.HandlerFunc(handlers.session.Me)))

	mux.Handle("GET /api/organizers", requireAuth(http.HandlerFunc(handlers.organizer.ReadAll)))
	mux.Handle("GET /api/organizers/{id}", requireAuth(http.HandlerFunc(handlers.organizer.Read)))
	mux.Handle("POST /api/organizers", requireAuth(http.HandlerFunc(handlers.organizer.Create)))
	mux.Handle("DELETE /api/organizers/{id}", requireAuth(http.HandlerFunc(handlers.organizer.Delete)))
	mux.Handle("PATCH /api/organizers/{id}", requireAuth(http.HandlerFunc(handlers.organizer.Update)))

	mux.Handle("POST /api/tours", requireAuth(http.HandlerFunc(handlers.tour.Create)))
	mux.Handle("GET /api/tours", requireAuth(http.HandlerFunc(handlers.tour.ReadAll)))
	mux.Handle("GET /api/tours/{id}", requireAuth(http.HandlerFunc(handlers.tour.Read)))
	mux.Handle("PATCH /api/tours/{id}", requireAuth(http.HandlerFunc(handlers.tour.Update)))
	mux.Handle("DELETE /api/tours/{id}", requireAuth(http.HandlerFunc(handlers.tour.Delete)))

	mux.Handle("PATCH /api/tourstops/{id}", requireAuth(http.HandlerFunc(handlers.tourStop.Update)))

	return mux
}

func getHandlers(db *sql.DB) handlers {
	sessionRepository := repository.NewSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepository)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	organizerRepository := repository.NewOrganizerRepository(db)
	organizerService := service.NewOrganizerService(organizerRepository)

	tourRepository := repository.NewTourRepository(db)
	tourService := service.NewTourService(tourRepository)

	tourStopRepository := repository.NewTourStopRepository(db)
	tourStopService := service.NewTourStopService(tourStopRepository)

	return handlers{
		session: &httpapi.SessionHandler{
			SessionService: sessionService,
			UserService:    userService,
		},
		organizer: &httpapi.OrganizerHandler{Service: organizerService},
		tour:      &httpapi.TourHandler{Service: tourService},
		tourStop:  &httpapi.TourStopHandler{Service: tourStopService},
	}
}
