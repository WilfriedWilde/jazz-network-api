package service

import (
	"errors"
	"fmt"
	"jazz-network-api/model"
	"jazz-network-api/repository"
	"net/http"
)

var (
	ErrNotFound           = errors.New("service: not found")
	ErrDuplicate          = errors.New("service: duplicate value")
	ErrUsernameTaken      = errors.New("service: username taken")
	ErrEmailTaken         = errors.New("service: email taken")
	ErrInvalidCredentials = errors.New("service: invalid email or password")
	ErrOrganizerConflict  = errors.New("service: organizer referenced by one or more tour_stops")
	ErrNothingToUpdate    = errors.New("service: nothing to update")
	ErrDateOrder          = &model.ValidationError{
		Status: http.StatusBadRequest,
		Code:   "validation_error",
		Fields: []model.FieldError{
			model.FieldError{
				Field:   "end date",
				Message: "should be on the same day or after start date",
			},
		},
	}
)

func convertToServiceError(operation, resource string, err error) error {
	var serviceErr error

	switch {
	case errors.Is(err, repository.ErrNotFound):
		serviceErr = ErrNotFound

	case errors.Is(err, repository.ErrDuplicate):
		serviceErr = ErrDuplicate

	case errors.Is(err, repository.ErrNothingToUpdate):
		serviceErr = ErrNothingToUpdate
	}

	if serviceErr != nil {
		err = errors.Join(serviceErr, err)
	}

	return fmt.Errorf("%s %s: %w", operation, resource, err)
}
