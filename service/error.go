package service

import (
	"errors"
	"fmt"
	"jazz-network-api/repository"
)

var (
	ErrNotFound           = errors.New("service: not found")
	ErrDuplicate          = errors.New("service: duplicate value")
	ErrUsernameTaken      = errors.New("service: username taken")
	ErrEmailTaken         = errors.New("service: email taken")
	ErrPasswordTooShort   = errors.New("service: password too short")
	ErrPasswordTooWeak    = errors.New("service: password too weak")
	ErrInvalidCredentials = errors.New("service: invalid email or password")
	ErrOrganizerConflict  = errors.New("service: organizer referenced by one or more tour_stops")
)

func convertToServiceError(operation, resource string, err error) error {
	var serviceErr error

	switch {
	case errors.Is(err, repository.ErrNotFound):
		serviceErr = ErrNotFound

	case errors.Is(err, repository.ErrDuplicate):
		serviceErr = ErrDuplicate
	}

	if serviceErr != nil {
		err = errors.Join(serviceErr, err)
	}

	return fmt.Errorf("%s %s: %w", operation, resource, err)
}
