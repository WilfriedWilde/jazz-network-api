package httpapi

import (
	"errors"
	"jazz-network-api/service"
	"log"
	"net/http"
)

type requestError struct {
	Status  int
	Code    string
	Message string
}

func (e *requestError) GetStatus() int {
	return e.Status
}

func (e *requestError) GetCode() string {
	return e.Code
}

func (e *requestError) Error() string {
	return e.Message
}

func writeError(w http.ResponseWriter, err error) {
	var requestErr *requestError

	if errors.As(err, &requestErr) {
		writeJSON(w, requestErr.GetStatus(), errBody{
			Code:    requestErr.GetCode(),
			Message: requestErr.Error(),
		})
	}

	writeAPIerror(w, err)
}

func writeAPIerror(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		log.Print(err)
		writeJSON(w, http.StatusNotFound, errBody{
			Code:    "not_found",
			Message: "resource not found",
		})
	default:
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, errBody{
			Code:    "internal_error",
			Message: "an unexpected error occurred",
		})
	}
}
