package repository

import "errors"

var (
	ErrNotFound        = errors.New("repository: not found")
	ErrDuplicate       = errors.New("repository: duplicate value")
	ErrNothingToUpdate = errors.New("repository: nothing to update")
)
