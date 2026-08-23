package httpapi

import (
	"errors"
	"jazz-network-api/repository"
	"jazz-network-api/service"
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
	var requestError *requestError
	var serviceError *service.ServiceError
	var repositoryError *repository.RepositoryError

	switch {
	case errors.As(err, &requestError):
		writeJSON(w, requestError.GetStatus(), errBody{
			Code:    requestError.GetCode(),
			Message: requestError.Error(),
		})

	case errors.As(err, &serviceError):
		writeJSON(w, serviceError.GetStatus(), errBody{
			Code:    serviceError.GetCode(),
			Message: serviceError.Error(),
		})

	case errors.As(err, &repositoryError):
		writeJSON(w, repositoryError.GetStatus(), errBody{
			Code:    repositoryError.GetCode(),
			Message: repositoryError.Error(),
		})
	}
}
