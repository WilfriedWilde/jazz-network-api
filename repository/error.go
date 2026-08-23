package repository

import "errors"

var (
	ErrNotFound = errors.New("not_found")
)

type RepositoryError struct {
	Status  int
	Code    string
	Message string
}

func (e *RepositoryError) GetStatus() int {
	return e.Status
}

func (e *RepositoryError) GetCode() string {
	return e.Code
}

func (e *RepositoryError) Error() string {
	return e.Message
}
