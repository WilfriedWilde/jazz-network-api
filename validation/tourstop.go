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

	isFreeDay := update.Type != nil && *update.Type == "free_day"
	if isFreeDay {
		if update.Country == nil {
			addFieldErr("country", "required for type = 'free_day'")
		} else {
			*update.Country = strings.TrimSpace(*update.Country)
			if *update.Country == "" {
				addFieldErr("country", "nonempty")
			}

			if !isValidCountry(*update.Country) {
				addFieldErr("country", "invalid country")
			}
		}

		if update.City == nil {
			addFieldErr("city", "required for type = 'free_day'")
		} else {
			*update.City = strings.TrimSpace(*update.City)
			if *update.City == "" {
				addFieldErr("city", "nonempty")
			}
		}
	} else if update.Country != nil || update.City != nil {
		addFieldErr("country and city", "only allowed together with type = 'free_day'")
	}

	isGig := update.Type != nil && *update.Type == "gig"
	if isGig && update.OrganizerID == nil {
		addFieldErr("organizer_id", "required for type = 'gig'")
	}

	if update.Type != nil {
		*update.Type = strings.TrimSpace(*update.Type)
		if _, found := tourStopTypes[*update.Type]; !found {
			addFieldErr("type", "should be 'gig', 'travel_day', 'free_day' or 'unplanned'")
		}
	}

	if validationErr == nil {
		return nil
	}

	return validationErr
}
