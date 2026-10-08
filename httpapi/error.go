package httpapi

import (
	"errors"
	"jazz-network-api/model"
	"jazz-network-api/service"
	"log"
	"net/http"
)

func writeError(w http.ResponseWriter, err error) {
	var requestErr *model.RequestError
	var validationErr *model.ValidationError

	switch {
	case errors.As(err, &requestErr):
		writeJSON(w, requestErr.Status, model.ErrBody{
			Code:    requestErr.Code,
			Message: requestErr.Error(),
		})

	case errors.As(err, &validationErr):
		writeJSON(w, validationErr.Status, model.ErrFields{
			Code:   validationErr.Code,
			Errors: validationErr.Fields,
		})

	default:
		writeAPIerror(w, err)
	}
}

func writeAPIerror(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		log.Print(err)
		writeJSON(w, http.StatusNotFound, model.ErrBody{
			Code:    "not_found",
			Message: "resource not found",
		})

	case errors.Is(err, service.ErrOrganizerConflict):
		log.Print(err)
		writeJSON(w, http.StatusConflict, model.ErrBody{
			Code:    "conflict",
			Message: "organizer referenced by one or more tour_stops",
		})

	default:
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, model.ErrBody{
			Code:    "internal_error",
			Message: "an unexpected error occurred",
		})
	}
}
