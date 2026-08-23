package service

import (
	"errors"
	"jazz-network-api/repository"
	"net/http"
)

type ServiceError struct {
	Status  int
	Code    string
	Message string
}

func (e *ServiceError) GetStatus() int {
	return e.Status
}

func (e *ServiceError) GetCode() string {
	return e.Code
}

func (e *ServiceError) Error() string {
	return e.Message
}

func formatServiceError(err error, entity string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return &ServiceError{
			Status:  http.StatusNotFound,
			Code:    "not_found",
			Message: err.Error(),
		}
	}
}
