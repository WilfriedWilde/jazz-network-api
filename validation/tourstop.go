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

// TODO: returns wrong error when trying to update type = gig without organizer_id
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
		if _, found := tourStopTypes[*update.Type]; !found {
			addFieldErr("type", "should be 'gig', 'travel_day', 'free_day' or 'unplanned'")
		}

		switch *update.Type {
		case "free_day":
			if update.Country == nil {
				addFieldErr("country", "required for type = 'free_day'")
			}

			if update.City == nil {
				addFieldErr("city", "required for type = 'free_day'")
			}
		case "gig":
			if update.OrganizerID == nil {
				addFieldErr("organizer_id", "required for type = 'gig'")
			}

			if update.Country != nil || update.City != nil {
				addFieldErr("country and city", "not allowed for type != 'free_day'")
			}
		default:
			if update.Country != nil || update.City != nil {
				addFieldErr("country and city", "not allowed for type != 'free_day'")
			}
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
