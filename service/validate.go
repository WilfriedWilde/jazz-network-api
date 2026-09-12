package service

import (
	"jazz-network-api/model"
	"jazz-network-api/validation"
	"net/http"
	"strings"
	"unicode"
)

var (
	organizerTypes = map[string]struct{}{
		"venue":    {},
		"festival": {},
	}
)

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var hasLetter, hasNumber bool
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsNumber(r):
			hasNumber = true
		}
	}

	if !hasLetter || !hasNumber {
		return ErrPasswordTooWeak
	}

	return nil
}

func validateOrganizerInput(input model.OrganizerInput) error {
	var validationErr *model.ValidationError

	addFieldErr := func(field, message string) {
		if validationErr == nil {
			validationErr = &model.ValidationError{
				Status: http.StatusBadRequest,
				Code:   "vallidation_error",
			}
		}

		validationErr.Fields = append(validationErr.Fields, model.FieldError{
			Field:   field,
			Message: message,
		})
	}

	// type in ["venue", "festival"], required
	input.Type = strings.TrimSpace(input.Type)
	if input.Type == "" {
		addFieldErr("type", "required")
	}

	if _, found := organizerTypes[input.Type]; !found {
		addFieldErr("type", "should be either 'venue' or 'festival'")
	}

	// name required
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		addFieldErr("name", "required")
	}

	// country in existing real country, required
	input.Country = strings.TrimSpace(input.Country)

	// city required
	input.City = strings.TrimSpace(input.City)
	if input.City == "" {
		addFieldErr("city", "required")
	}

	// address required
	input.Address = strings.TrimSpace(input.Address)
	if input.Address == "" {
		addFieldErr("address", "required")
	}

	// websites valid format?
	for _, w := range input.Websites {
		if !validation.IsValidWebsite(w) {
			addFieldErr("websites", "wrong format")
		}
	}

	// emails valid format
	for _, e := range input.Emails {
		if !validation.IsValidEmail(e) {
			addFieldErr("emails", "wrong format")
		}
	}

	// phone numbers valid format (maybe depending on country?), no letters
	for _, p := range input.PhoneNumbers {
		if !validation.IsValidPhoneNumber(p) {
			addFieldErr("phone_numbers", "wrong format")
		}
	}

	// -90º <= latitude <= 90º, required
	if input.Latitude == nil {
		addFieldErr("latitude", "required")
	} else if *input.Latitude < -90. || *input.Latitude > 90. {
		addFieldErr("latitude", "should be in range [-90, 90]")
	}

	// -180º <= longitude <= 180º, required
	if input.Longitude == nil {
		addFieldErr("longitude", "required")
	} else if *input.Longitude < -180. || *input.Longitude > 180. {
		addFieldErr("longitude", "should be in range [-180, 180]")
	}

	if validationErr == nil {
		return nil
	}

	return validationErr
}
