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

	case errors.Is(err, service.ErrDuplicate):
		log.Print(err)
		writeJSON(w, http.StatusConflict, model.ErrBody{
			Code:    "duplicate",
			Message: "entity already existing",
		})

	case errors.Is(err, service.ErrOrganizerConflict):
		log.Print(err)
		writeJSON(w, http.StatusConflict, model.ErrBody{
			Code:    "conflict",
			Message: "organizer referenced by one or more tour_stops",
		})

	case errors.Is(err, service.ErrTourConflict):
		log.Print(err)
		writeJSON(w, http.StatusConflict, model.ErrBody{
			Code:    "conflict",
			Message: "tour has one or more planned tour_stops",
		})

	case errors.Is(err, service.ErrNothingToUpdate):
		log.Print(err)
		writeJSON(w, http.StatusBadRequest, model.ErrBody{
			Code:    "nothing_to_update",
			Message: "update request body must not be empty",
		})

	case errors.Is(err, service.ErrTourStopConstraintLocation):
		log.Print(err)
		writeJSON(w, http.StatusBadRequest, model.ErrBody{
			Code:    "constraint_location",
			Message: "country and city accepted and required only for type = 'free_day'",
		})

	case errors.Is(err, service.ErrTourStopConstraintOrganizer):
		log.Print(err)
		writeJSON(w, http.StatusBadRequest, model.ErrBody{
			Code:    "constraint_organizer",
			Message: "organizer_id is required only for type = 'gig' and should not be present for the rest of types",
		})

	default:
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, model.ErrBody{
			Code:    "internal_error",
			Message: "an unexpected error occurred",
		})
	}
}
