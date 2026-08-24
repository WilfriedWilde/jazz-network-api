package service

import (
	"errors"
	"fmt"
	"jazz-network-api/repository"
)

var (
	ErrNotFound = errors.New("service: not found")
)

func convertToServiceError(operation, resource string, err error) error {
	var serviceErr error

	switch {
	case errors.Is(err, repository.ErrNotFound):
		serviceErr = ErrNotFound
	}

	if serviceErr != nil {
		err = errors.Join(serviceErr, err)
	}

	return fmt.Errorf("%s %s: %w", operation, resource, err)
}
