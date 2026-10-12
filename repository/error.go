package repository

import "errors"

var (
	ErrNotFound                    = errors.New("repository: not found")
	ErrDuplicate                   = errors.New("repository: duplicate value")
	ErrNothingToUpdate             = errors.New("repository: nothing to update")
	ErrTourStopConstraintLocation  = errors.New("repository: country and city required for type in ['gig', 'free_day']")
	ErrTourStopConstraintOrganizer = errors.New("repository: organizer_id is required for type = 'gig'")
)
