package validation

import (
	"jazz-network-api/model"
	"net/http"
	"strings"
	"time"
)

func ValidateTourInput(input model.TourInput) error {
	var validationErr *model.ValidationError

	addFieldErr := func(field, message string) {
		if validationErr == nil {
			validationErr = &model.ValidationError{
				Status: http.StatusBadRequest,
				Code:   "validation_error",
			}
		}

		validationErr.Fields = append(validationErr.Fields, model.FieldError{
			Field:   field,
			Message: message,
		})
	}

	// name required
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		addFieldErr("name", "required")
	}

	// start required
	input.Start = strings.TrimSpace(input.Start)
	if input.Start == "" {
		addFieldErr("start", "required")
	} else if ok, msg := isValidDate(input.Start); !ok {
		addFieldErr("start date", msg)
	}

	// end required
	input.End = strings.TrimSpace(input.End)
	if input.End == "" {
		addFieldErr("end", "required")
	} else if okStart, _ := isValidDate(input.Start); okStart {
		if okEnd, msg := isValidDate(input.End); !okEnd {
			addFieldErr("end date", msg)
		} else {
			// check if end >= start
			if !IsEndAfterStart(input.Start, input.End) {
				addFieldErr("end date", "should be on the same day or after start date")
			}
		}
	}

	if validationErr == nil {
		return nil
	}

	return validationErr
}

func ValidateTourUpdate(update model.TourUpdate) error {
	var validationErr *model.ValidationError

	addFieldErr := func(field, message string) {
		if validationErr == nil {
			validationErr = &model.ValidationError{
				Status: http.StatusBadRequest,
				Code:   "validation_error",
			}
		}

		validationErr.Fields = append(validationErr.Fields, model.FieldError{
			Field:   field,
			Message: message,
		})
	}

	// name nonempty
	if update.Name != nil {
		*update.Name = strings.TrimSpace(*update.Name)
		if *update.Name == "" {
			addFieldErr("name", "nonempty")
		}
	}

	// start valid date
	if update.Start != nil {
		if ok, msg := isValidDate(*update.Start); !ok {
			addFieldErr("start date", msg)
		}

		if update.End == nil {
			addFieldErr("end date", "non empty")
		}
	}

	// end valid date
	if update.End != nil {
		if okEnd, msg := isValidDate(*update.End); !okEnd {
			addFieldErr("end date", msg)
		}

		if update.Start == nil {
			addFieldErr("end date", "non empty")
		}
	}

	if update.Start != nil && update.End != nil {
		// check if end >= start
		if !IsEndAfterStart(*update.Start, *update.End) {
			addFieldErr("end date", "should be on the same day or after start date")
		}
	}

	if validationErr == nil {
		return nil
	}

	return validationErr
}

// isValidDate checks if date has format YYYY-MM-DD, is a valid calendar date and is situated in the future
func isValidDate(date string) (bool, string) {
	d, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return false, err.Error()
	}

	if d.Before(time.Now()) {
		return false, "should be in the future"
	}

	return true, ""
}

func IsEndAfterStart(start, end string) bool {
	startDate, _ := time.Parse(time.DateOnly, start)
	endDate, _ := time.Parse(time.DateOnly, end)

	return endDate.After(startDate)
}
