package validation

import (
	"jazz-network-api/model"
	"net/http"
	"strings"
)

var (
	tourStopTypes = map[string]struct{}{
		"gig":        {},
		"travel_day": {},
		"free_day":   {},
		"unplanned":  {},
	}
)

func ValidateTourStopUpdate(update model.TourStopUpdate) error {
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

	// type in ["gig", "travel_day", "free_day", "unplanned"]
	if update.Type != nil {
		*update.Type = strings.TrimSpace(*update.Type)
		if _, found := organizerTypes[*update.Type]; !found {
			addFieldErr("type", "should be 'gig', 'travel_day', 'free_day' or 'unplanned'")
		}
	}

	// country in existing real country
	if update.Country != nil {
		if !isValidCountry(*update.Country) {
			addFieldErr("country", "invalid country")
		}
	}

	// city nonempty
	if update.City != nil {
		*update.City = strings.TrimSpace(*update.City)
		if *update.City == "" {
			addFieldErr("city", "nonempty")
		}
	}

	if validationErr == nil {
		return nil
	}

	return validationErr
}
